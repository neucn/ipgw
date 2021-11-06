package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseVersion(t *testing.T) {
	a := assert.New(t)
	cases := []struct {
		v      string
		parsed *Semver
	}{
		{v: "v1.0.1", parsed: &Semver{
			Major:      1,
			Minor:      0,
			Patch:      1,
			Prerelease: "",
		}},
		{v: "v1.0.1-beta", parsed: &Semver{
			Major:      1,
			Minor:      0,
			Patch:      1,
			Prerelease: "beta",
		}},
		{v: "0.2.0-alpha.3", parsed: &Semver{
			Major:      0,
			Minor:      2,
			Patch:      0,
			Prerelease: "alpha.3",
		}},
		{v: "0.1", parsed: nil},
	}
	for _, c := range cases {
		a.Equal(ParseVersion(c.v), c.parsed)
	}
}

func TestCompareVersion(t *testing.T) {
	a := assert.New(t)
	cases := []struct {
		a        string
		b        string
		expected bool
	}{
		{a: "v1.0.0", b: "v0.2.2", expected: true},
		{a: "v1.0.0-alpha", b: "v1.2.2", expected: false},
		{a: "v1.0.0-alpha", b: "v1.0.0", expected: false},
		{a: "v1.0.0", b: "v1.0.0-beta.2", expected: true},
		{a: "v1.0.0-beta", b: "v1.0.0-alpha", expected: true},
		{a: "v1", b: "v1.0.0-alpha", expected: false},
		{a: "1.0.0-beta", b: "", expected: true},
	}
	for _, c := range cases {
		a.Equal(CompareVersion(ParseVersion(c.a), ParseVersion(c.b)), c.expected)
	}
}
