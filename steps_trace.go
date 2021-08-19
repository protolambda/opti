package opti

import "fmt"

type Steps []*StepView

func (steps *Steps) Last() *StepView {
	if steps == nil {
		return nil
	}
	if len(*steps) == 0 {
		return nil
	}
	return (*steps)[len(*steps)-1]
}

func (steps *Steps) ByIndex(i uint64) *StepView {
	if steps == nil {
		return nil
	}
	if i >= uint64(len(*steps)) {
		return nil
	}
	return (*steps)[i]
}

var _ StepsTrace = (*Steps)(nil)

func (steps *Steps) Append(step *StepView) error {
	if steps == nil {
		return fmt.Errorf("cannot extend nil steps")
	}
	*steps = append(*steps, step)
	return nil
}

func (steps *Steps) Length() uint64 {
	if steps == nil {
		return 0
	}
	return uint64(len(*steps))
}

type StepsTrace interface {
	// nil if empty steps
	Last() *StepView
	// nil if invalid index
	ByIndex(i uint64) *StepView
	Length() uint64
}

type WritableStepsTrace interface {
	StepsTrace
	Append(step *StepView) error
}

// Takes a trace, and produces the next step
type Processor func(trac StepsTrace) (*StepView, error)
