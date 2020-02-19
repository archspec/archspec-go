package cpu

import (
	"fmt"

	"github.com/scylladb/go-set/strset"
)

type aliasPredicate struct {
	Reason   string
	AnyOf    []string
	Families strset.Set
}

var featureAliases = parseFeatureAliases()

func parseFeatureAliases() map[string]aliasPredicate {
	aliases := make(map[string]aliasPredicate)
	for name, value := range JSONData.FeatureAliases {
		aliases[name] = aliasPredicate{
			Reason:   value.Reason,
			AnyOf:    value.AnyOf,
			Families: *strset.New(value.Families...),
		}
	}
	fmt.Println(aliases)
	return aliases
}

func (a aliasPredicate) Evaluate(m Microarchitecture) bool {
	result := true
	// This predicate is true if any of the flags are supported.
	if len(a.AnyOf) != 0 {
		result = result && m.Features.HasAny(a.AnyOf...)
	}
	// This predicate is true if the family is among the allowed ones
	if !a.Families.IsEmpty() {
		result = result && a.Families.Has(m.Family().Name)
	}

	return result
}
