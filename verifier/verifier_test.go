package verifier_test

import (
	"testing"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"github.com/succinctlabs/gnark-plonky2-verifier/field"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier"
	"github.com/succinctlabs/gnark-plonky2-verifier/verifier/utils"
)

type TestVerifierCircuit struct {
	proofWithPIsFilename            string `gnark:"-"`
	commonCircuitDataFilename       string `gnark:"-"`
	verifierOnlyCircuitDataFilename string `gnark:"-"`
}

func (c *TestVerifierCircuit) Define(api frontend.API) error {
	proofWithPis := utils.DeserializeProofWithPublicInputs(c.proofWithPIsFilename)
	commonCircuitData := utils.DeserializeCommonCircuitData(c.commonCircuitDataFilename)
	verfierOnlyCircuitData := utils.DeserializeVerifierOnlyCircuitData(c.verifierOnlyCircuitDataFilename)

	verifierChip := verifier.NewVerifierChip(api, commonCircuitData)
	verifierChip.Verify(proofWithPis, verfierOnlyCircuitData, commonCircuitData)
	return nil
}

func TestDecodeBlockVerifier(t *testing.T) {
	assert := test.NewAssert(t)

	testCase := func() {
		circuit := TestVerifierCircuit{
			proofWithPIsFilename:            "./data/decode_block/proof_with_public_inputs.json",
			commonCircuitDataFilename:       "./data/decode_block/common_circuit_data.json",
			verifierOnlyCircuitDataFilename: "./data/decode_block/verifier_only_circuit_data.json",
		}

		witness := TestVerifierCircuit{}
		err := test.IsSolved(&circuit, &witness, field.TEST_CURVE.ScalarField())
		assert.NoError(err)
	}
	testCase()
}