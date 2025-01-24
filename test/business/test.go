package main

import (
	"fmt"
	"math/big"

	"github.com/DelphinusLab/zkwasm-minirollup-rpc-go/zkwasm"
)

var (
	InitPlayerCmd  = big.NewInt(1)
	BuyElfCmd      = big.NewInt(2)
	CollectCoinCmd = big.NewInt(11)
	CleanRanchCmd  = big.NewInt(4)
	DepositCmd     = big.NewInt(8)
)

func main() {
	prikey := "1234"
	pid1, pid2 := zkwasm.GetPid(prikey)

	fmt.Println("pid1:", pid1.Int64())
	fmt.Println("pid2:", pid2.Int64())

	data := zkwasm.Query(prikey)
	fmt.Println("data:", data)

	zkwamRpc := zkwasm.NewZKWasmAppRpc("http://localhost:3000")

	state, _ := zkwamRpc.QueryState(prikey)
	fmt.Println("state:", state)
	// 收集金币
	//nonce, _ := zkwamRpc.GetNonce(prikey)
	//cmd := zkwamRpc.CreateCommand(nonce, big.NewInt(11), big.NewInt(0))
	//fmt.Println("cmd:", cmd)
	//transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, big.NewInt(1), big.NewInt(1), big.NewInt(0)}, prikey)
	//fmt.Println("transaction:", transaction)
	// 初始化玩家
	initPlayer(zkwamRpc, prikey)
	// 购买宠物
	//buyElf(zkwamRpc, prikey)
	// 收集金币
	//collectCoin(zkwamRpc, prikey)

	// 清理牧场
	//clearRanch(zkwamRpc, prikey)

	// 充值
	//deposit(zkwamRpc, "123123123")

	// 查询状态
	state, _ = zkwamRpc.QueryState(prikey)
	fmt.Println("state:", state)
}

// 初始化玩家
func initPlayer(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	// 获取配置
	cmd := zkwamRpc.CreateCommand(big.NewInt(0), InitPlayerCmd, []*big.Int{})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("transaction:", transaction)
}

// 购买宠物
func buyElf(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elfType := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, BuyElfCmd, []*big.Int{ranchId, elfType})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("transaction:", transaction)
}

// 收集金币
func collectCoin(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elfId := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, CollectCoinCmd, []*big.Int{ranchId, elfId})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("transaction:", transaction)
}

// 清理牧场
func clearRanch(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, CleanRanchCmd, []*big.Int{ranchId})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("transaction:", transaction)
}

func deposit(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {

	initPlayer(zkwamRpc, prikey)
	pid1, pid2 := zkwasm.GetPid(prikey)
	ranchId := big.NewInt(1)
	propType := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(big.NewInt(0), DepositCmd, []*big.Int{pid1, pid2, ranchId, propType})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("transaction:", transaction)
}
