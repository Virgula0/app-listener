package certificates

import _ "embed"

//go:embed app-listener-release.pub
var ReleasePublicKeyPEM []byte
