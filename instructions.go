package opti

import "github.com/ethereum/go-ethereum/common"

func opStart(prev *StepView) (*StepView, error) {
	// construct the first (and maybe only sub-step)
	step, err := prev.CopyView()
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

func OpPop(trac StepsTrace) error {
	step, err := opStart(trac.Prev())
	if err != nil {
		return err
	}

	stack, err := step.GetStack()
	if err != nil {
		return err
	}
	if err := stack.PopWord(); err != nil {
		return err
	}
	return trac.Append(step)
}

func MakePush(size uint64, pushByteSize uint64) Processor {
	return func(trac StepsTrace) error {
		step, err := opStart(trac.Prev())
		if err != nil {
			return err
		}

		code, err := step.GetCode()
		if err != nil {
			return err
		}
		codeLen, err := code.Length()
		if err != nil {
			return err
		}

		pc, err := step.GetPc()
		if err != nil {
			return err
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
			return err
		}
		var w [32]byte
		// TODO: padding is probably wrong
		copy(w[:], common.RightPadBytes(data, 32))

		stack, err := step.GetStack()
		if err != nil {
			return err
		}
		if err := stack.PushWord(w); err != nil {
			return err
		}
		if err := step.SetPc(pc + size); err != nil {
			return err
		}
		return trac.Append(step)
	}
}
