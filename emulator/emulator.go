package emulator

import "github.com/protolambda/opti"

// This emulator is not part of the fraud proof generation code, but rather consumes it, to verify the fraud proof is correct.
// This emulator needs to be translated to solidity to make the fraud-proof run on L1
func Verify(pre *opti.StepView, postStepRoot [32]byte) (ok bool, err error) {
	// TODO include dependent step data and MPT proofs
	return false, nil
}
