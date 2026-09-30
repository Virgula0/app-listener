package integrationtests

// The edit-protected editor runs as root over a user-owned vault. Its root-only checks live in
// internal/tui/harness_test.go (tuiTestAmd64Bin); these run them in a pooled container.

func (s *IntegrationSuite) runEditorHarness(subtest string) {
	c := s.guardContainer()
	// pooled: terminated at suite end
	s.Require().NoError(c.CopyFileToContainer(s.ctx, tuiTestAmd64Bin, "/tui.test", 0o755))
	code, out := s.exec(c, []string{"sh", "-c",
		"APPLISTENER_TUI_HARNESS=1 /tui.test -test.v -test.run " + shQuote("^TestEditorHarness/"+subtest+"$")})
	s.Require().Equalf(0, code, "tui.test %s failed: %s", subtest, out)
	s.Require().Containsf(out, "--- PASS: TestEditorHarness/"+subtest, "tui.test %s did not run: %s", subtest, out)
}

// Offline edit-protected: the vault's user swaps a directory inside it for a symlink while the
// editor is open; a chown must not hand the link target (a system file) to that user.
func (s *IntegrationSuite) TestEditor_ChownRefusesSymlinkedParent() {
	s.runEditorHarness("TestChownRefusesSymlinkedParent")
}

func (s *IntegrationSuite) TestEditor_ChownInsideVault() {
	s.runEditorHarness("TestChownInsideVault")
}
