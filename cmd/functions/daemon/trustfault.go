//go:build !ci

package daemon

func injectedTrustFault(string) error { return nil }
