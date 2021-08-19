package opti

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
)

type (
	// produce a step that consumes the dynamic gas for the current opcode, by calling contract.UseGas(dynamicCost)
	// Them either produce an error execution mode (like ExecErrOutOfGas),
	// or continue the interpreter by setting the execution mode to ExecUpdateMemorySize
	gasFunc Processor

	// memorySizeFunc returns the required size, and whether the operation overflowed a uint64
	memorySizeFunc func(view *StackView) (size uint64, overflow bool, err error)
)

type Operation struct {
	Proc Processor

	constantGas uint64
	dynamicGas  gasFunc
	// minStack tells how many stack items are required
	minStack uint64
	// maxStack specifies the max length the stack can have for this operation
	// to not overflow the stack.
	maxStack uint64

	// memorySize returns the memory size required for the operation
	memorySize memorySizeFunc

	halts   bool // indicates whether the operation should halt further execution
	jumps   bool // indicates whether the program counter should not increment
	writes  bool // determines whether this a state modifying operation
	reverts bool // determines whether the operation reverts state (implicitly halts)
	returns bool // determines whether the operations sets the return data content
}

// The interpreter is a pure function: given the last step, and optional access to a previous step, it can compute the next step.
// The execution mode distinguishes the various tasks around just regular opcode execution.
type ExecMode byte

const (
	ExecBlockPre = iota

	// TODO: check transaction is included, check fee params, signature validity, etc.
	ExecTxInclusion
	ExecTxSig
	ExecTxFeesPre

	// Call pre-processing (increment call depth)
	ExecCallPre

	// Interpreter loop consists of stack/memory/gas checks, and then opcode execution.
	ExecOpcodeLoad
	ExecValidateStack
	ExecReadOnlyCheck
	ExecConstantGas
	ExecCalcMemorySize
	ExecDynamicGas
	ExecUpdateMemorySize
	// when done running, continue with ExecOpcodeLoad. Or any error
	ExecOpcodeRun

	// Any error should follow up with running call-post processing
	ExecCallPost

	// TODO: charge tx fees
	ExecTxFeesPost

	// TODO: maybe enumerate errors separately, generalize them?
	ExecErrStackUnderflow
	ExecErrStackOverflow
	ExecErrWriteProtection
	ExecErrOutOfGas
	ExecErrGasUintOverflow
)

type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
	IsBerlin, IsLondon                                      bool
}

type Config struct {
	JumpTable [256]*Operation
}

func (cfg *Config) Rules(blockNum uint64) Rules {
	// TODO: adjust rules based on chain config like geth. Maybe just reuse their config code?
	return Rules{
		ChainID:          new(big.Int).SetUint64(12345), // TODO: chain ID
		IsHomestead:      true,
		IsEIP150:         true,
		IsEIP155:         true,
		IsEIP158:         true,
		IsByzantium:      true,
		IsConstantinople: true,
		IsPetersburg:     true,
		IsIstanbul:       true,
		IsBerlin:         true,
		IsLondon:         true,
	}
}

func (cfg *Config) NextStep(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	mode, err := last.GetExecMode()
	if err != nil {
		return nil, err
	}
	switch mode {
	case ExecBlockPre:
		// TODO
		return nil, errors.New("fix me")
	case ExecTxInclusion:
		// TODO
		return nil, errors.New("fix me")
	case ExecTxSig:
		// TODO
		return nil, errors.New("fix me")
	case ExecTxFeesPre:
		// TODO
		return nil, errors.New("fix me")
	// TODO: tx init step, to initialize contract scope
	case ExecCallPre:
		return cfg.execCallPre(trac)
	case ExecOpcodeLoad:
		return cfg.execOpcodeLoad(trac)
	case ExecValidateStack:
		return cfg.execValidateStack(trac)
	case ExecReadOnlyCheck:
		return cfg.execReadOnlyCheck(trac)
	case ExecConstantGas:
		return cfg.execConstantGas(trac)
	case ExecCalcMemorySize:
		return cfg.execCalcMemorySize(trac)
	case ExecDynamicGas:
		return cfg.execDynamicGas(trac)
	case ExecUpdateMemorySize:
		return cfg.execUpdateMemorySize(trac)
	case ExecOpcodeRun:
		return cfg.execOpcodeRun(trac)
	case ExecCallPost:
		return cfg.execCallPost(trac)
	case ExecTxFeesPost:
		// TODO
		return nil, errors.New("fix me")
	default:
		return nil, fmt.Errorf("unrecognized execution mode: %v", mode)
	}
}

func (cfg *Config) execCallPre(trac StepsTrace) (*StepView, error) {
	// Increment the call depth which is restricted to 1024
	last := trac.Last()

	callDepth, err := last.GetCallDepth()
	if err != nil {
		return nil, err
	}

	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}
	if err := next.SetCallDepth(callDepth + 1); err != nil {
		return nil, err
	}

	// TODO readonly mode support

	// reset returndata
	if err := next.SetRetData(EmptyReturnData()); err != nil {
		return nil, err
	}

	// TODO check code size, go directly to execCallPost if 0.

	// TODO: reset stack
	// TODO: reset memory
	// TODO: reset PC

	// Move to the first opcode of the code in this call
	if err := next.SetExecMode(ExecOpcodeLoad); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execOpcodeLoad(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	pc, err := last.GetPc()
	if err != nil {
		return nil, err
	}
	code, err := last.GetCode()
	if err != nil {
		return nil, err
	}
	// 0 if PC is too high, same hack as in Geth
	op, err := code.GetOp(pc)
	if err != nil {
		return nil, err
	}
	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}
	if err := next.SetOp(op); err != nil {
		return nil, err
	}
	if err := next.SetExecMode(ExecValidateStack); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execValidateStack(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}

	// validate stack
	stack, err := last.GetStack()
	if err != nil {
		return nil, err
	}
	sLen, err := stack.Length()
	if err != nil {
		return nil, err
	}

	op, err := last.GetOp()
	if err != nil {
		return nil, err
	}

	// Get the operation from the jump table and validate the stack to ensure there are
	// enough stack items available to perform the operation.
	operation := cfg.JumpTable[op]

	if sLen < operation.minStack {
		if err := next.SetExecMode(ExecErrStackUnderflow); err != nil {
			return nil, err
		}
		return next, nil
	} else if sLen > operation.maxStack {
		if err := next.SetExecMode(ExecErrStackOverflow); err != nil {
			return nil, err
		}
		return next, nil
	}

	// valid stack, continue interpreter
	if err := next.SetExecMode(ExecReadOnlyCheck); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execReadOnlyCheck(trac StepsTrace) (*StepView, error) {
	last := trac.Last()

	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}

	isReadOnly, err := last.GetReadOnly()
	if err != nil {
		return nil, err
	}

	// If the operation is valid, enforce write restrictions
	if isReadOnly {
		op, err := last.GetOp()
		if err != nil {
			return nil, err
		}
		operation := cfg.JumpTable[op]
		// If the interpreter is operating in readonly mode, make sure no
		// state-modifying operation is performed. The 3rd stack item
		// for a call operation is the value. Transferring value from one
		// account to the others means the state is modified and should also
		// return with an error.
		if operation.writes {
			if err := next.SetExecMode(ExecErrWriteProtection); err != nil {
				return nil, err
			}
			return next, nil
		}
		if op == vm.CALL {
			stack, err := last.GetStack()
			if err != nil {
				return nil, err
			}
			value, err := stack.Back(2)
			if err != nil {
				return nil, err
			}
			if value != ([32]byte{}) {
				if err := next.SetExecMode(ExecErrWriteProtection); err != nil {
					return nil, err
				}
				return next, nil
			}
		}
	}

	// Continue interpreter to next step
	if err := next.SetExecMode(ExecConstantGas); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execConstantGas(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	op, err := last.GetOp()
	if err != nil {
		return nil, err
	}
	operation := cfg.JumpTable[op]

	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}
	// Static portion of gas
	if ok, err := next.UseGas(operation.constantGas); err != nil {
		return nil, err
	} else if !ok {
		if err := next.SetExecMode(ExecErrOutOfGas); err != nil {
			return nil, err
		}
		return next, nil
	}

	// Continue interpreter to next step
	if err := next.SetExecMode(ExecCalcMemorySize); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execCalcMemorySize(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	op, err := last.GetOp()
	if err != nil {
		return nil, err
	}
	operation := cfg.JumpTable[op]

	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}

	var memorySize uint64
	// calculate the new memory size and expand the memory to fit
	// the operation
	// Memory check needs to be done prior to evaluating the dynamic gas portion,
	// to detect calculation overflows
	if operation.memorySize != nil {
		stack, err := last.GetStack()
		if err != nil {
			return nil, err
		}
		memSize, overflow, err := operation.memorySize(stack)
		if err != nil {
			return nil, err
		}
		if overflow {
			if err := next.SetExecMode(ExecErrGasUintOverflow); err != nil {
				return nil, err
			}
			return next, nil
		}
		// memory is expanded in words of 32 bytes. Gas
		// is also calculated in words.
		if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
			if err := next.SetExecMode(ExecErrGasUintOverflow); err != nil {
				return nil, err
			}
			return next, nil
		}
	}

	if err := next.SetMemDesired(memorySize); err != nil {
		return nil, err
	}

	if err := next.SetExecMode(ExecDynamicGas); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execDynamicGas(trac StepsTrace) (*StepView, error) {
	last := trac.Last()

	op, err := last.GetOp()
	if err != nil {
		return nil, err
	}
	operation := cfg.JumpTable[op]

	if operation.dynamicGas != nil {
		return operation.dynamicGas(trac)
	}

	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}
	if err := next.SetExecMode(ExecUpdateMemorySize); err != nil {
		return nil, err
	}
	return next, nil
}

func (cfg *Config) execUpdateMemorySize(trac StepsTrace) (*StepView, error) {
	last := trac.Last()

	memorySize, err := last.GetMemDesired()
	if err != nil {
		return nil, err
	}

	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}

	if memorySize > 0 {
		memory, err := next.GetMemory()
		if err != nil {
			return nil, err
		}

		mLen, err := memory.Length()
		if err != nil {
			return nil, err
		}

		if mLen >= memorySize {
			// no extension to do
			if err := next.SetExecMode(ExecOpcodeRun); err != nil {
				return nil, err
			}
			return next, nil
		}

		// If we can efficiently add 32 bytes of memory, go for it.
		// Otherwise just align by adding a byte, or finishing trailing non-aligning bytes.
		if (memorySize-mLen > 32) && (mLen%32 == 0) {
			memorySize -= 32
			if err := memory.AppendZeroBytes32(); err != nil {
				return nil, err
			}
		} else {
			memorySize -= 1
			if err := memory.AppendZeroByte(); err != nil {
				return nil, err
			}
		}
		if err := next.SetMemDesired(memorySize); err != nil {
			return nil, err
		}

		if mLen >= memorySize {
			// no extension to do
			if err := next.SetExecMode(ExecOpcodeRun); err != nil {
				return nil, err
			}
			return next, nil
		} else {
			// leave the execution mode on ExecUpdateMemorySize. The next step will further extend the memory size.
			return next, nil
		}
	} else {
		if err := next.SetExecMode(ExecOpcodeRun); err != nil {
			return nil, err
		}
		return next, nil
	}
}

func (cfg *Config) execOpcodeRun(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	op, err := last.GetOp()
	if err != nil {
		return nil, err
	}
	operation := cfg.JumpTable[op]
	return operation.Proc(trac)
}

func (cfg *Config) execCallPost(trac StepsTrace) (*StepView, error) {
	last := trac.Last()

	// note; call-depth will unwind itself, since we copy it from the caller.

	retStepIndex, err := last.GetReturnToStep()
	if err != nil {
		return nil, err
	}
	callerStep := trac.ByIndex(retStepIndex)
	if callerStep == nil {
		return nil, fmt.Errorf("unknown return-step: %d", retStepIndex)
	}

	// Next step is a lot like the caller, but we preserve the return data, return unused gas, and preserve the state
	// TODO maybe also track past log events, to reconstruct receipt root for block fraud proof
	next, err := callerStep.CopyView()
	if err != nil {
		return nil, err
	}

	stateRoot, err := last.GetStateRoot()
	if err != nil {
		return nil, err
	}
	if err := next.SetStateRoot(stateRoot); err != nil {
		return nil, err
	}

	retData, err := last.GetRetData()
	if err != nil {
		return nil, err
	}
	if err := next.SetRetData(retData); err != nil {
		return nil, err
	}

	unusedGas, err := last.GetGas()
	if err != nil {
		return nil, err
	}
	if err := next.ReturnGas(unusedGas); err != nil {
		return nil, err
	}

	// Continue processing in caller, exactly where we left of,
	// but indicating the return by incrementing the sub-index.
	//
	// I.e. we don't want to keep repeating the same call forever, but start post-processing of this call as caller.
	subIndex, err := next.GetSubIndex()
	if err != nil {
		return nil, err
	}
	if err := next.SetSubIndex(subIndex + 1); err != nil {
		return nil, err
	}

	return next, nil
}
