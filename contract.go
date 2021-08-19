package opti

func (step *StepView) UseGas(delta uint64) (ok bool, err error) {
	gas, err := step.GetGas()
	if err != nil {
		return false, err
	}
	if delta > gas {
		// out of gas
		return false, nil
	}
	return true, step.SetGas(gas - delta)
}

func (step *StepView) ReturnGas(delta uint64) error {
	gas, err := step.GetGas()
	if err != nil {
		return err
	}
	// no overflow, assuming gas total is capped within uint64
	return step.SetGas(gas + delta)
}

func (step *StepView) IncrementPC() error {
	pc, err := step.GetPc()
	if err != nil {
		return err
	}
	return step.SetPc(pc + 1)
}
