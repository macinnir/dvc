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

	rcString := ""
	if strings.Contains(inputVersion, "-RC") {
		// Grab the suffix
		rcString = inputVersion[strings.Index(inputVersion, "-RC")+3:]
		// Remove the suffix
		inputVersion = inputVersion[0:strings.Index(inputVersion, "-RC")]
	}

	r := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

	if !r.Match([]byte(inputVersion)) {
		return "", errors.New("invalid format")
	}

	parts := strings.Split(inputVersion, ".")

	major, _ := strconv.ParseInt(parts[0], 10, 64)
	minor, _ := strconv.ParseInt(parts[1], 10, 64)
	patch, _ := strconv.ParseInt(parts[2], 10, 64)

	var rc int64 = 0
	if len(rcString) > 0 {
		rc, _ = strconv.ParseInt(rcString, 10, 64)
	}

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
		suffix = fmt.Sprintf("-RC%d", rc)
	case VersionTypeRelease:
		// All this does is remove the "RC"
	}

	// fmt.Println("Prefix: ", prefix, "Major", major, "Minor", minor, "Patch", patch, "Suffix", suffix)
	newVersion := prefix + fmt.Sprintf("%d.%d.%d", major, minor, patch) + suffix
	// fmt.Println("New Version", newVersion)

	return newVersion, nil

}
