package opti

import (
	"github.com/ethereum/go-ethereum/core/vm"
	. "github.com/protolambda/ztyp/view"
)

const (
	fieldStateRoot = iota
	fieldExecMode
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
	fieldMemLastGas
	fieldMemDesired
	fieldStack
	fieldRetData
	fieldCode
	fieldCodeHash
	fieldCodeAddr
	fieldInput
	fieldGas
	fieldValue
	fieldReadOnly
	fieldOp
	fieldPc
	fieldSubIndex
	fieldSubRemaining
	fieldSubData
	fieldReturnToStep
)

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

func (v *StepView) GetStateRoot() ([32]byte, error)    { return AsRoot(v.Get(fieldStateRoot)) }
func (v *StepView) SetStateRoot(p [32]byte) error      { return v.Set(fieldStateRoot, (*Bytes32View)(&p)) }
func (v *StepView) GetExecMode() (ExecMode, error)     { return asExecMode(v.Get(fieldExecMode)) }
func (v *StepView) SetExecMode(p ExecMode) error       { return v.Set(fieldExecMode, ByteView(p)) }
func (v *StepView) GetBlockHashes() (*HistView, error) { return AsHistView(v.Get(fieldBlockHashes)) }
func (v *StepView) SetBlockHashes(p *HistView) error   { return v.Set(fieldBlockHashes, p) }
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
func (v *StepView) GetTxIndex() (uint64, error)        { return asUint64(v.Get(fieldTxIndex)) }
func (v *StepView) SetTxIndex(p uint64) error          { return v.Set(fieldTxIndex, Uint64View(p)) }
func (v *StepView) GetGasPrice() (uint64, error)       { return asUint64(v.Get(fieldGasPrice)) }
func (v *StepView) SetGasPrice(p uint64) error         { return v.Set(fieldGasPrice, Uint64View(p)) }
func (v *StepView) GetTo() (Address, error)            { return AsAddress(v.Get(fieldTo)) }
func (v *StepView) SetTo(p Address) error              { return v.Set(fieldTo, p.View()) }
func (v *StepView) GetCreate() (bool, error)           { return asBool(v.Get(fieldCreate)) }
func (v *StepView) SetCreate(p bool) error             { return v.Set(fieldCreate, BoolView(p)) }
func (v *StepView) GetCallDepth() (uint64, error)      { return asUint64(v.Get(fieldCallDepth)) }
func (v *StepView) SetCallDepth(p uint64) error        { return v.Set(fieldCallDepth, Uint64View(p)) }
func (v *StepView) GetCaller() (Address, error)        { return AsAddress(v.Get(fieldCaller)) }
func (v *StepView) SetCaller(p Address) error          { return v.Set(fieldCaller, p.View()) }
func (v *StepView) GetMemory() (*MemoryView, error)    { return AsMemoryView(v.Get(fieldMemory)) }
func (v *StepView) SetMemory(p *MemoryView) error      { return v.Set(fieldMemory, p) }
func (v *StepView) GetMemLastGas() (uint64, error)     { return asUint64(v.Get(fieldMemLastGas)) }
func (v *StepView) SetMemLastGas(p uint64) error       { return v.Set(fieldMemLastGas, Uint64View(p)) }
func (v *StepView) GetMemDesired() (uint64, error)     { return asUint64(v.Get(fieldMemDesired)) }
func (v *StepView) SetMemDesired(p uint64) error       { return v.Set(fieldMemDesired, Uint64View(p)) }
func (v *StepView) GetStack() (*StackView, error)      { return AsStackView(v.Get(fieldStack)) }
func (v *StepView) SetStack(p *StackView) error        { return v.Set(fieldStack, p) }
func (v *StepView) GetRetData() (*RetDataView, error)  { return AsRetDataView(v.Get(fieldRetData)) }
func (v *StepView) SetRetData(p *RetDataView) error    { return v.Set(fieldRetData, p) }
func (v *StepView) GetCode() (*CodeView, error)        { return AsCodeView(v.Get(fieldCode)) }
func (v *StepView) SetCode(p *CodeView) error          { return v.Set(fieldCode, p) }
func (v *StepView) GetCodeHash() ([32]byte, error)     { return AsRoot(v.Get(fieldCodeHash)) }
func (v *StepView) SetCodeHash(p [32]byte) error       { return v.Set(fieldCodeHash, (*Bytes32View)(&p)) }
func (v *StepView) GetCodeAddr() (Address, error)      { return AsAddress(v.Get(fieldCodeAddr)) }
func (v *StepView) SetCodeAddr(p Address) error        { return v.Set(fieldCodeAddr, p.View()) }
func (v *StepView) GetInput() (*InputView, error)      { return AsInputView(v.Get(fieldInput)) }
func (v *StepView) SetInput(p *InputView) error        { return v.Set(fieldInput, p) }
func (v *StepView) GetGas() (uint64, error)            { return asUint64(v.Get(fieldGas)) }
func (v *StepView) SetGas(p uint64) error              { return v.Set(fieldGas, Uint64View(p)) }
func (v *StepView) GetValue() ([32]byte, error)        { return AsRoot(v.Get(fieldValue)) }
func (v *StepView) SetValue(p [32]byte) error          { return v.Set(fieldValue, (*Bytes32View)(&p)) }
func (v *StepView) GetReadOnly() (bool, error)         { return asBool(v.Get(fieldReadOnly)) }
func (v *StepView) SetReadOnly(p bool) error           { return v.Set(fieldReadOnly, BoolView(p)) }
func (v *StepView) GetOp() (vm.OpCode, error)          { return asOpCode(v.Get(fieldOp)) }
func (v *StepView) SetOp(p vm.OpCode) error            { return v.Set(fieldOp, ByteView(p)) }
func (v *StepView) GetPc() (uint64, error)             { return asUint64(v.Get(fieldPc)) }
func (v *StepView) SetPc(p uint64) error               { return v.Set(fieldPc, Uint64View(p)) }
func (v *StepView) GetSubIndex() (uint64, error)       { return asUint64(v.Get(fieldSubIndex)) }
func (v *StepView) SetSubIndex(p uint64) error         { return v.Set(fieldSubIndex, Uint64View(p)) }
func (v *StepView) GetSubRemaining() (bool, error)     { return asBool(v.Get(fieldSubRemaining)) }
func (v *StepView) SetSubRemaining(p bool) error       { return v.Set(fieldSubRemaining, BoolView(p)) }
func (v *StepView) GetSubData() (*SubDataView, error)  { return AsSubDataView(v.Get(fieldSubData)) }
func (v *StepView) SetSubData(p *SubDataView) error    { return v.Set(fieldSubData, p) }
func (v *StepView) GetReturnToStep() (uint64, error)   { return asUint64(v.Get(fieldReturnToStep)) }
func (v *StepView) SetReturnToStep(p uint64) error     { return v.Set(fieldReturnToStep, Uint64View(p)) }

func StepType() *ContainerTypeDef {
	return ContainerType("Step", []FieldDef{
		// Unused in the step itself, but important as output, to claim a state-root,
		// which can then be trusted by the next step.
		// Steps that access memory need to supply a separate (outside of the step sub-tree)
		// MPT-proof of the account and/or storage to access the data.
		{"state_root", Bytes32Type},

		{"exec_mode", Uint8Type},

		// TODO: operation-mode field, to switch between:
		// - block processing
		// - transaction inclusion check
		// - transaction validation
		// - interpreter loop: "Run"
		// - opcode processing
		// - block post-processing (build MPT from buffered receipt roots, handle logs-bloom, etc.)

		// History scope
		// ------------------
		// Most recent 256 blocks (excluding the block itself)
		{"block_hashes", HistType},
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
		{"memory_last_gas", Uint64Type},
		// We compute the memory size, charge for it first, and only then allocate it.
		{"memory_desired", Uint64Type},
		{"stack", StackType},
		{"ret_data", RetDataType},
		{"code", CodeType},
		{"code_hash", Bytes32Type},
		{"code_addr", AddressType},
		{"input", InputType},
		{"gas", Uint64Type},
		{"value", Bytes32Type},
		// Make storage read-only, to support STATIC-CALL
		{"read_only", BoolType},
		// Execution scope
		// ------------------
		// We generalize the opcode read from the code at PC,
		// and cache it here to not re-read the code every step of the opcode.
		// Ignored for starting-state (zeroed).
		{"op", ByteType},
		// The program-counter, index pointing to current executed opcode in the code
		{"pc", Uint64Type},
		// when splitting up operations further
		{"sub_index", Uint64Type},
		// true when the sub-operation is ongoing and must be completed still.
		{"sub_remaining", BoolType},
		// sub-computations need a place to track their inner state
		{"sub_data", SubDataType},
		// When doing a return, continue with the operations after this step.
		{"return_to_step", Uint64Type},
		// Depending on the unwind mode: continue/return/out-of-gas/etc.
		{"unwind", Uint64Type},
	})
}
