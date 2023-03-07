package verifier

import (
	"github.com/taramakage/gon-verifier/internal/chain"
)

type Registry struct {
	vs map[string]Verifier
}

func NewRegistry(r *chain.Registry) *Registry {
	vs := map[string]Verifier{
		"A1": A1Verifier{r},
		"A2": A2Verifier{r},
		"A3": A3Verifier{r},
		"A4": A4Verifier{r},
		"A5": A5Verifier{r},
		"A6": A6Verifier{r},
	}
	return &Registry{vs}
}

// Get returns a verifier by key.
func (r *Registry) Get(key string) Verifier {
	return r.vs[key]
}
