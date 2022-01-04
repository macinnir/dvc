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
	VersionTypeNone    VersionType = "none"
	VersionTypeMajor   VersionType = "major"
	VersionTypeMinor   VersionType = "minor"
	VersionTypePatch   VersionType = "patch"
	VersionTypeRC      VersionType = "rc"
	VersionTypeRelease VersionType = "release"
	VersionTypeAlpha   VersionType = "alpha"
	VersionTypeBeta    VersionType = "beta"
)

// 1. Develop
// 2. Merge to staging branch and tag with next minor version and RC (1.2.3-RC1)

func NextVersion(inputVersion string, inputVersionType string) (string, error) {

	var versionType VersionType = VersionTypeNone
	// versionType := VersionTypePatch

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
		case VersionTypeAlpha:
			versionType = VersionTypeAlpha
		case VersionTypeBeta:
			versionType = VersionTypeBeta
		default:
			return "", errors.New("invalid version type")
		}
	}

	// fmt.Println("VersionType: ", versionType, "Version: ", inputVersion)

	// normalize casing
	inputVersion = strings.ToLower(inputVersion)

	// Prefix
	var prefix string = ""
	if inputVersion[0:1] == "v" {
		prefix = "v"
		inputVersion = inputVersion[1:]
	}

	var rc int64 = 0
	var alpha int64 = 0
	var beta int64 = 0

	// v1.8.5-rc.1
	inputVersion, rc = ParsePreRelease(inputVersion, "rc")
	inputVersion, alpha = ParsePreRelease(inputVersion, "alpha")
	inputVersion, beta = ParsePreRelease(inputVersion, "beta")

	if versionType == VersionTypeNone {
		if rc > 0 {
			versionType = VersionTypeRC
		} else if alpha > 0 {
			versionType = VersionTypeAlpha
		} else if beta > 0 {
			versionType = VersionTypeBeta
		} else {
			versionType = VersionTypePatch
		}
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
	// if rc > 0 && versionType != VersionTypeRelease {
	// 	versionType = VersionTypeRC
	// }

	// If an ALPHA is found in the version, and this is not a release, iterate the ALPHA
	// if alpha > 0 && versionType != VersionTypeRelease {
	// 	// if a 'beta' release was specified, skip to the beta
	// 	if versionType != VersionTypeBeta {
	// 		versionType = VersionTypeAlpha
	// 	}

	// 	// If an BETA is found in the version, and this is not a release, iterate the BETA
	// }

	// if beta > 0 && versionType != VersionTypeRelease {
	// 	versionType = VersionTypeBeta
	// }

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
	case VersionTypeAlpha:

		// If we're actually on a beta or a release, this is considered regressive
		if beta > 0 || rc > 0 {
			return "", errors.New("regressive tagging")
		}

		// This is the first alpha, iterate the minor version
		if alpha == 0 {
			minor++
			patch = 0
		}
		alpha++
		suffix = fmt.Sprintf("-alpha.%d", alpha)
	case VersionTypeBeta:

		// If we're actually on a release, this is considered regressive
		if rc > 0 {
			return "", errors.New("regressive tagging")
		}

		// This is the first beta (and we're not coming from an alpha), iterate the minor version
		if beta == 0 && alpha == 0 {
			minor++
			patch = 0
		}
		beta++
		suffix = fmt.Sprintf("-beta.%d", beta)
	case VersionTypeRC:
		// This is the first RC (and we're not coming from beta or alpha), iterate the minor version
		if rc == 0 && alpha == 0 && beta == 0 {
			minor++
			patch = 0
		}
		rc++
		suffix = fmt.Sprintf("-rc.%d", rc)
	case VersionTypeRelease:
		// All this does is remove the alpha/beta/rc suffix
	}

	// fmt.Println("Prefix: ", prefix, "Major", major, "Minor", minor, "Patch", patch, "Suffix", suffix)
	newVersion := prefix + fmt.Sprintf("%d.%d.%d", major, minor, patch) + suffix
	// fmt.Println("New Version", newVersion)

	return newVersion, nil

}

func ParsePreRelease(inputVersion, suffixType string) (string, int64) {

	var suffixNum int64 = 0

	suffix := fmt.Sprintf("-%s.", suffixType)
	suffixLen := len(suffix)

	// v1.8.5-rc.1
	if strings.Contains(inputVersion, suffix) {

		// Grab the suffix
		suffixNumString := inputVersion[strings.Index(inputVersion, suffix)+suffixLen:]
		// fmt.Println("SuffixNumString: ", suffixNumString)

		// Remove any trailing dashed nonsense
		// v1.8.5-rc.1-12312312312
		if strings.Contains(suffixNumString, "-") {
			parts := strings.Split(suffixNumString, "-")
			suffixNumString = parts[0]
		}

		suffixNum, _ = strconv.ParseInt(suffixNumString, 10, 64)

		// Remove the suffix (and anything after it)
		inputVersion = inputVersion[0:strings.Index(inputVersion, suffix)]
	}

	return inputVersion, suffixNum
}
