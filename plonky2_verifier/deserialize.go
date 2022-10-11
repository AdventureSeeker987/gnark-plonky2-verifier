package plonky2_verifier

import (
	"encoding/json"
	. "gnark-ed25519/field"
	"gnark-ed25519/utils"
	"io/ioutil"
	"os"
)

type ProofWithPublicInputsRaw struct {
	Proof struct {
		WiresCap []struct {
			Elements []uint64 `json:"elements"`
		} `json:"wires_cap"`
		PlonkZsPartialProductsCap []struct {
			Elements []uint64 `json:"elements"`
		} `json:"plonk_zs_partial_products_cap"`
		QuotientPolysCap []struct {
			Elements []uint64 `json:"elements"`
		} `json:"quotient_polys_cap"`
		Openings struct {
			Constants       [][]uint64 `json:"constants"`
			PlonkSigmas     [][]uint64 `json:"plonk_sigmas"`
			Wires           [][]uint64 `json:"wires"`
			PlonkZs         [][]uint64 `json:"plonk_zs"`
			PlonkZsNext     [][]uint64 `json:"plonk_zs_next"`
			PartialProducts [][]uint64 `json:"partial_products"`
			QuotientPolys   [][]uint64 `json:"quotient_polys"`
		} `json:"openings"`
		OpeningProof struct {
			CommitPhaseMerkleCaps []interface{} `json:"commit_phase_merkle_caps"`
			QueryRoundProofs      []struct {
				InitialTreesProof struct {
					EvalsProofs [][]interface{} `json:"evals_proofs"`
				} `json:"initial_trees_proof"`
				Steps []interface{} `json:"steps"`
			} `json:"query_round_proofs"`
			FinalPoly struct {
				Coeffs [][]uint64 `json:"coeffs"`
			} `json:"final_poly"`
			PowWitness uint64 `json:"pow_witness"`
		} `json:"opening_proof"`
	} `json:"proof"`
	PublicInputs []uint64 `json:"public_inputs"`
}

type CommonCircuitDataRaw struct {
	Config struct {
		NumWires                uint64 `json:"num_wires"`
		NumRoutedWires          uint64 `json:"num_routed_wires"`
		NumConstants            uint64 `json:"num_constants"`
		UseBaseArithmeticGate   bool   `json:"use_base_arithmetic_gate"`
		SecurityBits            uint64 `json:"security_bits"`
		NumChallenges           uint64 `json:"num_challenges"`
		ZeroKnowledge           bool   `json:"zero_knowledge"`
		MaxQuotientDegreeFactor uint64 `json:"max_quotient_degree_factor"`
		FriConfig               struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []int `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"fri_config"`
	} `json:"config"`
	FriParams struct {
		Config struct {
			RateBits          uint64 `json:"rate_bits"`
			CapHeight         uint64 `json:"cap_height"`
			ProofOfWorkBits   uint64 `json:"proof_of_work_bits"`
			ReductionStrategy struct {
				ConstantArityBits []uint64 `json:"ConstantArityBits"`
			} `json:"reduction_strategy"`
			NumQueryRounds uint64 `json:"num_query_rounds"`
		} `json:"config"`
		Hiding             bool          `json:"hiding"`
		DegreeBits         uint64        `json:"degree_bits"`
		ReductionArityBits []interface{} `json:"reduction_arity_bits"`
	} `json:"fri_params"`
	DegreeBits    uint64 `json:"degree_bits"`
	SelectorsInfo struct {
		SelectorIndices []uint64 `json:"selector_indices"`
		Groups          []struct {
			Start uint64 `json:"start"`
			End   uint64 `json:"end"`
		} `json:"groups"`
	} `json:"selectors_info"`
	QuotientDegreeFactor uint64   `json:"quotient_degree_factor"`
	NumGateConstraints   uint64   `json:"num_gate_constraints"`
	NumConstants         uint64   `json:"num_constants"`
	NumPublicInputs      uint64   `json:"num_public_inputs"`
	KIs                  []uint64 `json:"k_is"`
	NumPartialProducts   uint64   `json:"num_partial_products"`
	CircuitDigest        struct {
		Elements []uint64 `json:"elements"`
	} `json:"circuit_digest"`
}

type VerifierOnlyCircuitDataRaw struct {
	ConstantsSigmasCap []struct {
		Elements []uint64 `json:"elements"`
	} `json:"constants_sigmas_cap"`
}

func DeserializeMerkleCap(merkleCapRaw []struct{ Elements []uint64 }) MerkleCap {
	n := len(merkleCapRaw)
	merkleCap := make([]Hash, n)
	for i := 0; i < n; i++ {
		copy(merkleCap[i][:], utils.Uint64ArrayToFArray(merkleCapRaw[i].Elements))
	}
	return merkleCap
}

func DeserializeOpeningSet(openingSetRaw struct {
	Constants       [][]uint64
	PlonkSigmas     [][]uint64
	Wires           [][]uint64
	PlonkZs         [][]uint64
	PlonkZsNext     [][]uint64
	PartialProducts [][]uint64
	QuotientPolys   [][]uint64
}) OpeningSet {
	return OpeningSet{
		Constants:       utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.Constants),
		PlonkSigmas:     utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.PlonkSigmas),
		Wires:           utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.Wires),
		PlonkZs:         utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.PlonkZs),
		PlonkZsNext:     utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.PlonkZsNext),
		PartialProducts: utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.PartialProducts),
		QuotientPolys:   utils.Uint64ArrayToQuadraticExtensionArray(openingSetRaw.QuotientPolys),
	}
}

func DeserializeFriProof(openingProofRaw struct {
	CommitPhaseMerkleCaps []interface{}
	QueryRoundProofs      []struct {
		InitialTreesProof struct {
			EvalsProofs [][]interface{}
		}
		Steps []interface{}
	}
	FinalPoly struct {
		Coeffs [][]uint64
	}
	PowWitness uint64
}) FriProof {
	var openingProof FriProof
	openingProof.PowWitness = NewFieldElement(openingProofRaw.PowWitness)
	openingProof.FinalPoly.Coeffs = utils.Uint64ArrayToQuadraticExtensionArray(openingProofRaw.FinalPoly.Coeffs)
	return openingProof
}

func DeserializeProofWithPublicInputs(path string) ProofWithPublicInputs {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()
	rawBytes, _ := ioutil.ReadAll(jsonFile)

	var raw ProofWithPublicInputsRaw
	err = json.Unmarshal(rawBytes, &raw)
	if err != nil {
		panic(err)
	}

	var proofWithPis ProofWithPublicInputs
	proofWithPis.Proof.WiresCap = DeserializeMerkleCap([]struct{ Elements []uint64 }(raw.Proof.WiresCap))
	proofWithPis.Proof.PlonkZsPartialProductsCap = DeserializeMerkleCap([]struct{ Elements []uint64 }(raw.Proof.PlonkZsPartialProductsCap))
	proofWithPis.Proof.QuotientPolysCap = DeserializeMerkleCap([]struct{ Elements []uint64 }(raw.Proof.QuotientPolysCap))
	proofWithPis.Proof.Openings = DeserializeOpeningSet(struct {
		Constants       [][]uint64
		PlonkSigmas     [][]uint64
		Wires           [][]uint64
		PlonkZs         [][]uint64
		PlonkZsNext     [][]uint64
		PartialProducts [][]uint64
		QuotientPolys   [][]uint64
	}(raw.Proof.Openings))
	proofWithPis.Proof.OpeningProof = DeserializeFriProof(struct {
		CommitPhaseMerkleCaps []interface{}
		QueryRoundProofs      []struct {
			InitialTreesProof struct{ EvalsProofs [][]interface{} }
			Steps             []interface{}
		}
		FinalPoly  struct{ Coeffs [][]uint64 }
		PowWitness uint64
	}(raw.Proof.OpeningProof))
	proofWithPis.PublicInputs = utils.Uint64ArrayToFArray(raw.PublicInputs)

	return proofWithPis
}

func DeserializeCommonCircuitData(path string) CommonCircuitDataRaw {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()
	rawBytes, _ := ioutil.ReadAll(jsonFile)

	var raw CommonCircuitDataRaw
	err = json.Unmarshal(rawBytes, &raw)
	if err != nil {
		panic(err)
	}

	return raw
}

func DeserializeVerifierOnlyCircuitData(path string) VerifierOnlyCircuitData {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()
	rawBytes, _ := ioutil.ReadAll(jsonFile)

	var raw VerifierOnlyCircuitDataRaw
	err = json.Unmarshal(rawBytes, &raw)
	if err != nil {
		panic(err)
	}

	return VerifierOnlyCircuitData{
		ConstantSigmasCap: DeserializeMerkleCap([]struct{ Elements []uint64 }(raw.ConstantsSigmasCap)),
	}
}
