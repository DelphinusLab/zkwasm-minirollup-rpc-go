package zkwasm

import (
	"math/big"
)

func pow5(a *Field) *Field {
	return a.Mul(a.Mul(a.Mul(a.Mul(a))))
}

func sboxFull(state []*Field) []*Field {
	for i := 0; i < len(state); i++ {
		tmp := state[i].Mul(state[i])
		state[i] = state[i].Mul(tmp)
		state[i] = state[i].Mul(tmp)
	}
	return state
}

func addConstants(a, b []*Field) []*Field {
	for i := range a {
		a[i] = a[i].Add(b[i])
	}
	return a
}

func apply(matrix [][]*Field, vector []*Field) []*Field {
	result := make([]*Field, len(matrix))
	for i := range result {
		result[i] = NewField(big.NewInt(0))
	}
	for i := range matrix {
		for j := range matrix[i] {
			result[i] = result[i].Add(matrix[i][j].Mul(vector[j]))
		}
	}
	return result
}

func applySparseMatrix(matrix map[string][]string, state []*Field) []*Field {
	words := make([]*Field, len(state))
	for i := range state {
		words[i] = NewField(state[i].v)
	}
	state0 := NewField(big.NewInt(0))

	for i := range words {
		rowI := matrix["row"][i]
		f := NewField(HexToBigInt(rowI))
		state0 = state0.Add(f.Mul(words[i]))
	}
	state[0] = state0
	for i := 1; i < len(words); i++ {
		hat := NewField(HexToBigInt(matrix["col_hat"][i-1]))
		state[i] = hat.Mul(words[0]).Add(words[i])
	}
	return state
}

func toFieldArray(arr []string) []*Field {
	fields := make([]*Field, len(arr))
	for i, value := range arr {
		fields[i] = NewField(HexToBigInt(value))
	}
	return fields
}

func convertTo2DStringSlice(arr []interface{}) [][]string {
	result := make([][]string, len(arr))
	for i, v := range arr {
		result[i] = toStringSlice(v.([]interface{}))
	}
	return result
}

func toFieldMatrix(arr [][]string) [][]*Field {
	matrix := make([][]*Field, len(arr))
	for i, row := range arr {
		matrix[i] = toFieldArray(row)
	}
	return matrix
}

type Poseidon struct {
	state      []*Field
	absortbing []*Field
	squeezed   bool
	config     map[string]interface{}
}

func NewPoseidon(config map[string]interface{}) *Poseidon {
	state := make([]*Field, int(config["t"].(float64)))
	for i := range state {
		state[i] = NewField(big.NewInt(0))
	}
	state[0] = NewField(HexToBigInt("0000000000000000000000000000000000000000000000010000000000000000"))
	return &Poseidon{
		state:      state,
		absortbing: []*Field{},
		squeezed:   false,
		config:     config,
	}
}

func (p *Poseidon) getState() []*Field {
	return p.state
}

func (p *Poseidon) permute() {
	rf := int(p.config["r_f"].(float64)) / 2

	// First half of full rounds
	startConstants := p.config["constants"].(map[string]interface{})["start"].([]interface{})
	startConstantsStr := make([][]string, len(startConstants))
	for i, v := range startConstants {
		startConstantsStr[i] = toStringSlice(v.([]interface{}))
	}
	p.state = addConstants(p.state, toFieldArray(startConstantsStr[0]))
	for i := 1; i < rf; i++ {
		p.state = sboxFull(p.state)
		p.state = addConstants(p.state, toFieldArray(startConstantsStr[i]))
		p.state = apply(toFieldMatrix(convertTo2DStringSlice(p.config["mds_matrices"].(map[string]interface{})["mds"].([]interface{}))), p.state)
	}
	p.state = sboxFull(p.state)
	p.state = addConstants(p.state, toFieldArray(startConstantsStr[len(startConstantsStr)-1]))
	p.state = apply(toFieldMatrix(convertTo2DStringSlice(p.config["mds_matrices"].(map[string]interface{})["pre_sparse_mds"].([]interface{}))), p.state)

	// Partial rounds
	partialConstants := p.config["constants"].(map[string]interface{})["partial"].([]interface{})
	partialConstantsStr := toStringSlice(partialConstants)
	sparseMatrices := p.config["mds_matrices"].(map[string]interface{})["sparse_matrices"].([]interface{})
	sparseMatricesMap := make([]map[string][]string, len(sparseMatrices))
	for i, v := range sparseMatrices {
		sparseMatricesMap[i] = toStringMap(v.(map[string]interface{}))
	}
	for i := 0; i < len(partialConstantsStr) && i < len(sparseMatricesMap); i++ {
		p.state[0] = pow5(p.state[0])
		p.state[0] = p.state[0].Add(NewField(HexToBigInt(partialConstantsStr[i])))
		applySparseMatrix(sparseMatricesMap[i], p.state)
	}

	// Second half of the full rounds
	endConstants := p.config["constants"].(map[string]interface{})["end"].([]interface{})
	endConstantsStr := make([][]string, len(endConstants))
	for i, v := range endConstants {
		endConstantsStr[i] = toStringSlice(v.([]interface{}))
	}
	for _, constants := range endConstantsStr {
		p.state = sboxFull(p.state)
		p.state = addConstants(p.state, toFieldArray(constants))
		p.state = apply(toFieldMatrix(convertTo2DStringSlice(p.config["mds_matrices"].(map[string]interface{})["mds"].([]interface{}))), p.state)
	}
	p.state = sboxFull(p.state)
	p.state = apply(toFieldMatrix(convertTo2DStringSlice(p.config["mds_matrices"].(map[string]interface{})["mds"].([]interface{}))), p.state)
}

func (p *Poseidon) updateExact(elements []*Field) *Field {
	if p.squeezed {
		panic("Cannot update after squeeze")
	}
	if len(elements) != int(p.config["rate"].(float64)) {
		panic("Invalid input size")
	}

	for j := 0; j < int(p.config["rate"].(float64)); j++ {
		p.state[j+1] = p.state[j+1].Add(elements[j])
	}
	p.permute()
	return p.state[1]
}

func (p *Poseidon) update(elements []*Field) {
	if p.squeezed {
		panic("Cannot update after squeeze")
	}
	for i := 0; i < len(elements); i += int(p.config["rate"].(float64)) {
		if i+int(p.config["rate"].(float64)) > len(elements) {
			p.absortbing = elements[i:]
		} else {
			chunk := elements[i : i+int(p.config["rate"].(float64))]
			for j := 0; j < int(p.config["rate"].(float64)); j++ {
				p.state[j+1] = p.state[j+1].Add(chunk[j])
			}
			p.permute()
			p.absortbing = []*Field{}
		}
	}
}

func (p *Poseidon) squeeze() *Field {
	lastChunk := p.absortbing
	lastChunk = append(lastChunk, NewField(big.NewInt(1)))
	for i := 0; i < len(lastChunk) && i < len(p.state)-1; i++ {
		p.state[i+1] = p.state[i+1].Add(lastChunk[i])
	}
	p.permute()
	p.absortbing = []*Field{}
	return p.state[1]
}

func toStringSlice(arr []interface{}) []string {
	result := make([]string, len(arr))
	for i, v := range arr {
		result[i] = v.(string)
	}
	return result
}

func toStringMap(m map[string]interface{}) map[string][]string {
	result := make(map[string][]string)
	for k, v := range m {
		result[k] = toStringSlice(v.([]interface{}))
	}
	return result
}

func PoseidonHash(inputs []*Field) *Field {
	if len(inputs) == 0 {
		panic("Invalid input size")
	}
	hasher := NewPoseidon(Config)
	hasher.update(inputs)
	return hasher.squeeze()
}
