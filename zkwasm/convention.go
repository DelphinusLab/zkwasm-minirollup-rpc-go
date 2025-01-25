package zkwasm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
	"strings"
)

func bytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func bytesToDecimal(bytes []byte) string {
	var sb strings.Builder
	for _, b := range bytes {
		sb.WriteString(fmt.Sprintf("%02d", b))
	}
	return sb.String()
}

func (pc *PlayerConvention) ComposeWithdrawParams(address common.Address, nonce, command, amount, tokenIndex *big.Int) ([]*big.Int, error) {
	addressBytes := address.Bytes()
	firstLimb := new(big.Int).SetBytes(reverseBytes(addressBytes[:4]))
	sndLimb := new(big.Int).SetBytes(reverseBytes(addressBytes[4:12]))
	thirdLimb := new(big.Int).SetBytes(reverseBytes(addressBytes[12:20]))
	one := new(big.Int).Add(new(big.Int).Lsh(firstLimb, 32), amount)
	return pc.CreateCommand(nonce, command, []*big.Int{tokenIndex, one, sndLimb, thirdLimb}), nil
}

func decodeWithdraw(txdata []byte) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if len(txdata) > 1 {
		for i := 0; i < len(txdata); i += 32 {
			extra := txdata[i : i+4]
			address := txdata[i+4 : i+24]
			amount := txdata[i+24 : i+32]
			amountInWei, err := strconv.ParseInt(bytesToDecimal(amount), 10, 64)
			if err != nil {
				return nil, err
			}
			result = append(result, map[string]interface{}{
				"op":      extra[0],
				"index":   extra[1],
				"address": "0x" + bytesToHex(address),
				"amount":  amountInWei,
			})
		}
	}
	return result, nil
}

func reverseBytes(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

type PlayerConvention struct {
	processingKey   string
	rpc             *ZKWasmAppRpc
	commandDeposit  *big.Int
	commandWithdraw *big.Int
}

func NewPlayerConvention(key string, rpc *ZKWasmAppRpc, commandDeposit, commandWithdraw *big.Int) *PlayerConvention {
	return &PlayerConvention{
		processingKey:   key,
		rpc:             rpc,
		commandDeposit:  commandDeposit,
		commandWithdraw: commandWithdraw,
	}
}

func (pc *PlayerConvention) CreateCommand(nonce, command *big.Int, params []*big.Int) []*big.Int {
	cmd := new(big.Int).Lsh(nonce, 16)
	cmd.Add(cmd, new(big.Int).Lsh(big.NewInt(int64(len(params)+1)), 8))
	cmd.Add(cmd, command)

	buf := []*big.Int{cmd}
	buf = append(buf, params...)
	return buf
}

func (pc *PlayerConvention) getConfig() (map[string]interface{}, error) {
	return pc.rpc.QueryConfig()
}

func (pc *PlayerConvention) getState() (map[string]interface{}, error) {
	state, err := pc.rpc.QueryState(pc.processingKey)
	if err != nil {
		return nil, err
	}
	var parsedState map[string]interface{}
	if err := json.Unmarshal([]byte(state["data"].(string)), &parsedState); err != nil {
		return nil, err
	}
	return parsedState, nil
}

func (pc *PlayerConvention) getNonce() (*big.Int, error) {
	state, err := pc.rpc.QueryState(pc.processingKey)
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

func (pc *PlayerConvention) Deposit(pid1, pid2, amount *big.Int) (string, error) {
	nonce, err := pc.getNonce()
	if err != nil {
		return "", err
	}
	return pc.rpc.SendTransaction(
		pc.CreateCommand(nonce, pc.commandDeposit, []*big.Int{pid1, pid2, amount}), pc.processingKey)
}

func (pc *PlayerConvention) WithdrawRewards(address common.Address, amount *big.Int) (string, error) {
	nonce, err := pc.getNonce()
	if err != nil {
		return "", err
	}
	params, err := pc.ComposeWithdrawParams(address, nonce, pc.commandWithdraw, amount, big.NewInt(0))
	if err != nil {
		return "", err
	}
	return pc.rpc.SendTransaction(params, pc.processingKey)
}
