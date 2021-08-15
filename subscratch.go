package opti

import (
	. "github.com/protolambda/ztyp/view"
)

// 1024 words to track sub-step progress. Not to be confused with the memory scratchpad slots.
// TODO: any operations that need more scratch space?
var SubDataType = VectorType(Bytes32Type, 1024)

type SubDataView struct {
	*ComplexVectorView
}

func AsSubDataView(v View, err error) (*SubDataView, error) {
	vec, err := AsComplexVector(v, err)
	if err != nil {
		return nil, err
	}
	return &SubDataView{ComplexVectorView: vec}, err
}

func (v *SubDataView) SetWord(i uint64, w [32]byte) error {
	return v.Set(i, (*Bytes32View)(&w))
}

func (v *SubDataView) GetWord(i uint64) ([32]byte, error) {
	return AsRoot(v.Get(i))
}
