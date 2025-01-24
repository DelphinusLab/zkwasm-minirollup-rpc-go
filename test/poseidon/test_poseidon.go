package main

import (
	"fmt"
	"github.com/DelphinusLab/zkwasm-minirollup-rpc-go/zkwasm"
	"math/big"
)

func main() {
	input := []*zkwasm.Field{zkwasm.NewField(big.NewInt(123)), zkwasm.NewField(big.NewInt(456))}
	hashResult := zkwasm.PoseidonHash(input)
	fmt.Println(hashResult.String())
}
