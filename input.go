package opti

import (
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
