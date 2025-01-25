package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/DelphinusLab/zkwasm-minirollup-rpc-go/zkwasm"
)

var (
	InitPlayerCmd  = big.NewInt(1)
	BuyElfCmd      = big.NewInt(2)
	CollectCoinCmd = big.NewInt(11)
	CleanRanchCmd  = big.NewInt(4)
	DepositCmd     = big.NewInt(8)
	BuyPropCmd     = big.NewInt(12)
	FeedElfCmd     = big.NewInt(3)
	//const CLEAN_RANCH: u64 = 4; // 清洁牧场
	//const TREAT_ELF: u64 = 5; // 治疗宠物
	//const SELL_ELF: u64 = 6; // 卖出精灵
	//
	//const WITHDRAW: u64 = 7; // 提现
	//const DEPOSIT: u64 = 8; // 充值
	//const BOUNTY: u64 = 9;
	//const BUY_RANCH: u64 = 10; // 购买牧场
	//const COLLECT_GOLD: u64 = 11; // 收集金币
	//
	//const BUY_SLOT: u64 = 13; // 购买宠物槽位
	//
	//const BUY_PROP: u64 = 12; // 购买道具

	SellElfCmd      = big.NewInt(6)
	WithdrawCmd     = big.NewInt(7)
	BuyRanchCmd     = big.NewInt(10)
	BuySlotCmd      = big.NewInt(13)
	TreatmentElfCmd = big.NewInt(5)
)

func main() {
	prikey := "12345"
	pid1, pid2 := zkwasm.GetPid(prikey)

	fmt.Println("pid1:", pid1.Uint64())
	fmt.Println("pid2:", pid2.Uint64())

	zkwamRpc := zkwasm.NewZKWasmAppRpc("http://localhost:3000")
	//zkwamRpc := zkwasm.NewZKWasmAppRpc("https://zk-server.pumpelf.ai")

	// 收集金币
	//nonce, _ := zkwamRpc.GetNonce(prikey)
	//cmd := zkwamRpc.CreateCommand(nonce, big.NewInt(11), big.NewInt(0))
	//fmt.Println("cmd:", cmd)
	//transaction, _ := zkwamRpc.SendTransaction([4]*big.Int{cmd, big.NewInt(1), big.NewInt(1), big.NewInt(0)}, prikey)
	//fmt.Println("transaction:", transaction)
	// 初始化玩家
	//initPlayer(zkwamRpc, prikey)
	// 购买宠物
	//buyElf(zkwamRpc, prikey)

	// 购买食物道具
	//buyFoodProp(zkwamRpc, prikey)
	// 喂食精灵
	//feedElf(zkwamRpc, prikey)

	// 购买健康道具
	//buyHealthProp(zkwamRpc, prikey)

	// 治疗宠物
	//healthElf(zkwamRpc, prikey)
	// 收集金币
	//collectCoin(zkwamRpc, prikey)

	// 清理牧场
	//clearRanch(zkwamRpc, prikey)

	// 充值
	//deposit(zkwamRpc, "126532")
	//withdraw(zkwamRpc, prikey)

	// 查询状态
	getState(zkwamRpc, prikey)
}

func getState(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	state, _ := zkwamRpc.QueryState(prikey)
	fmt.Println("state:", state)
}

// 初始化玩家
func initPlayer(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("初始化玩家")
	cmd := zkwamRpc.CreateCommand(big.NewInt(0), InitPlayerCmd, []*big.Int{})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("初始化玩家:", transaction)
}

// 购买宠物
func buyElf(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("购买宠物")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elfType := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, BuyElfCmd, []*big.Int{ranchId, elfType})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("购买宠物:", transaction)
}

// 收集金币
func collectCoin(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("收集金币")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elfId := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, CollectCoinCmd, []*big.Int{ranchId, elfId})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("收集金币:", transaction)
}

// 清理牧场
func clearRanch(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("清理牧场")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	cmd := zkwamRpc.CreateCommand(nonce, CleanRanchCmd, []*big.Int{ranchId})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("清理牧场:", transaction)
}

// 购买食物道具
func buyFoodProp(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("购买食物道具")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	propType := big.NewInt(4)
	cmd := zkwamRpc.CreateCommand(nonce, BuyPropCmd, []*big.Int{ranchId, propType})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("购买食物道具:", transaction)
}

// 购买道具
func buyHealthProp(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("购买健康道具")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	propType := big.NewInt(6)
	cmd := zkwamRpc.CreateCommand(nonce, BuyPropCmd, []*big.Int{ranchId, propType})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("购买健康道具:", transaction)
}

func healthElf(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("治疗宠物")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranch_id := big.NewInt(1)
	elf_id := big.NewInt(1)
	prop_type := big.NewInt(6)
	cmd := zkwamRpc.CreateCommand(nonce, TreatmentElfCmd, []*big.Int{ranch_id, elf_id, prop_type})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("治疗宠物:", transaction)
}

// 喂食精灵
func feedElf(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("喂食精灵")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	ranchId := big.NewInt(1)
	elf_id := big.NewInt(1)
	prop_type := big.NewInt(4)
	cmd := zkwamRpc.CreateCommand(nonce, FeedElfCmd, []*big.Int{ranchId, elf_id, prop_type})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("喂食精灵:", transaction)
}

func deposit(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("充值")

	initPlayer(zkwamRpc, prikey)
	pid1, pid2 := zkwasm.GetPid("12345")
	fmt.Println("pid1:", pid1.Uint64())
	fmt.Println("pid2:", pid2.Uint64())
	ranchId := big.NewInt(1)
	propType := big.NewInt(1)
	nonce, _ := zkwamRpc.GetNonce(prikey)
	cmd := zkwamRpc.CreateCommand(nonce, DepositCmd, []*big.Int{pid1, pid2, ranchId, propType})
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("充值:", transaction)
}

// 提现
func withdraw(zkwamRpc *zkwasm.ZKWasmAppRpc, prikey string) {
	fmt.Println("提现")
	nonce, _ := zkwamRpc.GetNonce(prikey)
	address := common.HexToAddress("0xae1e3ffa0a95b7c11cfd0a8f02d3250f20b51ff2")
	cmd, _ := zkwamRpc.ComposeWithdrawParams(address, nonce, WithdrawCmd, big.NewInt(1), big.NewInt(0))
	transaction, _ := zkwamRpc.SendTransaction(cmd, prikey)
	fmt.Println("提现:", transaction)
}
