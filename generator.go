package opti

import (
	"bytes"
	"fmt"
	"github.com/protolambda/ztyp/codec"
)

func GenerateProof(code []byte, value [32]byte, from Address, to Address, gas uint64) (StepsTrace, error) {
	// TODO: do proper tx initialization, instead of hacking basics into an initial step
	step, err := AsStep(StepType().Default(nil), nil)
	if err != nil {
		return nil, err
	}
	codeView, err := AsCodeView(CodeType.Deserialize(codec.NewDecodingReader(bytes.NewReader(code), uint64(len(code)))))
	if err != nil {
		return nil, err
	}
	if err := step.SetCode(codeView); err != nil {
		return nil, err
	}
	if err := step.SetValue(value); err != nil {
		return nil, err
	}
	if err := step.SetCaller(from); err != nil {
		return nil, err
	}
	if err := step.SetTo(to); err != nil {
		return nil, err
	}
	if err := step.SetGas(gas); err != nil {
		return nil, err
	}
	// hack in an initial call frame
	if err := step.SetExecMode(ExecCallPre); err != nil {
		return nil, err
	}
	// Now wrap the initial step in a trace
	trac := &Steps{step}

	// Create an EVM config
	conf := Config{JumpTable: DefaultJumpTable()}

	// And start extending trace until we have all steps
	for i := 0; true; i++ {
		fmt.Printf("processing step %d now\n", i)
		// check if we're done evaluating
		last := trac.Last()
		execMode, err := last.GetExecMode()
		if err != nil {
			return nil, err
		}
		// check if done
		if execMode == ExecSTOP {
			break
		}

		next, err := conf.NextStep(trac)
		if err != nil {
			return nil, fmt.Errorf("failed to extend trace at index %d: %v", i, err)
		}
		if err := trac.Append(next); err != nil {
			return nil, err
		}
	}

	// trace is complete!
	return trac, nil
}
