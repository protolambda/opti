package opti

import (
	. "github.com/protolambda/ztyp/view"
)

// TODO: 64 MB memory maximum enough or too much? Every 2x makes the tree a layer deeper,
// but otherwise not much cost for unused space
var MemoryType = ListType(ByteType, 64<<20)

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

func (v *MemoryView) SetByte(i uint64, w byte) error {
	return v.Set(i, (*ByteView)(&w))
}

func (v *MemoryView) GetByte(i uint64) (byte, error) {
	return asByte(v.Get(i))
}

func (v *MemoryView) AppendZeroBytes32() error {
	// TODO: if the length aligns to 32, this can be optimized
	for i := 0; i < 32; i++ {
		if err := v.Append(ByteView(0)); err != nil {
			return err
		}
	}
	return nil
}

func (v *MemoryView) AppendZeroByte() error {
	return v.Append(ByteView(0))
}
