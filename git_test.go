package main

import (
	"testing"
)

func TestVerifyCommitMessage(t *testing.T) {
	testCases := []struct {
		name  string
		msg   string
		valid bool
	}{
		{
			msg:   "feat: valid message",
			valid: true,
		},
		{
			msg:   "feat:",
			valid: false,
		},
		{
			msg:   "invalid",
			valid: false,
		},
		{
			msg:   "feat commit message",
			valid: false,
		},
		{
			msg:   "feat",
			valid: false,
		},
	}

	for _, tc := range testCases {
		res := verifyCommitMessage(tc.msg)

		if tc.valid != res {
			t.Fail()
		}
	}
}
