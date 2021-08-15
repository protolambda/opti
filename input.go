package opti

import (
	"fmt"
	. "github.com/protolambda/ztyp/view"
)

// Assuming a tx input can be max 400M gas, and 4 gas is paid per zero byte, then put a 100M limit on input.
var InputType = ListType(ByteType, 100_000_000)

type InputView struct {
	*BasicListView
}

func AsInputView(v View, err error) (*InputView, error) {
	li, err := AsBasicList(v, err)
	if err != nil {
		return nil, err
	}
	return &InputView{BasicListView: li}, err
}

func (v *InputView) Slice(start uint64, end uint64) ([]byte, error) {
	if start > end {
		return nil, fmt.Errorf("invalid slice, start %d > end %d", start, end)
	}
	length := end - start
	out := make([]byte, length, length)
	for i := uint64(0); i < length; i++ {
		if b, err := asByte(v.Get(start + i)); err != nil {
			return nil, fmt.Errorf("failed to get byte %d from input", start+i)
		} else {
			out[i] = b
		}
	}
	return out, nil
}
