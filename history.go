package opti

import (
	. "github.com/protolambda/ztyp/view"
)

// 256 most recent block hashes, "recent history", are available for direct access in the EVM
var HistType = VectorType(RootType, 256)

type HistView struct {
	*ComplexVectorView
}

func AsHistView(v View, err error) (*HistView, error) {
	vec, err := AsComplexVector(v, err)
	if err != nil {
		return nil, err
	}
	return &HistView{ComplexVectorView: vec}, err
}

func (v *HistView) SetWord(i uint64, w [32]byte) error {
	return v.Set(i, (*Bytes32View)(&w))
}

func (v *HistView) GetWord(i uint64) ([32]byte, error) {
	return AsRoot(v.Get(i))
}
