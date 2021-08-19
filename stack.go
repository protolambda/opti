package opti

import (
	"fmt"
	. "github.com/protolambda/ztyp/view"
)

// EVM stack is max 1024 words
var StackType = ListType(Bytes32Type, 1024)

type StackView struct {
	*ComplexListView
}

func AsStackView(v View, err error) (*StackView, error) {
	li, err := AsComplexList(v, err)
	if err != nil {
		return nil, err
	}
	return &StackView{ComplexListView: li}, err
}

func (v *StackView) SetWord(i uint64, w [32]byte) error {
	return v.Set(i, (*Bytes32View)(&w))
}

func (v *StackView) GetWord(i uint64) ([32]byte, error) {
	return AsRoot(v.Get(i))
}

func (v *StackView) Back(n uint64) ([32]byte, error) {
	length, err := v.Length()
	if err != nil {
		return [32]byte{}, err
	}
	if n+1 > length {
		return [32]byte{}, fmt.Errorf("bad stack back access", n)
	}
	return v.GetWord(length - n - 1)
}

func (v *StackView) PushWord(w [32]byte) error {
	return v.Append((*Bytes32View)(&w))
}

func (v *StackView) PopWord() error {
	return v.Pop()
}
