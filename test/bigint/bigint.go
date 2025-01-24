package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func main() {
	hexStr := "0x0190f922d97c8a7dcf0a142a3be27749d1c64bc22f1c556aaa24925d158cac56"
	// 使用 go-ethereum 的 common.FromHex 函数将十六进制字符串转换为字节切片
	bytes := common.FromHex(hexStr)

	// 将字节切片转换为 big.Int
	bigInt := new(big.Int).SetBytes(bytes)

	// 打印十进制结果
	fmt.Println("Decimal:", bigInt.String())

	// 打印十六进制结果（带 0x 前缀）
	fmt.Println("Hex:", "0x"+bigInt.Text(16))
}
