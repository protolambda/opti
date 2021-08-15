package opti

import (
	. "github.com/protolambda/ztyp/view"
)

// Needs to be as big as memory, all of it can be returned
var RetDataType = ListType(Bytes32Type, 64<<20)

type RetDataView struct {
	*ComplexListView
}

func EmptyReturnData() *RetDataView {
	return &RetDataView{ComplexListView: RetDataType.Default(nil).(*ComplexListView)}
}

func AsRetDataView(v View, err error) (*RetDataView, error) {
	li, err := AsComplexList(v, err)
	if err != nil {
		return nil, err
	}
	return &RetDataView{ComplexListView: li}, err
}

func (v *RetDataView) SetWord(i uint64, w [32]byte) error {
	return v.Set(i, (*Bytes32View)(&w))
}

func (v *RetDataView) GetWord(i uint64) ([32]byte, error) {
	return AsRoot(v.Get(i))
}
