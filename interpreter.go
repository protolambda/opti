package opti

type Operation struct {
	Proc Processor
	// TODO: properties and gas costs etc.
}

type Config struct {
	JumpTable [256]*Operation
}

func (cfg *Config) Run(trac StepsTrace) (err error) {

	// TODO set step mode to "interpreter"

	// Increment the call depth which is restricted to 1024
	{
		prev := trac.Prev()
		step, err := prev.CopyView()
		if err != nil {
			return err
		}
		callDepth, err := step.GetCallDepth()
		if err != nil {
			return err
		}
		if err := step.SetCallDepth(callDepth + 1); err != nil {
			return err
		}
		if err := trac.Append(step); err != nil {
			return err
		}
	}

	// Decrement the call depth when we are done with the execution
	defer func() {
		// if we exit gracefully, unwind the callstack
		// TODO: geth uses errors for EVM execution flow, in our case those would be steps, we return errors for non-input problems.
		if err != nil {
			return
		}

		err = func() error {
			prev := trac.Prev()
			step, err := prev.CopyView()
			if err != nil {
				return err
			}
			callDepth, err := step.GetCallDepth()
			if err != nil {
				return err
			}
			if err := step.SetCallDepth(callDepth - 1); err != nil {
				return err
			}
			return nil
		}()
	}()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This also makes sure that the readOnly flag isn't removed for child calls.
	// TODO

	// Geth: Reset the previous call's return data. It's unimportant to preserve the old buffer
	//       as every returning call will return new data anyway.
	// Opti: Ok to overwrite, can reset back from previous step by index.
	{
		prev := trac.Prev()
		step, err := prev.CopyView()
		if err != nil {
			return err
		}
		if err := step.SetRetData(EmptyReturnData()); err != nil {
			return err
		}
		if err := trac.Append(step); err != nil {
			return err
		}
	}

	// TODO: check code size == 0, exit early

	// TODO: init new call context

	// TODO: defer stack reset

	// TODO: switch to pre-operation mode

	// The Interpreter main run loop (contextual). This loop runs until either an
	// explicit STOP, RETURN or SELFDESTRUCT is executed, an error occurred during
	// the execution of one of the operations or until the done flag is set by the
	// parent context.
	for {
		prev := trac.Prev()
		pc, err := prev.GetPc()
		if err != nil {
			return err
		}
		code, err := prev.GetCode()
		if err != nil {
			return err
		}
		op, err := code.GetOp(pc)
		if err != nil {
			return err
		}
		// Get the operation from the jump table and validate the stack to ensure there are
		// enough stack items available to perform the operation.
		operation := cfg.JumpTable[op]

		// validate stack
		// TODO

		// If the operation is valid, enforce write restrictions
		// TODO

		// Static portion of gas
		// TODO

		// calculate the new memory size and expand the memory to fit
		// the operation
		// Memory check needs to be done prior to evaluating the dynamic gas portion,
		// to detect calculation overflows
		// TODO

		// Dynamic portion of gas
		// consume the gas and return an error if not enough gas is available.
		// cost is explicitly set so that the capture state defer method can get the proper cost
		// TODO

		// TODO: resize memory

		// TODO: switch step to operation mode

		// TODO: execute operation
		if err := operation.Proc(trac); err != nil {
			return err
		}

		// Operation processor should already have placed new return data in the last step if it outputs anything
		// Thus no need to copy

		// TODO: switch step to post-operation mode
	}

	return
}
