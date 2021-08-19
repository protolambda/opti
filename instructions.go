package opti

import "github.com/ethereum/go-ethereum/common"

func opStart(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	// construct the first (and maybe only sub-step)
	step, err := last.CopyView()
	if err != nil {
		return nil, err
	}

	// History scope: stays the same
	// Block scope: stays the same
	// Tx scope: stays the same

	// Contract scope
	// TODO: update contract scope?

	// Execution scope:
	if err := step.SetSubIndex(0); err != nil {
		return nil, err
	}
	if err := step.SetSubRemaining(false); err != nil {
		return nil, err
	}

	//must(step.Data.SetSubScratch(default))  // TODO: clean sub scratch

	return step, nil
}

func OpStop(trac StepsTrace) (*StepView, error) {
	last := trac.Last()
	next, err := last.CopyView()
	if err != nil {
		return nil, err
	}
	// halt the EVM, we're done with the tx
	if err := next.SetExecMode(ExecSTOP); err != nil {
		return nil, err
	}
	return next, nil
}

func OpPop(trac StepsTrace) (*StepView, error) {
	next, err := opStart(trac)
	if err != nil {
		return nil, err
	}

	stack, err := next.GetStack()
	if err != nil {
		return nil, err
	}
	if err := stack.PopWord(); err != nil {
		return nil, err
	}
	if err := next.IncrementPC(); err != nil {
		return nil, err
	}
	if err := next.SetExecMode(ExecOpcodeLoad); err != nil {
		return nil, err
	}
	return next, nil
}

func MakePush(size uint64, pushByteSize uint64) Processor {
	return func(trac StepsTrace) (*StepView, error) {
		next, err := opStart(trac)
		if err != nil {
			return nil, err
		}

		code, err := next.GetCode()
		if err != nil {
			return nil, err
		}
		codeLen, err := code.Length()
		if err != nil {
			return nil, err
		}

		pc, err := next.GetPc()
		if err != nil {
			return nil, err
		}

		startMin := codeLen
		if pc+1 < startMin {
			startMin = pc + 1
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		data, err := code.Slice(startMin, endMin)
		if err != nil {
			return nil, err
		}
		var w [32]byte
		// TODO: padding is probably wrong
		copy(w[:], common.RightPadBytes(data, 32))

		stack, err := next.GetStack()
		if err != nil {
			return nil, err
		}
		if err := stack.PushWord(w); err != nil {
			return nil, err
		}
		if err := next.SetPc(pc + size); err != nil {
			return nil, err
		}
		if err := next.SetExecMode(ExecOpcodeLoad); err != nil {
			return nil, err
		}
		return next, nil
	}
}
