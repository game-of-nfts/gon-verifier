package verifier

import (
	"github.com/taramakage/gon-verifier/internal/chain"
)

type Registry struct {
	vs map[string]Verifier
}

func NewRegistry(r *chain.Registry) *Registry {
	vs := map[string]Verifier{
		"A1":  A1Verifier{r},
		"A2":  A2Verifier{r},
		"A3":  A3Verifier{r},
		"A4":  A4Verifier{r},
		"A5":  A5Verifier{r},
		"A6":  A6Verifier{r},
		"A7":  NewFlowVerifier(r, "a01", true),
		"A8":  NewFlowVerifier(r, "a02", true),
		"A9":  NewFlowVerifier(r, "a03", true),
		"A10": NewFlowVerifier(r, "a04", true),
		"A11": NewFlowVerifier(r, "a05", true),
		"A12": NewFlowVerifier(r, "a06", true),
		"A13": NewFlowVerifier(r, "b01", false),
		"A14": NewFlowVerifier(r, "b02", false),
		"A15": NewFlowVerifier(r, "b03", false),
		"A16": NewFlowVerifier(r, "b04", false),
		"A17": NewFlowVerifier(r, "c01", false),
		"A18": NewFlowVerifier(r, "c02", false),
		"A19": NewFlowVerifier(r, "c03", false),
		"A20": NewFlowVerifier(r, "c04", false),
	}
	return &Registry{vs}
}

// Get returns a verifier by key.
func (r *Registry) Get(key string) Verifier {
	return r.vs[key]
}
