// Copyright 2019-2020 Lawrence Livermore National Security, LLC and other
// Archspec Project Developers. See the top-level COPYRIGHT file for details.
//
// SPDX-License-Identifier: (Apache-2.0 OR MIT)

package cpu

import (
	"encoding/json"
	"log"

	"github.com/archspec/archspec-go/archspec"
)

// JSONData is an in memory representation of the JSON file
// microarchitectures.json
var JSONData = parseJSONDocument()

type jMicroarchitecture struct {
	From       interface{}            `json:"from"`
	Vendor     string                 `json:"vendor"`
	Features   []string               `json:"features"`
	Generation int                    `json:"generation"`
	Compilers  map[string][]jCompiler `json:"compilers,omitempty"`
}

type jCompiler struct {
	Name     string `json:"name,omitempty"`
	Flags    string `json:"flags,omitempty"`
	Versions string `json:"versions,omitempty"`
}

type jFeatureAlias struct {
	Reason   string   `json:"reason"`
	AnyOf    []string `json:"any_of"`
	Families []string `json:"families"`
}

type jConversions struct {
	Description string            `json:"description"`
	ArmVendors  map[string]string `json:"arm_vendors"`
	DarwinFlags map[string]string `json:"darwin_flags"`
}

type jDocument struct {
	Microarchitectures map[string]jMicroarchitecture `json:"microarchitectures"`
	FeatureAliases     map[string]jFeatureAlias      `json:"feature_aliases"`
	Conversions        jConversions                  `json:"conversions"`
}

func parseJSONDocument() *jDocument {
	filename := "json/cpu/microarchitectures.json"
	file, err := archspec.JSONdirectory.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var document *jDocument
	json.NewDecoder(file).Decode(&document)
	return document
}
