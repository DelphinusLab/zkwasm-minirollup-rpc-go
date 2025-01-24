package zkwasm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

type ZKWasmAppRpc struct {
	baseURL string
	client  *http.Client
}

func NewZKWasmAppRpc(baseURL string) *ZKWasmAppRpc {
	return &ZKWasmAppRpc{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (rpc *ZKWasmAppRpc) sendRawTransaction(cmd []*big.Int, prikey string) (map[string]interface{}, error) {
	data := Sign(cmd, prikey)
	jsonData, err := json.Marshal(data)
	fmt.Println("sendRawTransaction:", string(jsonData))
	if err != nil {
		return nil, err
	}
	resp, err := rpc.client.Post(fmt.Sprintf("%s/send", rpc.baseURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("SendTransactionError")
}

func (rpc *ZKWasmAppRpc) SendTransaction(cmd []*big.Int, prikey string) (string, error) {
	resp, err := rpc.sendRawTransaction(cmd, prikey)
	if err != nil {
		return "", err
	}
	fmt.Println("resp:", resp)
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		jobStatus, err := rpc.queryJobStatus(resp["jobid"].(string))
		if err != nil {
			continue
		}
		if jobStatus != nil {
			if _, ok := jobStatus["finishedOn"]; ok && jobStatus["failedReason"] == nil {
				returnValue := jobStatus["returnvalue"].(map[string]interface{})
				marshal, jsonErr := json.Marshal(returnValue)
				if jsonErr != nil {
					return "", jsonErr
				}
				return string(marshal), nil
			} else if jobStatus["failedReason"] != nil {
				return "", errors.New(jobStatus["failedReason"].(string))
			}
		}
	}
	return "", errors.New("MonitorTransactionFail")
}

func (rpc *ZKWasmAppRpc) QueryState(prikey string) (map[string]interface{}, error) {
	data := Query(prikey)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	fmt.Println("QueryState:", string(jsonData))
	resp, err := rpc.client.Post(fmt.Sprintf("%s/query", rpc.baseURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("UnexpectedResponseStatus")
}

func (rpc *ZKWasmAppRpc) QueryConfig() (map[string]interface{}, error) {
	resp, err := rpc.client.Post(fmt.Sprintf("%s/config", rpc.baseURL), "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("QueryConfigError")
}

func (rpc *ZKWasmAppRpc) CreateCommand(nonce, command *big.Int, params []*big.Int) []*big.Int {
	cmd := new(big.Int).Lsh(nonce, 16)
	cmd.Add(cmd, new(big.Int).Lsh(big.NewInt(int64(len(params)+1)), 8))
	cmd.Add(cmd, command)

	buf := []*big.Int{cmd}
	buf = append(buf, params...)
	return buf
}

func (rpc *ZKWasmAppRpc) queryJobStatus(jobID string) (map[string]interface{}, error) {
	resp, err := rpc.client.Get(fmt.Sprintf("%s/job/%s", rpc.baseURL, jobID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("QueryJobError")
}

func (rpc *ZKWasmAppRpc) GetNonce(prikey string) (*big.Int, error) {
	state, err := rpc.QueryState(prikey)
	if err != nil {
		return big.NewInt(0), err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(state["data"].(string)), &data)
	if err != nil {
		return big.NewInt(0), err
	}
	player, ok := data["player"]
	if !ok {
		fmt.Println("player field does not exist")
		return big.NewInt(0), nil
	} else if player == nil {
		return big.NewInt(0), nil
	} else {
		playerMap := player.(map[string]interface{})
		fmt.Println("player:", player)
		return big.NewInt(int64(playerMap["nonce"].(float64))), nil
	}
}
