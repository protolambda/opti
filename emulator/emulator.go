package emulator

import (
	"fmt"
	"github.com/protolambda/opti"
	"github.com/protolambda/ztyp/tree"
)

type IndexedStep struct {
	Index uint64
	Step  *opti.StepView
}

func (p *IndexedStep) String() string {
	if p == nil {
		return "<nil>"
	}
	stepRoot := p.Step.HashTreeRoot(tree.GetHashFn())
	return fmt.Sprintf("%d (%x)", p.Index, stepRoot[:7])
}

type partialTrace struct {
	Ref *IndexedStep
	Pre *IndexedStep
}

func (p *partialTrace) Last() *opti.StepView {
	return p.Pre.Step
}

func (p *partialTrace) ByIndex(i uint64) *opti.StepView {
	if i == p.Ref.Index {
		return p.Ref.Step
	}
	if i == p.Pre.Index {
		return p.Pre.Step
	}
	panic(fmt.Errorf("invalid witness, cannot access step %d (executing step %s, optional step %s)", i, p.Pre, p.Ref))
}

func (p *partialTrace) Length() uint64 {
	return p.Pre.Index + 1
}

var _ opti.StepsTrace = (*partialTrace)(nil)

// This emulator is not part of the fraud proof generation code, but rather consumes it, to verify the fraud proof is correct.
// This emulator needs to be translated to solidity to make the fraud-proof run on L1
//
// Warning: this verifies if the transition after a pre-step, along with an optional referenced step, results in a step with a matching merkle root.
// It does not check if the first step, out of all, is correct, or if the last one matches the state-root as claimed at the end of the block processing
// (TODO: might actually be able to put state-root comparison into it)
func Verify(conf *opti.Config, pre *IndexedStep, ref *IndexedStep, expectedPostStepRoot [32]byte) (ok bool, err error) {
	partial := &partialTrace{Ref: ref, Pre: pre}
	out, err := conf.NextStep(partial)
	if err != nil {
		return false, fmt.Errorf("failed to produce next step: %v", err)
	}
	outRoot := out.HashTreeRoot(tree.GetHashFn())
	return outRoot == expectedPostStepRoot, nil
}
