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
		{input: "v1.2.3-RC1", versionType: "rc", e: nil, output: "v1.2.3-RC2"},
		{input: "1.2.3-RC1", versionType: "rc", e: nil, output: "1.2.3-RC2"},
		{input: "1.2.3-RC2", versionType: "release", e: nil, output: "1.2.3"},
		{input: "v1.2.3", versionType: "rc", e: nil, output: "v1.3.0-RC1"},
		{input: "1.2.3", versionType: "rc", e: nil, output: "1.3.0-RC1"},
		{input: "1.2", versionType: "", e: errors.New("invalid format"), output: ""},
		{input: "1234.1234.1234", versionType: "", e: nil, output: "1234.1234.1235"},
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
