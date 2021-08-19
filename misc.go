package opti

import (
	"github.com/ethereum/go-ethereum/core/vm"
	. "github.com/protolambda/ztyp/view"
)

const Bytes32Type = RootType

type Bytes32View = RootView

const AddressType = SmallByteVecMeta(20)

func AsAddress(v View, err error) (Address, error) {
	c, err := AsSmallByteVec(v, err)
	if err != nil {
		return Address{}, err
	}
	var out Address
	copy(out[:], c)
	return out, nil
}

type Address [20]byte

func (addr Address) View() SmallByteVecView {
	// overlay on a copy of the value
	return SmallByteVecView(addr[:])
}

func asExecMode(v View, err error) (ExecMode, error) {
	out, errOut := AsByte(v, err)
	return ExecMode(out), errOut
}

func asOpCode(v View, err error) (vm.OpCode, error) {
	out, errOut := AsByte(v, err)
	return vm.OpCode(out), errOut
}

func asByte(v View, err error) (byte, error) {
	out, errOut := AsByte(v, err)
	return byte(out), errOut
}

func asUint64(v View, err error) (uint64, error) {
	out, errOut := AsUint64(v, err)
	return uint64(out), errOut
}

func asBool(v View, err error) (bool, error) {
	out, errOut := AsBool(v, err)
	return bool(out), errOut
}
