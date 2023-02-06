package cpu

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	haswell = map[string]string{
		"vendor_id": "GenuineIntel",
		"flags":     "fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 fma cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm epb invpcid_single tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid cqm xsaveopt cqm_llc cqm_occup_llc ibpb ibrs stibp dtherm arat pln pts spec_ctrl intel_stibp",
	}
)

func TestProcInfo(t *testing.T) {
	_, err := ParseCPUInfo(map[string]string{})
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("wrong error: %v", err)
	}
	_, err = ParseCPUInfo(map[string]string{"vendor_id": "GenuineIntel"})
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("wrong error: %v", err)
	}
	_, err = ParseCPUInfo(map[string]string{"flags": "sse"})
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("wrong error: %v", err)
	}
	maHW, err := ParseCPUInfo(haswell)
	assert.NoError(t, err)
	assert.Equal(t, "haswell", maHW.Name)
}

func TestProcInfoFiles(t *testing.T) {
	cpuinfo, err := ParseTestFile("../json/tests/targets/linux-centos7-cascadelake")
	assert.NoError(t, err)
	maHW, err := ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "cascadelake", maHW.Name)
	cpuinfo, err = ParseTestFile("../json/tests/targets/linux-rhel7-haswell")
	assert.NoError(t, err)
	maHW, err = ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "haswell", maHW.Name)
	cpuinfo, err = ParseTestFile("../json/tests/targets/linux-rhel7-skylake_avx512")
	assert.NoError(t, err)
	maHW, err = ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "skylake_avx512", maHW.Name)
	cpuinfo, err = ParseTestFile("../json/tests/targets/linux-rhel7-zen")
	assert.NoError(t, err)
	maHW, err = ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "zen", maHW.Name, "wrong arch: %v", maHW)
	// Zen3
	cpuinfo, err = ParseTestFile("../json/tests/targets/linux-ubuntu20.04-zen3")
	assert.NoError(t, err)
	maHW, err = ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "zen3", maHW.Name)
}

func TestParse_Family(t *testing.T) {
	cpuinfo, err := ParseTestFile("../json/tests/targets/linux-rhel7-zen")
	assert.NoError(t, err)
	maHW, err := ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "zen", maHW.Name, "wrong arch: %v", maHW)
	assert.Len(t, maHW.Parents, 1)
	assert.Equal(t, "x86_64_v3", maHW.Parents[0].Name)
	// Zen3
	cpuinfo, err = ParseTestFile("../json/tests/targets/linux-ubuntu20.04-zen3")
	assert.NoError(t, err)
	maHW, err = ParseCPUInfo(cpuinfo)
	assert.NoError(t, err)
	assert.Equal(t, "zen3", maHW.Name)
	// TODO: How can I get the genric architecture from the cpuinfo?
	// -> Iterate through the parents and find a generic one?
	assert.Equal(t, "zen2", maHW.Parents[0].Name)
}

func ParseTestFile(fname string) (map[string]string, error) {
	cpuinfo := map[string]string{}
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		if len(s) == 2 {
			cpuinfo[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
		}
	}
	return cpuinfo, nil
}
