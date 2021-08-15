package opti

import (
	. "github.com/protolambda/ztyp/view"
)

const (
	fieldStateRoot = iota
	fieldBlockHashes
	fieldCoinbase
	fieldGasLimit
	fieldBlockNumber
	fieldTime
	fieldDifficulty
	fieldBaseFee
	fieldOrigin
	fieldTxIndex
	fieldGasPrice
	fieldTo
	fieldCreate
	fieldCallDepth
	fieldCaller
	fieldMemory
	fieldMemLastGasCost
	fieldStack
	fieldReturnData
	fieldCode
	fieldCodeHash
	fieldCodeAddr
	fieldInput
	fieldGas
	fieldValue
	fieldOp
	fieldPc
	fieldSubIndex
	fieldSubRemaining
	fieldSubScratch
	fieldReturnToStep
)

const Bytes32Type = RootType

type Bytes32View = RootView

var BlockHashesType = VectorType(RootType, 256)

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

// TODO: 64 MB memory maximum enough or too much? Every 2x makes the tree a layer deeper,
// but otherwise not much cost for unused space
var MemoryType = ListType(Bytes32Type, 64<<20)

// EVM stack is max 1024 words
var StackType = ListType(Bytes32Type, 1024)

// Needs to be as big as memory, all of it can be returned
var ReturnDataType = ListType(Bytes32Type, 64<<20)

// See https://github.com/ethereum/EIPs/blob/master/EIPS/eip-170.md
// ~24.5 KB
var CodeType = ListType(ByteType, 0x6000)

// Assuming a tx input can be max 400M gas, and 4 gas is paid per zero byte, then put a 100M limit on input.
var InputType = ListType(ByteType, 100_000_000)

// 1024 words to track sub-step progress. Not to be confused with the memory scratchpad slots.
// TODO: any operations that need more scratch space?
var ScratchType = VectorType(Bytes32Type, 1024)

// TODO: views over memory, return data, stack, code, input and scratch

type StepView struct {
	*ContainerView
}

func AsStep(v View, err error) (*StepView, error) {
	cv, err := AsContainer(v, err)
	return &StepView{ContainerView: cv}, err

}
func (v *StepView) CopyView() (*StepView, error) {
	return AsStep(v.Copy())
}

func asUint64(v View, err error) (uint64, error) {
	out, errOut := AsUint64(v, err)
	return uint64(out), errOut
}

func asBool(v View, err error) (bool, error) {
	out, errOut := AsBool(v, err)
	return bool(out), errOut
}

func (v *StepView) GetStateRoot() ([32]byte, error)    { return AsRoot(v.Get(fieldStateRoot)) }
func (v *StepView) SetStateRoot(p [32]byte) error      { return v.Set(fieldStateRoot, (*Bytes32View)(&p)) }
func (v *StepView) GetBlockHashes() (View, error)      { return v.Get(fieldBlockHashes) }
func (v *StepView) SetBlockHashes(p View) error        { return v.Set(fieldBlockHashes, p) }
func (v *StepView) GetCoinbase() (Address, error)      { return AsAddress(v.Get(fieldCoinbase)) }
func (v *StepView) SetCoinbase(p Address) error        { return v.Set(fieldCoinbase, p.View()) }
func (v *StepView) GetGasLimit() (uint64, error)       { return asUint64(v.Get(fieldGasLimit)) }
func (v *StepView) SetGasLimit(p uint64) error         { return v.Set(fieldGasLimit, Uint64View(p)) }
func (v *StepView) GetBlockNumber() (uint64, error)    { return asUint64(v.Get(fieldBlockNumber)) }
func (v *StepView) SetBlockNumber(p uint64) error      { return v.Set(fieldBlockNumber, Uint64View(p)) }
func (v *StepView) GetTime() (uint64, error)           { return asUint64(v.Get(fieldTime)) }
func (v *StepView) SetTime(p uint64) error             { return v.Set(fieldTime, Uint64View(p)) }
func (v *StepView) GetDifficulty() ([32]byte, error)   { return AsRoot(v.Get(fieldDifficulty)) }
func (v *StepView) SetDifficulty(p [32]byte) error     { return v.Set(fieldDifficulty, (*Bytes32View)(&p)) }
func (v *StepView) GetBaseFee() ([32]byte, error)      { return AsRoot(v.Get(fieldBaseFee)) }
func (v *StepView) SetBaseFee(p [32]byte) error        { return v.Set(fieldBaseFee, (*Bytes32View)(&p)) }
func (v *StepView) GetOrigin() (Address, error)        { return AsAddress(v.Get(fieldOrigin)) }
func (v *StepView) SetOrigin(p Address) error          { return v.Set(fieldOrigin, p.View()) }
func (v *StepView) GetTxIndex() (View, error)          { return v.Get(fieldTxIndex) }
func (v *StepView) SetTxIndex(p uint64) error          { return v.Set(fieldTxIndex, Uint64View(p)) }
func (v *StepView) GetGasPrice() (View, error)         { return v.Get(fieldGasPrice) }
func (v *StepView) SetGasPrice(p uint64) error         { return v.Set(fieldGasPrice, Uint64View(p)) }
func (v *StepView) GetTo() (Address, error)            { return AsAddress(v.Get(fieldTo)) }
func (v *StepView) SetTo(p Address) error              { return v.Set(fieldTo, p.View()) }
func (v *StepView) GetCreate() (bool, error)           { return asBool(v.Get(fieldCreate)) }
func (v *StepView) SetCreate(p bool) error             { return v.Set(fieldCreate, BoolView(p)) }
func (v *StepView) GetCallDepth() (View, error)        { return v.Get(fieldCallDepth) }
func (v *StepView) SetCallDepth(p uint64) error        { return v.Set(fieldCallDepth, Uint64View(p)) }
func (v *StepView) GetCaller() (Address, error)        { return AsAddress(v.Get(fieldCaller)) }
func (v *StepView) SetCaller(p Address) error          { return v.Set(fieldCaller, p.View()) }
func (v *StepView) GetMemory() (View, error)           { return v.Get(fieldMemory) }
func (v *StepView) SetMemory(p View) error             { return v.Set(fieldMemory, p) }
func (v *StepView) GetMemLastGasCost() (uint64, error) { return asUint64(v.Get(fieldMemLastGasCost)) }
func (v *StepView) SetMemLastGasCost(p uint64) error   { return v.Set(fieldMemLastGasCost, Uint64View(p)) }
func (v *StepView) GetStack() (View, error)            { return v.Get(fieldStack) }
func (v *StepView) SetStack(p View) error              { return v.Set(fieldStack, p) }
func (v *StepView) GetReturnData() (View, error)       { return v.Get(fieldReturnData) }
func (v *StepView) SetReturnData(p View) error         { return v.Set(fieldReturnData, p) }
func (v *StepView) GetCode() (View, error)             { return v.Get(fieldCode) }
func (v *StepView) SetCode(p View) error               { return v.Set(fieldCode, p) }
func (v *StepView) GetCodeHash() (View, error)         { return v.Get(fieldCodeHash) }
func (v *StepView) SetCodeHash(p [32]byte) error       { return v.Set(fieldCodeHash, (*Bytes32View)(&p)) }
func (v *StepView) GetCodeAddr() (Address, error)      { return AsAddress(v.Get(fieldCodeAddr)) }
func (v *StepView) SetCodeAddr(p Address) error        { return v.Set(fieldCodeAddr, p.View()) }
func (v *StepView) GetInput() (View, error)            { return v.Get(fieldInput) }
func (v *StepView) SetInput(p View) error              { return v.Set(fieldInput, p) }
func (v *StepView) GetGas() (View, error)              { return v.Get(fieldGas) }
func (v *StepView) SetGas(p uint64) error              { return v.Set(fieldGas, Uint64View(p)) }
func (v *StepView) GetValue() ([32]byte, error)        { return AsRoot(v.Get(fieldValue)) }
func (v *StepView) SetValue(p [32]byte) error          { return v.Set(fieldValue, (*Bytes32View)(&p)) }
func (v *StepView) GetOp() (View, error)               { return v.Get(fieldOp) }
func (v *StepView) SetOp(p byte) error                 { return v.Set(fieldOp, ByteView(p)) }
func (v *StepView) GetPc() (uint64, error)             { return asUint64(v.Get(fieldPc)) }
func (v *StepView) SetPc(p uint64) error               { return v.Set(fieldPc, Uint64View(p)) }
func (v *StepView) GetSubIndex() (uint64, error)       { return asUint64(v.Get(fieldSubIndex)) }
func (v *StepView) SetSubIndex(p uint64) error         { return v.Set(fieldSubIndex, Uint64View(p)) }
func (v *StepView) GetSubRemaining() (bool, error)     { return asBool(v.Get(fieldSubRemaining)) }
func (v *StepView) SetSubRemaining(p bool) error       { return v.Set(fieldSubRemaining, BoolView(p)) }
func (v *StepView) GetSubScratch() (View, error)       { return v.Get(fieldSubScratch) }
func (v *StepView) SetSubScratch(p View) error         { return v.Set(fieldSubScratch, p) }
func (v *StepView) GetReturnToStep() (uint64, error)   { return asUint64(v.Get(fieldReturnToStep)) }
func (v *StepView) SetReturnToStep(p uint64) error     { return v.Set(fieldReturnToStep, Uint64View(p)) }

func StepType() *ContainerTypeDef {
	return ContainerType("Step", []FieldDef{
		// Unused in the step itself, but important as output, to claim a state-root,
		// which can then be trusted by the next step.
		// Steps that access memory need to supply a separate (outside of the step sub-tree)
		// MPT-proof of the account and/or storage to access the data.
		{"state_root", Bytes32Type},
		// History scope
		// ------------------
		// Most recent 256 blocks (excluding the block itself)
		{"block_hashes", BlockHashesType},
		// Block scope
		// ------------------
		// TODO: origin balance check for fee payment and value transfer
		{"coinbase", AddressType},
		{"gas_limit", Uint64Type},
		{"block_number", Uint64Type},
		{"time", Uint64Type},
		{"difficulty", Bytes32Type},
		// TODO: maybe as uint256 using geth type?
		{"base_fee", Bytes32Type},
		// Tx scope
		// ------------------
		{"origin", AddressType},
		{"tx_index", Uint64Type},
		{"gas_price", Uint64Type},
		// Contract scope
		// ------------------
		{"to", AddressType},
		{"create", BoolType},
		{"call_depth", Uint64Type},
		{"caller", AddressType},
		{"memory", MemoryType},
		// expanding memory costs exponentially more gas, for the difference in length
		{"memory_last_gas_cost", Uint64Type},
		{"stack", StackType},
		{"return_data", ReturnDataType},
		{"code", CodeType},
		{"code_hash", Bytes32Type},
		{"code_addr", AddressType},
		{"input", InputType},
		{"gas", Uint64Type},
		{"value", Bytes32Type},
		// Execution scope
		// ------------------
		// ignored for starting-state (zeroed)
		{"op", ByteType},
		{"pc", Uint64Type},
		// when splitting up operations further
		{"sub_index", Uint64Type},
		// true when the sub-operation is ongoing and must be completed still.
		{"sub_remaining", BoolType},
		// sub-computations need a place to track their inner state
		{"sub_scratch", ScratchType},
		// When doing a return, continue with the operations after this step.
		{"return_to_step", Uint64Type},
	})
}
