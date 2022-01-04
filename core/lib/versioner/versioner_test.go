package versioner_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/macinnir/dvc/core/lib/versioner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersioner(t *testing.T) {

	tests := []struct {
		input       string
		versionType string
		e           error
		output      string
	}{
		{input: "v1.2.3", versionType: "", e: nil, output: "v1.2.4"},
		{input: "1.2.3", versionType: "", e: nil, output: "1.2.4"},
		{input: "1.2.3", versionType: "major", e: nil, output: "2.0.0"},
		{input: "1.2.3", versionType: "minor", e: nil, output: "1.3.0"},
		{input: "1.2.3", versionType: "patch", e: nil, output: "1.2.4"},
		{input: "v1.2.3", versionType: "major", e: nil, output: "v2.0.0"},
		{input: "v1.2.3", versionType: "minor", e: nil, output: "v1.3.0"},
		{input: "v1.2.3", versionType: "patch", e: nil, output: "v1.2.4"},
		{input: "v1.2.3", versionType: "invalid type", e: errors.New("invalid version type"), output: ""},

		{input: "v1.2.3-rc.1", versionType: "rc", e: nil, output: "v1.2.3-rc.2"},
		{input: "v1.2.3-RC.1", versionType: "rc", e: nil, output: "v1.2.3-rc.2"},
		{input: "v1.2.3-RC1", versionType: "rc", e: nil, output: "v1.3.0-rc.1"},
		{input: "1.2.3-rc.1", versionType: "rc", e: nil, output: "1.2.3-rc.2"},
		{input: "1.2.3-rc.1", versionType: "rc", e: nil, output: "1.2.3-rc.2"},
		{input: "1.2.3-rc.2", versionType: "release", e: nil, output: "1.2.3"},
		{input: "1.2.3-rc.2", versionType: "", e: nil, output: "1.2.3-rc.3"},
		{input: "1.2.3-rc.2-5-g98e9b2b", versionType: "", e: nil, output: "1.2.3-rc.3"},
		{input: "1.2.3-rc.2-5-g98e9b2b", versionType: "major", e: nil, output: "2.0.0"},

		{input: "v1.2.3", versionType: "alpha", e: nil, output: "v1.3.0-alpha.1"},
		{input: "v1.2.3-alpha.1", versionType: "alpha", e: nil, output: "v1.2.3-alpha.2"},
		{input: "v1.2.3-ALPHA.1", versionType: "alpha", e: nil, output: "v1.2.3-alpha.2"},
		{input: "v1.2.3-ALPHA1", versionType: "alpha", e: nil, output: "v1.3.0-alpha.1"},
		{input: "1.2.3-alpha.1", versionType: "alpha", e: nil, output: "1.2.3-alpha.2"},
		{input: "1.2.3-alpha.1", versionType: "alpha", e: nil, output: "1.2.3-alpha.2"},
		{input: "1.2.3-alpha.2", versionType: "release", e: nil, output: "1.2.3"},
		{input: "1.2.3-alpha.2", versionType: "", e: nil, output: "1.2.3-alpha.3"},
		{input: "1.2.3-alpha.2-5-g98e9b2b", versionType: "", e: nil, output: "1.2.3-alpha.3"},

		{input: "v1.2.3-alpha.1", versionType: "beta", e: nil, output: "v1.2.3-beta.1"},
		{input: "v1.2.3-beta.1", versionType: "beta", e: nil, output: "v1.2.3-beta.2"},
		{input: "v1.2.3-BETA.1", versionType: "beta", e: nil, output: "v1.2.3-beta.2"},
		{input: "v1.2.3-BETA1", versionType: "beta", e: nil, output: "v1.3.0-beta.1"},
		{input: "1.2.3-beta.1", versionType: "beta", e: nil, output: "1.2.3-beta.2"},
		{input: "1.2.3-beta.1", versionType: "beta", e: nil, output: "1.2.3-beta.2"},
		{input: "1.2.3-beta.2", versionType: "release", e: nil, output: "1.2.3"},
		{input: "1.2.3-beta.2", versionType: "", e: nil, output: "1.2.3-beta.3"},
		{input: "1.2.3-beta.2-5-g98e9b2b", versionType: "", e: nil, output: "1.2.3-beta.3"},

		{input: "0.4.0-rc.2-1-gfd4c000", versionType: "patch", e: nil, output: "0.4.1"},
		{input: "v1.2.3", versionType: "rc", e: nil, output: "v1.3.0-rc.1"},
		{input: "1.2.3", versionType: "rc", e: nil, output: "1.3.0-rc.1"},
		{input: "1.2", versionType: "", e: errors.New("invalid format"), output: ""},
		{input: "1234.1234.1234", versionType: "", e: nil, output: "1234.1234.1235"},
		{input: "v1.8.58-5-g98e9b2b", versionType: "", e: nil, output: "v1.8.59"},

		// Develop branch
		{input: "0.1.0-5-g98e9b2b", versionType: "alpha", e: nil, output: "0.2.0-alpha.1"},
		{input: "0.2.0-alpha.1", versionType: "", e: nil, output: "0.2.0-alpha.2"},
		{input: "0.2.0-alpha.2", versionType: "", e: nil, output: "0.2.0-alpha.3"},
		{input: "0.2.0-alpha.3", versionType: "", e: nil, output: "0.2.0-alpha.4"},
		{input: "0.2.0-alpha.4", versionType: "beta", e: nil, output: "0.2.0-beta.1"},
		{input: "0.2.0-beta.1", versionType: "", e: nil, output: "0.2.0-beta.2"},
		{input: "0.2.0-beta.2", versionType: "", e: nil, output: "0.2.0-beta.3"},
		{input: "0.2.0-beta.2", versionType: "alpha", e: errors.New("regressive tagging"), output: ""},

		// Release branch
		{input: "0.2.0-beta.3", versionType: "rc", e: nil, output: "0.2.0-rc.1"},
		{input: "0.2.0-rc.1", versionType: "", e: nil, output: "0.2.0-rc.2"},
		{input: "0.2.0-rc.2", versionType: "", e: nil, output: "0.2.0-rc.3"},
		{input: "0.2.0-rc.2", versionType: "alpha", e: errors.New("regressive tagging"), output: ""},
		{input: "0.2.0-rc.2", versionType: "beta", e: errors.New("regressive tagging"), output: ""},

		// Master branch
		{input: "0.2.0-rc.3", versionType: "release", e: nil, output: "0.2.0"},
		{input: "0.2.0", versionType: "", e: nil, output: "0.2.1"},
		{input: "0.2.1", versionType: "", e: nil, output: "0.2.2"},
		{input: "0.2.2", versionType: "", e: nil, output: "0.2.3"},
	}

	for k := range tests {

		test := tests[k]
		nextVersion, e := versioner.NextVersion(test.input, test.versionType)

		if test.e == nil {
			require.Nil(t, e, fmt.Sprintf("For #%d %s - %s", k, test.input, test.versionType))
			assert.Equal(t, test.output, nextVersion, fmt.Sprintf("For #%d %s - %s", k, test.input, test.versionType))
		} else {
			require.NotNil(t, e, fmt.Sprintf("For #%d %s - %s", k, test.input, test.versionType))
			assert.Equal(t, test.e.Error(), e.Error(), fmt.Sprintf("For #%d %s - %s", k, test.input, test.versionType))
			assert.Equal(t, test.output, nextVersion, fmt.Sprintf("For #%d %s - %s", k, test.input, test.versionType))
		}

	}

}

func TestParsePreRelease(t *testing.T) {

	tests := []struct {
		input       string
		versionType string
		output      string
		outputNum   int64
	}{
		{input: "1.2.3", versionType: "alpha", output: "1.2.3", outputNum: 0},
		{input: "1.2.3-12345", versionType: "alpha", output: "1.2.3-12345", outputNum: 0},
		{input: "1.2.3-alpha.1", versionType: "alpha", output: "1.2.3", outputNum: 1},
		{input: "1.2.3-beta.1", versionType: "beta", output: "1.2.3", outputNum: 1},
		{input: "1.2.3-beta.1", versionType: "alpha", output: "1.2.3-beta.1", outputNum: 0},
		{input: "1.2.3-rc.2", versionType: "rc", output: "1.2.3", outputNum: 2},
	}

	for k := range tests {
		test := tests[k]

		output, outputNum := versioner.ParsePreRelease(test.input, test.versionType)

		assert.Equal(t, test.output, output, "For #%d %s (%s) --> output `%s`", k, test.input, test.versionType, output)
		assert.Equal(t, test.outputNum, outputNum, "For #%d %s (%s) --> outputNum %d", k, test.input, test.versionType, outputNum)
	}
}
