// Copyright 2019-2020 Lawrence Livermore National Security, LLC and other
// Archspec Project Developers. See the top-level COPYRIGHT file for details.
//
// SPDX-License-Identifier: (Apache-2.0 OR MIT)

package cpu

import (
	"errors"

	"github.com/scylladb/go-set/strset"
)

// Microarchitecture models a CPU microarchitecture
type Microarchitecture struct {
	Name       string
	Parents    []Microarchitecture
	Vendor     string
	Features   strset.Set
	Generation int
}

// TARGETS is a list of all the CPU microarchitectures known to the package
var TARGETS = parseTargets()

// parseTargets constructs the list of known targets from the information read
// from the embedded JSON file.
func parseTargets() map[string]Microarchitecture {
	targets := make(map[string]Microarchitecture)
	for name, value := range JSONData.Microarchitectures {
		// If this target is already constructed proceed to the next
		if _, ok := targets[name]; ok {
			continue
		}

		// Otherwise add it
		addItem(targets, name, value)
	}
	return targets
}

// addItem adds a single item to a list of targets
func addItem(targets map[string]Microarchitecture, name string, value jMicroarchitecture) {
	// Normalize the parents field to a list of strings
	parentNames := make([]string, 0)
	switch value.From.(type) {
	case nil:
		// Nothing to do, this is an architecture family like x86_64 or AArch64
	case string:
		// The microarchitecture has a single parent
		parentNames = append(parentNames, value.From.(string))
	case []interface{}:
		// The microarchitecture has a multiple parents (e.g. icelake)
		for _, item := range value.From.([]interface{}) {
			parentNames = append(parentNames, item.(string))
		}
	}

	// Ensure parents have been constructed already
	parents := make([]Microarchitecture, 0)
	for _, pname := range parentNames {
		if _, ok := targets[pname]; !ok {
			addItem(targets, pname, JSONData.Microarchitectures[pname])
		}
		parents = append(parents, targets[pname])
	}

	targets[name] = Microarchitecture{
		Name:       name,
		Vendor:     value.Vendor,
		Features:   *strset.New(value.Features...),
		Generation: value.Generation,
		Parents:    parents,
	}
}

// Ancestors returns the list of all the ancestors of the current microarchitecture
func (m Microarchitecture) Ancestors() []Microarchitecture {
	ancestors := make([]Microarchitecture, 0)
	seen := make(map[string]struct{})
	// First add parents
	for _, p := range m.Parents {
		seen[p.Name] = struct{}{}
		ancestors = append(ancestors, p)
	}
	// Then their ancestors
	for _, p := range m.Parents {
		for _, a := range p.Ancestors() {
			if _, thereAlready := seen[a.Name]; thereAlready {
				continue
			}
			seen[a.Name] = struct{}{}
			ancestors = append(ancestors, a)
		}
	}
	return ancestors
}

// Family returns the generic Architecture of this microarchitecture
func (m Microarchitecture) Family() Microarchitecture {
	for _, a := range m.Ancestors() {
		if len(a.Ancestors()) == 0 {
			return a
		}
	}
	panic(errors.New("microarchitecture has no valid family"))
}

// Supports returns True if the instructionSet is supported by the microarchitecture, false otherwise
func (m Microarchitecture) Supports(instructionSet string) bool {
	// TODO: might extend the API to support variadic args
	// Check if there's a feature that satisfies the request
	if ok := m.Features.Has(instructionSet); ok {
		return true
	}

	// If there's no verbatim feature check aliases
	if _, ok := featureAliases[instructionSet]; ok {
		return featureAliases[instructionSet].Evaluate(m)
	}
	return false
}

// CompatibleWith returns True if binaries that run on the current microarchitecture
// can also run on the microarchitecture passed as argument i.e. the current
// microarchitecture is one of the ancestor
func (m Microarchitecture) CompatibleWith(other Microarchitecture) bool {
	ancestorNames := make([]string, 0)
	ancestorNames = append(ancestorNames, other.Name)
	for _, a := range other.Ancestors() {
		ancestorNames = append(ancestorNames, a.Name)
	}

	otherSet := *strset.New(ancestorNames...)
	return otherSet.Has(m.Name)
}
