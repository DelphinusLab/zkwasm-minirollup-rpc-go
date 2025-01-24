package zkwasm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// Helper function to convert big-endian hex string to an integer
func BigEndianHexToInt(hexString string) *big.Int {
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}
	// Ensure even length by adding a leading zero if necessary
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}
	// Parse the hex string into a big integer
	result := new(big.Int)
	result.SetString(hexString, 16)
	return result
}

// Helper function to convert little-endian hex string to an integer
func LittleEndianHexToInt(hexString string) *big.Int {
	if strings.HasPrefix(hexString, "0x") {
		hexString = hexString[2:]
	}
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}
	// Reverse the byte order
	reversedHex := ""
	for i := len(hexString) - 2; i >= 0; i -= 2 {
		reversedHex += hexString[i : i+2]
	}
	// Parse the reversed hex string into a big integer
	result := new(big.Int)
	result.SetString(reversedHex, 16)
	return result
}

// Convert a u8 array (byte array) to a hex string
func U8ToHex(u8Array []byte) string {
	u8Array = reverseBytes(u8Array)
	return hex.EncodeToString(u8Array)
}

// Convert a big integer to a little-endian hex string
func BnToHexLe(n *big.Int, sz int) string {
	// Convert the integer to bytes in little-endian order
	bytesLe := n.Bytes()
	// Ensure the result is 32 bytes
	if len(bytesLe) < sz {
		pad := make([]byte, sz-len(bytesLe))
		bytesLe = append(pad, bytesLe...)
	}
	// Convert the bytes to hex
	return U8ToHex(bytesLe)
}

// LeHexInt represents a little-endian hex integer
type LeHexInt struct {
	HexStr string
}

// Convert LeHexInt to an integer
func (l *LeHexInt) ToInt() *big.Int {
	return LittleEndianHexToInt(l.HexStr)
}

// Convert LeHexInt to a u64 array
func (l *LeHexInt) ToU64Array() []*big.Int {
	// Convert to integer first
	num := l.ToInt()
	// Split the big integer into 4 uint64 values
	values := make([]*big.Int, 4)
	mask := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)
	for i := 0; i < 4; i++ {
		values[i] = new(big.Int).And(num, mask) // Take the least significant 64 bits
		num.Rsh(num, 64)                        // Shift the number by 64 bits for the next value
	}
	return values
}

// VerifySign verifies a signature
func VerifySign(msg *LeHexInt, pkx, pky, rx, ry *LeHexInt, s *LeHexInt) bool {
	l := PointBase().Mul((*Field)(NewCurveField(s.ToInt()))) // Mul 方法需要指针接收器
	pkey := &Point{(*Field)(NewCurveField(pkx.ToInt())), (*Field)(NewCurveField(pky.ToInt()))}

	// 将 Point 存储在变量中，避免临时值
	r := &Point{(*Field)(NewCurveField(rx.ToInt())), (*Field)(NewCurveField(ry.ToInt()))}

	// 使用 r 调用 Add 方法
	r = r.Add(pkey.Mul((*Field)(NewCurveField(msg.ToInt())))) // Add 方法返回一个新的 Point 指针

	// 直接在 r.x 上调用 Neg 方法
	negr := &Point{r.x.Neg(), r.y} // r 已经是指针类型，所以可以直接调用

	// 最后的 IsZero 检查
	return l.Add(negr).IsZero() // Add 和 IsZero 都是指针接收器的方法
}

// Sign signs a command using a private key
func Sign(cmd []*big.Int, prikey string) map[string]string {
	pkey := PrivateKeyFromString(prikey)
	r := pkey.R()
	R := PointBase().Mul((*Field)(r))

	var H *big.Int
	fvalues := []*Field{}
	h := big.NewInt(0)

	for i := 0; i < len(cmd); {
		currentI := i // 保存当前外层循环的 i
		v := big.NewInt(0)
		j := 0
		for ; j < 3; j++ {
			if currentI+j < len(cmd) {
				// 计算 v 和 h 的位移（使用 currentI 而非动态 i）
				shift := new(big.Int).Lsh(cmd[currentI+j], uint(64*j))
				v.Add(v, shift)

				hshiftBits := 64 * (currentI + j) // 正确位移位数
				hshift := new(big.Int).Lsh(cmd[currentI+j], uint(hshiftBits))
				h.Add(h, hshift)
			}
		}
		i += j // 更新外层循环的 i
		fvalues = append(fvalues, NewField(v))
	}
	H = PoseidonHash(fvalues).v
	fmt.Println("H:", H)
	hbn := H
	msgbn := h

	// 仅在曲线运算时使用 CurveField
	hbnCurve := NewCurveField(hbn)
	S := r.Add(pkey.key.Mul(hbnCurve))
	pubkey := pkey.PublicKey()
	fmt.Println("msgbn:", msgbn.String())
	fmt.Println("hbn:", hbn.String())
	data := map[string]string{
		"msg":  BnToHexLe(msgbn, len(cmd)*8),
		"hash": BnToHexLe(hbn, 32),
		"pkx":  BnToHexLe(pubkey.key.x.v, 32),
		"pky":  BnToHexLe(pubkey.key.y.v, 32),
		"sigx": BnToHexLe(R.x.v, 32),
		"sigy": BnToHexLe(R.y.v, 32),
		"sigr": BnToHexLe(S.v, 32),
	}
	return data
}

// Query queries the public key associated with a private key
func Query(prikey string) map[string]string {
	pkey := PrivateKeyFromString(prikey)
	pubkey := pkey.PublicKey()
	data := map[string]string{
		"pkx": BnToHexLe(pubkey.key.x.v, 32),
	}
	return data
}

// GetPid retrieves the PID associated with a private key
func GetPid(prikey string) (*big.Int, *big.Int) {
	pkey := PrivateKeyFromString(prikey)
	pubkey := pkey.PublicKey()
	pidKey := BnToHexLe(pubkey.key.x.v, 32)
	pidAll := (&LeHexInt{pidKey}).ToU64Array()
	pid1, pid2 := pidAll[1], pidAll[2]
	return pid1, pid2
}
