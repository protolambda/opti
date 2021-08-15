package opti

import (
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
