package opti

import (
	. "github.com/protolambda/ztyp/view"
)

// TODO: 64 MB memory maximum enough or too much? Every 2x makes the tree a layer deeper,
// but otherwise not much cost for unused space
var MemoryType = ListType(Bytes32Type, 64<<20)

type MemoryView struct {
	*ComplexListView
}

func AsMemoryView(v View, err error) (*MemoryView, error) {
	li, err := AsComplexList(v, err)
	if err != nil {
		return nil, err
	}
	return &MemoryView{ComplexListView: li}, err
}

func (v *MemoryView) SetWord(i uint64, w [32]byte) error {
	return v.Set(i, (*Bytes32View)(&w))
}

func (v *MemoryView) GetWord(i uint64) ([32]byte, error) {
	return AsRoot(v.Get(i))
}
