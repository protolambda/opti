package opti

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/vm"
	. "github.com/protolambda/ztyp/view"
)

// See https://github.com/ethereum/EIPs/blob/master/EIPS/eip-170.md
// ~24.5 KB
var CodeType = BasicListType(ByteType, 0x6000)

type CodeView struct {
	*BasicListView
}

func AsCodeView(v View, err error) (*CodeView, error) {
	li, err := AsBasicList(v, err)
	if err != nil {
		return nil, err
	}
	return &CodeView{BasicListView: li}, err
}

func (v *CodeView) GetOp(pc uint64) (vm.OpCode, error) {
	codeLen, err := v.Length()
	if err != nil {
		return 0, err
	}
	// completed all opcodes? Then the opcode will be 0.
	if pc >= codeLen {
		return 0, nil
	}

	return asOpCode(v.Get(pc))
}

func (v *CodeView) Slice(start uint64, end uint64) ([]byte, error) {
	if start > end {
		return nil, fmt.Errorf("invalid slice, start %d > end %d", start, end)
	}
	length := end - start
	out := make([]byte, length, length)
	for i := uint64(0); i < length; i++ {
		if b, err := asByte(v.Get(start + i)); err != nil {
			return nil, fmt.Errorf("failed to get byte %d from code", start+i)
		} else {
			out[i] = b
		}
	}
	return out, nil
}
