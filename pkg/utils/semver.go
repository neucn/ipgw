package utils

import (
	"strconv"
	"strings"
)

type Semver struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

// CompareVersion compares strings like "v1.1.1", returns true when versionA > versionB
func CompareVersion(versionA, versionB *Semver) bool {
	if versionA == nil {
		return false
	}
	if versionB == nil {
		return true
	}
	// semver:
	//   v1.0.1 > v1.0.0
	//   v1.0.0 > v1.0.0-beta
	//   v1.0.0-beta.1 > v1.0.0-beta
	//   v1.0.0-beta > v1.0.0-alpha
	if versionA.Major != versionB.Major {
		return versionA.Major > versionB.Major
	}

	if versionA.Minor != versionB.Minor {
		return versionA.Minor > versionB.Minor
	}

	if versionA.Patch != versionB.Patch {
		return versionA.Patch > versionB.Patch
	}

	if versionA.Prerelease != "" && versionB.Prerelease != "" {
		return versionA.Prerelease > versionB.Prerelease
	}

	if versionA.Prerelease != "" {
		// versionB.Prerelease == ""
		return false
	}

	if versionB.Prerelease != "" {
		// versionA.Prerelease == ""
		return true
	}

	return false
}

func ParseVersion(version string) *Semver {
	if version == "" {
		return nil
	}
	// remove leading 'v'
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, "-")
	if len(parts) < 1 {
		return nil
	}
	dot := strings.Split(parts[0], ".")
	if len(dot) < 3 {
		return nil
	}
	preRelease := ""
	if len(parts) > 1 {
		preRelease = parts[1]
	}
	major, _ := strconv.Atoi(dot[0])
	minor, _ := strconv.Atoi(dot[1])
	patch, _ := strconv.Atoi(dot[2])
	return &Semver{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: preRelease,
	}
}
