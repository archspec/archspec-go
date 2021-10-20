// Copyright 2019-2020 Lawrence Livermore National Security, LLC and other
// Archspec Project Developers. See the top-level COPYRIGHT file for details.
//
// SPDX-License-Identifier: (Apache-2.0 OR MIT)

package cpu

import (
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
