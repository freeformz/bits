package main

import "testing"

func TestTruncateGoVersion(t *testing.T) {
	for name, tc := range map[string]struct {
		input  string
		output string
	}{
		"goMAJ.MIN.PATCH": {
			input:  "go1.13.3",
			output: "go1.13",
		},
		"goMAJ.MIN": {
			input:  "go1.13",
			output: "go1.13",
		},
		"goMAJ.MINrc": {
			input:  "go1.13rc1",
			output: "go1.13",
		},
		"goMAJ.MINbeta": {
			input:  "go1.13beta1",
			output: "go1.13",
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if output := truncateGoVersion(tc.input); output != tc.output {
				t.Errorf("got %q, expected %q", output, tc.output)
			}
		})
	}
}
