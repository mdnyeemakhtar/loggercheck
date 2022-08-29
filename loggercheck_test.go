package loggercheck_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/timonwong/loggercheck"
	"github.com/timonwong/loggercheck/internal/rules"
)

type tracingTestingErrorf struct {
	*testing.T
}

func (t tracingTestingErrorf) Errorf(format string, args ...interface{}) {
	t.Logf(format, args...)
}

func TestLinter(t *testing.T) {
	testdata := analysistest.TestData()

	testCases := []struct {
		name      string
		patterns  string
		flags     []string
		wantError error
	}{
		{
			name:     "all",
			patterns: "a/all",
		},
		{
			name:     "klogonly",
			patterns: "a/klogonly",
			flags:    []string{"-disable=logr,zap"},
		},
		{
			name:     "custom-only",
			patterns: "a/customonly",
			flags: []string{
				"-disable=logr",
				"-rulefile",
				"testdata/custom-rules.txt",
			},
		},
		{
			name:     "wrong-rules",
			patterns: "a/customonly",
			flags: []string{
				"-rulefile",
				"testdata/wrong-rules.txt",
			},
			wantError: rules.ErrInvalidRule,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := loggercheck.NewAnalyzer()
			err := a.Flags.Parse(tc.flags)
			require.NoError(t, err)

			var result []*analysistest.Result
			if tc.wantError != nil {
				result = analysistest.Run(&tracingTestingErrorf{t}, testdata, a, tc.patterns)
			} else {
				result = analysistest.Run(t, testdata, a, tc.patterns)
			}
			require.Len(t, result, 1)

			if tc.wantError != nil {
				assert.Error(t, result[0].Err)
				assert.True(t, errors.Is(result[0].Err, tc.wantError))
			}
		})
	}
}

func TestOptions(t *testing.T) {
	testdata := analysistest.TestData()

	customRules := []string{
		"(*a/customonly.Logger).Debugw",
		"(*a/customonly.Logger).Infow",
		"(*a/customonly.Logger).Warnw",
		"(*a/customonly.Logger).Errorw",
		"(*a/customonly.Logger).With",
		"a/customonly.Debugw",
		"a/customonly.Infow",
		"a/customonly.Warnw",
		"a/customonly.Errorw",
		"a/customonly.With",
	}

	testCases := []struct {
		name    string
		options []loggercheck.Option
	}{
		{
			name: "disable-all-then-enable-mylogger",
			options: []loggercheck.Option{
				loggercheck.WithDisable([]string{"klog", "logr", "zap"}),
				loggercheck.WithRules(customRules),
			},
		},
		{
			name: "ignore-logr",
			options: []loggercheck.Option{
				loggercheck.WithDisable([]string{"logr"}),
				loggercheck.WithRules(customRules),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := loggercheck.NewAnalyzer(tc.options...)
			analysistest.Run(t, testdata, a, "a/customonly")
		})
	}
}
