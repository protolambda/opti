package opti

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
	"time"
)

type Step struct {
	Data *StepView

	// Used in step execution, not merkleized into step itself,
	//  must match state-root of previous step (unchecked if nil).
	// If more than 1 storage slot is touched, the step must be split up in sub-steps.

	// nil if unused
	AccountProof [][]byte
	// nil if unused, requires account-proof
	StorageProof [][]byte
}

func (step *Step) Copy() (*Step, error) {
	v, err := step.Data.CopyView()
	if err != nil {
		return nil, fmt.Errorf("failed to copy step data: %v", err)
	}
	return &Step{Data: v}, nil
}

func (step *Step) MultiProof() [][]common.Hash {
	// TODO: from the step, represented as binary tree, construct the proof
	return nil
}

func (step *Step) Root() common.Hash {
	// TODO Compute step root, should match root corresponding to proof of data within the step
	// Exclude StepIndices field, and state-snapshot integer
	return common.Hash{}
}

type OptiTracer struct {
	Steps []*Step
}

type VerifierStateDB interface {
	vm.StateDB
	GetProof(addr common.Address) ([][]byte, error)
	GetProofByHash(addrHash common.Hash) ([][]byte, error)
	GetStorageProof(a common.Address, key common.Hash) ([][]byte, error)
	IntermediateRoot(deleteEmptyObjects bool) common.Hash
}

func (o *OptiTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	// TODO: maybe move this to the initialization of the tracer or a separate function, to be called for each transaction before tracing said transaction?

	// TODO init first step; zero stack, depth, memory, initial balance, tx signature verification, etc. (can be sub-steps for fraud-proof simplicity)
	initialStep := &Step{
		// TODO: init first step, containing pre-state root, block data, tx data, value, etc.
	}
	o.Steps = append(o.Steps, initialStep)

	// TODO: add sub-steps to verify the transaction
}

func (o *OptiTracer) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	prev := o.Steps[len(o.Steps)-1]

	// TODO: maybe log errors for debugging? Generally all tree modifications can be "must", since the tree representation is in memory, so no IO errors
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// construct the first (and maybe only sub-step)
	// TODO: instead of copying the full step, represent it as binary tree of data with typed view overlay, and copy while data-sharing
	step, err := prev.Copy()
	must(err)

	// History scope: stays the same
	// Block scope: stays the same
	// Tx scope: stays the same

	// Contract scope
	// TODO: update contract scope?

	// Execution scope:
	must(step.Data.SetPc(pc))
	must(step.Data.SetSubIndex(0))
	must(step.Data.SetSubRemaining(false))
	//must(step.Data.SetSubScratch(default))  // TODO: clean sub scratch

	// TODO add step(s)
	// - a basic opcode may just be 1 step (cheat code: we can diff the steps, and not re-implement the EVM opcode behavior, to produce the new step)
	// - if modifying the stack, memory, etc. we can add a sub-step that checks the bounds / cost / etc.
	// - if an opcode is complex, use sub-steps to split up the computation, e.g.:
	//   - a large copy in many small copies
	//   - a hash with large input into smaller input-processing and final digest computation
	stateDB, ok := env.StateDB.(VerifierStateDB)
	if !ok {
		panic("EVM state db lacks proving and state-root abilities for verifier")
	}
	// Given the step output of the previous (sub)step, a smart-contract functioning as EVM-step-emulator should be able to produce the next output

	step.Data.SetStateRoot(stateDB.IntermediateRoot(false))
	// TODO: maybe more sub-steps, depending on opcode
	step.Data.SetSubRemaining(false)
	// add the step
	o.Steps = append(o.Steps, step)
}

func (o *OptiTracer) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
	// TODO this is only traced during debugging mode, and should overlap with the output of the CaptureState call
}

func (o *OptiTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	// add final accounting steps, and declare a special operation step that checks if the state-root matches the claimed one

	// TODO: transaction bookkeeping is not done within the EVM, but here we merge the EVM and blockchain runtime, to process batches of transactions

	// TODO: when depth==0 and tx ends, gas_used=gas_limit-gas_remaining, charge origin
}

var _ vm.Tracer = (*OptiTracer)(nil)
