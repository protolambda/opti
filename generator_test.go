package opti

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/protolambda/ztyp/tree"
	"testing"
)

func TestGenerateProof(t *testing.T) {
	// not many opcodes supported yet, try something simple to test the interpreter state machine
	code := []byte{byte(vm.PUSH1), byte(vm.POP), byte(vm.PUSH1)}
	value := [32]byte{31: 123}
	from := Address{0xaa}
	to := Address{0xbb}
	gas := uint64(100)
	executionTrace, err := GenerateProof(code, value, from, to, gas)
	if err != nil {
		t.Fatal(err)
	}
	last := executionTrace.Last()
	stack, err := last.GetStack()
	if err != nil {
		t.Fatal(err)
	}
	l, err := stack.Length()
	if err != nil {
		t.Fatal(err)
	}
	if l != 1 {
		t.Fatal("expected 1 item on the stack")
	}
	remainingGas, err := last.GetGas()
	if err != nil {
		t.Fatal(err)
	}
	// started with 100 gas
	// push1 == 3 gas
	// pop == 2 gas
	if remainingGas != 100-3-2-3 {
		t.Fatalf("expected different gas value than: %d", remainingGas)
	}
	for i, step := range *executionTrace.(*Steps) {
		stepRoot := step.HashTreeRoot(tree.GetHashFn())
		fmt.Printf("step %d merkleizes to %s\n", i, stepRoot)
		// TODO: extract witness, we need to use *Step instead of the *StepView.
		// This wraps the view, maintains state data proofs as well,
		// and we can inject a virtual wrapper in the step tree backing to capture
		// touched merkle branches.
		// Too verbose, too much work, switching prototype to python for experimentation first.
	}
}
