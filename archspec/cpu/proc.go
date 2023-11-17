package cpu

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/scylladb/go-set/strset"
)

var ErrKeyNotFound = errors.New("missing key in cpuinfo")

type cMap struct {
	Name   string
	Length int
	MArch  Microarchitecture
}

func ParseCPUInfo(cpuinfo map[string]string) (march Microarchitecture, err error) {
	candidateMap := []cMap{}
	vndID, ok := cpuinfo["vendor_id"]
	if !ok {
		err = fmt.Errorf("%q: %w", "vendor_id", ErrKeyNotFound)
		return
	}
	flags, ok := cpuinfo["flags"]
	if !ok {
		err = fmt.Errorf("%q: %w", "flags", ErrKeyNotFound)
		return
	}

	for name, value := range JSONData.Microarchitectures {
		if value.Vendor != vndID {
			continue
		}
		flagSlice := strings.Split(flags, " ")
		flagSet := strset.New(flagSlice...)
		entrySet := strset.New(value.Features...)
		superS := flagSet.IsSubset(entrySet)
		if superS {
			candidateMap = append(candidateMap, cMap{Name: name, Length: len(entrySet.List()), MArch: MicroarchFromJMA(name, value)})
		}
	}
	// Sort by age preserving name order
	sort.SliceStable(candidateMap, func(i, j int) bool { return candidateMap[i].Length > candidateMap[j].Length })
	if len(candidateMap) == 0 {
		err = errors.New("no matching microarchitecture found")
	}
	march = candidateMap[0].MArch
	return
}
