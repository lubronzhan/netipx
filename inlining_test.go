// Copyright 2020 The Inet.Af AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netaddr

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestInlining(t *testing.T) {
	if v := runtime.Version(); strings.HasPrefix(v, "go1.14") ||
		strings.HasPrefix(v, "go1.13") ||
		strings.HasPrefix(v, "go1.12") ||
		strings.HasPrefix(v, "go1.11") {
		t.Skipf("skipping test on old Go version %q", v)
	}
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
	out, err := exec.Command(
		filepath.Join(runtime.GOROOT(), "bin", "go"+exe),
		"build",
		"--gcflags=-m",
		"inet.af/netaddr").CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %v, %s", err, out)
	}
	got := map[string]bool{}
	regexp.MustCompile(` can inline (\S+)`).ReplaceAllFunc(out, func(match []byte) []byte {
		got[strings.TrimPrefix(string(match), " can inline ")] = true
		return nil
	})
	for _, want := range []string{
		"(*IPSet).Add",
		"(*IPSet).Remove",
		"(*IPSet).RemoveRange",
		"(*uint128).clear",
		"(*uint128).set",
		"IP.BitLen",
		"IP.IPAddr",
		"IP.Is4",
		"IP.Is4in6",
		"IP.Is6",
		"IP.IsLoopback",
		"IP.IsMulticast",
		"IP.IsZero",
		"IP.Less",
		"IP.MarshalText",
		"IP.Unmap",
		"IP.Zone",
		"IP.hi",
		"IP.lo",
		"IP.v4",
		"IP.v6",
		"IP.withInternedZone",
		"IPPort.IsZero",
		"IPPort.MarshalText",
		"IPPort.TCPAddr",
		"IPPrefix.IsSingleIP",
		"IPPrefix.IsZero",
		"IPPrefix.Masked",
		"IPPrefix.String",
		"IPv4",
		"IPv6LinkLocalAllNodes",
		"IPv6Unspecified",
		"MustParseIP",
		"MustParseIPPort",
		"MustParseIPPrefix",
		"PutUDPAddr",
		"discardf",
		"mask4",
		"mask6",
		"uint128.bitSet",
	} {
		if !got[want] {
			t.Errorf("%q is no longer inlinable", want)
			continue
		}
		delete(got, want)
	}
	for sym := range got {
		if strings.Contains(sym, ".func") {
			continue
		}
		t.Logf("not in expected set, but also inlinable: %q", sym)

	}
}