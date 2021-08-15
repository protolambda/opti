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
