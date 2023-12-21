// Copyright 2019-2020 Lawrence Livermore National Security, LLC and other
// Archspec Project Developers. See the top-level COPYRIGHT file for details.
//
// SPDX-License-Identifier: (Apache-2.0 OR MIT)

package cpu

import (
	"fmt"
	"testing"
)

const tick = "\u2713"
const cross = "\u2717"

func TestAncestors(t *testing.T) {

	testTargets := map[string]int{
		"icelake":     16,
		"k10":         1,
		"steamroller": 4,
	}

	for name, expected := range testTargets {
		t.Run(name, func(t *testing.T) {
			tgt := TARGETS[name]
			if len(tgt.Ancestors()) != expected {
				t.Error(name, "got", len(tgt.Ancestors()), ", but expected", expected, cross)
				return
			}
			t.Log(name, tick)
		})
	}

}

func TestFamily(t *testing.T) {

	testTargets := map[string]string{
		"icelake":     "x86_64",
		"k10":         "x86_64",
		"steamroller": "x86_64",
		"thunderx2":   "aarch64",
		"power8le":    "ppc64le",
	}

	for target, family := range testTargets {
		t.Run(target, func(t *testing.T) {
			tgt := TARGETS[target]
			if tgt.Family().Name != family {
				t.Error(target, cross)
				return
			}
			t.Log(target, tick)
		})
	}
}

func TestSupports(t *testing.T) {

	testSupported := map[string]string{
		"icelake":        "avx",
		"k10":            "sse4a",
		"steamroller":    "bmi1",
		"a64fx":          "sve",
		"thunderx2":      "neon",
		"power8le":       "altivec",
		"broadwell":      "sse3",
		"haswell":        "ssse3",
		"skylake_avx512": "avx512",
	}

	for target, isa := range testSupported {
		t.Run(target, func(t *testing.T) {
			tgt := TARGETS[target]
			if !tgt.Supports(isa) {
				t.Error(target, "does not support", isa, cross)
				return
			}
			t.Log(target, "supports", isa, tick)
		})
	}

	testNotSupported := map[string]string{
		"icelake":   "neon",
		"broadwell": "doesnotexist",
	}

	for target, isa := range testNotSupported {
		t.Run(target, func(t *testing.T) {
			tgt := TARGETS[target]
			if tgt.Supports(isa) {
				t.Error(target, "supports", isa, cross)
				return
			}
			t.Log(target, "does not support", isa, tick)
		})
	}
}

func TestCompatibleWith(t *testing.T) {

	testCompatible := map[string]string{
		"cannonlake":  "icelake",
		"aarch64":     "thunderx2",
		"steamroller": "steamroller",
		"x86_64":      "bulldozer",
	}

	for target, laterUarch := range testCompatible {
		t.Run(target, func(t *testing.T) {
			tgt := TARGETS[target]
			ltr := TARGETS[laterUarch]
			if !tgt.CompatibleWith(ltr) {
				t.Error(target, "not compatible with", laterUarch, cross)
				return
			}
			t.Log(target, "compatible with", laterUarch, tick)
		})
	}

	testNotCompatible := map[string]string{
		"icelake":     "thunderx2",
		"thunderx2":   "aarch64",
		"steamroller": "bulldozer",
	}

	for target, earlierUarch := range testNotCompatible {
		t.Run(target, func(t *testing.T) {
			tgt := TARGETS[target]
			ltr := TARGETS[earlierUarch]
			if tgt.CompatibleWith(ltr) {
				t.Error(target, "compatible with", earlierUarch, cross)
				return
			}
			t.Log(target, "not compatible with", earlierUarch, tick)
		})
	}
}

func TestCompilers(t *testing.T) {
	testTargets := map[string]int{
		"icelake":     8,
		"k10":         7,
		"steamroller": 7,
		"thunderx2":   2,
		"power8le":    3,
	}

	for target, count := range testTargets {
		t.Run(target, func(t *testing.T) {
			tgt := TARGETS[target]
			fmt.Printf("Found %d in family %s\n", len(tgt.Compilers), tgt.Name)

			// Ensure fields are not empty!
			for _, compilers := range tgt.Compilers {
				for _, compiler := range compilers {
					if compiler.Flags == "" {
						t.Error(target, cross)
						return
					}
					if compiler.Versions == "" {
						t.Error(target, cross)
						return
					}
				}
			}
			if len(tgt.Compilers) != count {
				t.Error(target, cross)
				return
			}
			t.Log(target, tick)
		})
	}
}
