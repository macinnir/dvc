package versioner

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Development
// v0.1.0-RC1, v0.1.0-RC2, ... 		-- Feature complete. Testing/patching as release candidates
// v0.1.0 							-- Testing complete. Release.
// v0.1.1 							-- Post-release hot-fix

type VersionType string

const (
	VersionTypeMajor   VersionType = "major"
	VersionTypeMinor   VersionType = "minor"
	VersionTypePatch   VersionType = "patch"
	VersionTypeRC      VersionType = "rc"
	VersionTypeRelease VersionType = "release"
)

// 1. Develop
// 2. Merge to staging branch and tag with next minor version and RC (1.2.3-RC1)

func NextVersion(inputVersion string, inputVersionType string) (string, error) {

	versionType := VersionTypePatch

	if inputVersionType != "" {

		switch VersionType(inputVersionType) {
		case VersionTypeMajor:
			versionType = VersionTypeMajor
		case VersionTypeMinor:
			versionType = VersionTypeMinor
		case VersionTypePatch:
			versionType = VersionTypePatch
		case VersionTypeRC:
			versionType = VersionTypeRC
		case VersionTypeRelease:
			versionType = VersionTypeRelease
		default:
			return "", errors.New("invalid version type")
		}
	}

	var prefix string = ""

	if strings.ToLower(inputVersion[0:1]) == "v" {
		prefix = "v"
		inputVersion = inputVersion[1:]
	}

	// normalize casing
	inputVersion = strings.ToLower(inputVersion)

	rcString := ""
	var rc int64 = 0

	// v1.8.5-rc.1
	if strings.Contains(inputVersion, "-rc.") {
		// Grab the suffix
		rcString = inputVersion[strings.Index(inputVersion, "-rc.")+4:]

		if strings.Contains(rcString, "-") {
			parts := strings.Split(rcString, "-")
			rcString = parts[0]
		}

		rc, _ = strconv.ParseInt(rcString, 10, 64)

		// Remove the suffix
		inputVersion = inputVersion[0:strings.Index(inputVersion, "-rc.")]
	}

	// Could be something like v1.8.58-5-g98e9b2b
	if strings.Contains(inputVersion, "-") {
		parts := strings.Split(inputVersion, "-")
		inputVersion = parts[0]
	}

	r := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

	if !r.Match([]byte(inputVersion)) {
		return "", errors.New("invalid format")
	}

	parts := strings.Split(inputVersion, ".")

	major, _ := strconv.ParseInt(parts[0], 10, 64)
	minor, _ := strconv.ParseInt(parts[1], 10, 64)
	patch, _ := strconv.ParseInt(parts[2], 10, 64)

	// If an RC is found in the version, and this is not a release, iterate the RC
	if rc > 0 && versionType != VersionTypeRelease {
		versionType = VersionTypeRC
	}

	// fmt.Println("VersionType: ", versionType, "RC:", rc)

	suffix := ""

	switch versionType {
	case VersionTypeMajor:
		major++
		minor = 0
		patch = 0
	case VersionTypeMinor:
		minor++
		patch = 0
	case VersionTypePatch:
		patch++
	case VersionTypeRC:
		// This is the first RC, iterate the minor version
		if rc == 0 {
			minor++
			patch = 0
		}
		rc++
		suffix = fmt.Sprintf("-rc.%d", rc)
	case VersionTypeRelease:
		// All this does is remove the "RC"
	}

	// fmt.Println("Prefix: ", prefix, "Major", major, "Minor", minor, "Patch", patch, "Suffix", suffix)
	newVersion := prefix + fmt.Sprintf("%d.%d.%d", major, minor, patch) + suffix
	// fmt.Println("New Version", newVersion)

	return newVersion, nil

}
