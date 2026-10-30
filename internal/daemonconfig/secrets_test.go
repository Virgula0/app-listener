package daemonconfig

import (
	"slices"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

const secretService = "[secret_service]\nprovider /usr/bin/true\nprompter /usr/bin/false\nmanager /usr/bin/env\n\n"

func TestParseSecretSections(t *testing.T) {
	cfg, err := Parse([]byte(secretService + "[secret gh]\nmatch service=gh:*\nmatch xdg:schema=x application=gh\n" +
		"/usr/bin/true READ\n/usr/bin/false\n"))
	if err != nil {
		t.Fatal(err)
	}
	svc := cfg.SecretService
	if !slices.Equal(svc.Providers, []string{"/usr/bin/true"}) || !slices.Equal(svc.Prompters, []string{"/usr/bin/false"}) ||
		!slices.Equal(svc.Managers, []string{"/usr/bin/env"}) {
		t.Fatalf("secret_service: %+v", svc)
	}
	if len(cfg.Resources) != 2 {
		t.Fatalf("want the scope and the relay, got %+v", cfg.Resources)
	}
	scope, relay := cfg.Resources[0], cfg.Resources[1]
	if scope.Secret == nil || scope.Secret.Name != "gh" || scope.Path != SecretAnchor("gh") || scope.NeedEncryption ||
		len(scope.Secret.Matches) != 2 || len(scope.Binaries) != 2 || scope.Binaries[0].Events[0] != ebpf.EventRead {
		t.Fatalf("scope: %+v", scope)
	}
	if !relay.SecretRelay || relay.Path != SecretRelayDir || len(relay.Binaries) != 1 ||
		!slices.Equal(relay.Binaries[0].Events, []ebpf.EventType{ebpf.EventConnect, ebpf.EventStat}) {
		t.Fatalf("relay: %+v", relay)
	}
}

func TestSecretMatch(t *testing.T) {
	m := SecretMatch{"service": "gh:*", "username": "octo"}
	for attrs, want := range map[*map[string]string]bool{
		{"service": "gh:github.com", "username": "octo"}:          true,
		{"service": "gh:", "username": "octo", "extra": "x"}:      true,
		{"service": "github.com", "username": "octo"}:             false,
		{"service": "gh:github.com"}:                              false,
		{"service": "gh:github.com", "username": "octocat"}:       false,
		{"Service": "gh:github.com", "username": "octo"}:          false,
		{"service": "gh:github.com", "username": "octo*"}:         false,
		{"service": "xgh:github.com", "username": "octo", "": ""}: false,
	} {
		if got := m.Matches(*attrs); got != want {
			t.Errorf("%v: got %v", *attrs, got)
		}
	}
}

func TestParseSecretSectionRefusals(t *testing.T) {
	scope := "[secret gh]\nmatch service=gh:*\n/usr/bin/true\n"
	for name, text := range map[string]string{
		"no provider":      scope,
		"prompter only":    "[secret_service]\nprompter /usr/bin/true\n",
		"no match":         secretService + "[secret gh]\n/usr/bin/true\n",
		"bare star":        secretService + "[secret gh]\nmatch service=*\n",
		"inner star":       secretService + "[secret gh]\nmatch service=g*h\n",
		"no value":         secretService + "[secret gh]\nmatch service=\n",
		"no pair":          secretService + "[secret gh]\nmatch service\n",
		"dup attribute":    secretService + "[secret gh]\nmatch a=1 a=2\n",
		"file event":       secretService + "[secret gh]\nmatch a=1\n/usr/bin/true STAT\n",
		"connect event":    secretService + "[secret gh]\nmatch a=1\n/usr/bin/true CONNECT\n",
		"watch path":       secretService + "[secret gh]\nmatch a=1\nwatch: /tmp\n",
		"encrypted":        secretService + "[secret gh]\nmatch a=1\nneed_encryption: true\n",
		"runtime":          secretService + "[secret gh]\nmatch a=1\nruntime: true\n",
		"bad name":         secretService + "[secret GH]\nmatch a=1\n",
		"slash in name":    secretService + "[secret a/b]\nmatch a=1\n",
		"duplicate":        secretService + scope + scope,
		"relative path":    "[secret_service]\nprovider gnome-keyring-daemon\n",
		"unknown":          "[secret_service]\nkeyring /usr/bin/true\n",
		"dup provider":     "[secret_service]\nprovider /usr/bin/true\nprovider /usr/bin/true\n",
		"provider + event": "[secret_service]\nprovider /usr/bin/true CONNECT\n",
	} {
		if _, err := Parse([]byte(text)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}
