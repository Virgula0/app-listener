package editprotected

import "testing"

func TestValidateForwardFlags(t *testing.T) {
	reset := func() {
		forwardFlag, yesFlag, setPasswordFlag = false, false, false
		whitelistFlags, blacklistFlags, eventsFlags = nil, nil, nil
		putFlag = ""
	}
	t.Cleanup(reset)
	cases := []struct {
		name string
		set  func()
		ok   bool
	}{
		{"no flags", func() {}, true},
		{"-w without --forward", func() { whitelistFlags = []string{"/bin/x"} }, false},
		{"-e without --forward", func() { eventsFlags = []string{"OPEN"} }, false},
		{"--forward alone", func() { forwardFlag = true }, false},
		{"--forward -w", func() { forwardFlag, whitelistFlags = true, []string{"/bin/x"} }, true},
		{"--forward -b -e", func() {
			forwardFlag, blacklistFlags, eventsFlags = true, []string{"/bin/x"}, []string{"write"}
		}, true},
		{"--forward -w -b", func() {
			forwardFlag, whitelistFlags, blacklistFlags = true, []string{"/bin/x"}, []string{"/bin/y"}
		}, false},
		{"--forward bad -e", func() {
			forwardFlag, whitelistFlags, eventsFlags = true, []string{"/bin/x"}, []string{"EXEC"}
		}, false},
		{"--forward --put", func() { forwardFlag, whitelistFlags, putFlag = true, []string{"/bin/x"}, "f" }, false},
		{"--forward --set-password", func() {
			forwardFlag, whitelistFlags, setPasswordFlag = true, []string{"/bin/x"}, true
		}, false},
	}
	for _, c := range cases {
		reset()
		c.set()
		if err := validateForwardFlags(); (err == nil) != c.ok {
			t.Errorf("%s: err = %v, want ok=%v", c.name, err, c.ok)
		}
	}
}
