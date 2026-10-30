package secretgate

import (
	"github.com/godbus/dbus/v5"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

const (
	ifaceService    = "org.freedesktop.Secret.Service"
	ifaceCollection = "org.freedesktop.Secret.Collection"
	ifaceItem       = "org.freedesktop.Secret.Item"
	ifaceSession    = "org.freedesktop.Secret.Session"
	ifacePrompt     = "org.freedesktop.Secret.Prompt"
	ifaceProps      = "org.freedesktop.DBus.Properties"
	// gnome-keyring's own: gcr's prompter answers through the first, Seahorse and the PAM module
	// create, unlock and re-key keyrings with a master password through the second.
	ifacePrompterCallback = "org.gnome.keyring.internal.Prompter.Callback"
	ifaceGuiltRidden      = "org.gnome.keyring.InternalUnsupportedGuiltRiddenInterface"

	propAttributes = "org.freedesktop.Secret.Item.Attributes"
)

// openCalls pass undecided: {interface, member}, member "" for the whole interface (a call always
// names its member).
var openCalls = map[[2]string]bool{
	{ifaceSession, "Close"}:                     true,
	{ifacePrompt, "Prompt"}:                     true,
	{ifacePrompt, "Dismiss"}:                    true,
	{"org.freedesktop.DBus.Introspectable", ""}: true,
	{"org.freedesktop.DBus.Peer", ""}:           true,
}

type verdictKind int

const (
	refuse verdictKind = iota
	allow
	needManager
	needPrompter
	// gated: the call touches items; each item's scopes must admit the caller for op.
	gated
)

// verdict is what classify decides from the message alone. For gated, items are the item paths
// whose current attributes decide, attrs the attribute sets the call would give an item.
type verdict struct {
	kind  verdictKind
	op    ebpf.EventType
	items []dbus.ObjectPath
	attrs []map[string]string
}

// classify maps a method call to the keyring onto the gatekeeper's policy. Anything not listed is
// refused: gnome-keyring grows interfaces (the Flatpak portal's RetrieveSecret hands out per-app
// secrets to any caller naming an app id) that this table must opt into, never inherit.
func classify(m *dbus.Message) verdict {
	iface := header[string](m, dbus.FieldInterface)
	member := header[string](m, dbus.FieldMember)
	path := header[dbus.ObjectPath](m, dbus.FieldPath)
	if member != "" && (openCalls[[2]string{iface, member}] || openCalls[[2]string{iface, ""}]) {
		return verdict{kind: allow}
	}
	switch iface {
	case ifaceService:
		return classifyService(m, member)
	case ifaceCollection:
		return classifyCollection(m, member)
	case ifaceItem:
		switch member {
		case "GetSecret":
			return verdict{kind: gated, op: ebpf.EventRead, items: []dbus.ObjectPath{path}}
		case "SetSecret":
			return verdict{kind: gated, op: ebpf.EventWrite, items: []dbus.ObjectPath{path}}
		case "Delete":
			return verdict{kind: gated, op: ebpf.EventDelete, items: []dbus.ObjectPath{path}}
		}
	case ifaceProps:
		return classifyProps(m, member, path)
	case ifacePrompterCallback:
		return verdict{kind: needPrompter}
	case ifaceGuiltRidden:
		return verdict{kind: needManager}
	}
	return verdict{kind: refuse}
}

func classifyService(m *dbus.Message, member string) verdict {
	switch member {
	case "OpenSession", "SearchItems", "Unlock", "Lock", "LockService", "ReadAlias", "CreateCollection":
		return verdict{kind: allow}
	case "GetSecrets":
		items, ok := arg[[]dbus.ObjectPath](m, 0)
		if !ok {
			return verdict{kind: refuse}
		}
		return verdict{kind: gated, op: ebpf.EventRead, items: items}
	case "SetAlias", "ChangeLock":
		return verdict{kind: needManager}
	}
	return verdict{kind: refuse}
}

func classifyCollection(m *dbus.Message, member string) verdict {
	switch member {
	case "SearchItems":
		return verdict{kind: allow}
	case "CreateItem":
		props, ok := arg[map[string]dbus.Variant](m, 0)
		if !ok {
			return verdict{kind: refuse}
		}
		attrs, ok := attributesOf(props[propAttributes])
		if !ok {
			return verdict{kind: refuse}
		}
		// replace=true overwrites the item with these attributes: the same scopes decide.
		return verdict{kind: gated, op: ebpf.EventWrite, attrs: []map[string]string{attrs}}
	case "Delete":
		return verdict{kind: needManager}
	}
	return verdict{kind: refuse}
}

// classifyProps: reading properties is metadata, open to all; an item's label or attributes are
// written under the scopes of both what it is and what it would become (re-attributing an item must
// not move it into or out of a scope).
func classifyProps(m *dbus.Message, member string, path dbus.ObjectPath) verdict {
	switch member {
	case "Get", "GetAll":
		return verdict{kind: allow}
	case "Set":
	default:
		return verdict{kind: refuse}
	}
	iface, ok := arg[string](m, 0)
	if !ok {
		return verdict{kind: refuse}
	}
	switch iface {
	case ifaceItem:
		v := verdict{kind: gated, op: ebpf.EventWrite, items: []dbus.ObjectPath{path}}
		if prop, _ := arg[string](m, 1); prop == "Attributes" {
			val, ok := arg[dbus.Variant](m, 2)
			if !ok {
				return verdict{kind: refuse}
			}
			attrs, ok := attributesOf(val)
			if !ok {
				return verdict{kind: refuse}
			}
			v.attrs = []map[string]string{attrs}
		}
		return v
	case ifaceCollection, ifaceService:
		return verdict{kind: needManager}
	}
	return verdict{kind: refuse}
}

// attributesOf reads an a{ss} variant; an absent one (zero Variant) is no attributes.
func attributesOf(v dbus.Variant) (map[string]string, bool) {
	if v.Value() == nil {
		return map[string]string{}, true
	}
	attrs, ok := v.Value().(map[string]string)
	return attrs, ok
}

func arg[T any](m *dbus.Message, i int) (T, bool) {
	var zero T
	if i >= len(m.Body) {
		return zero, false
	}
	t, ok := m.Body[i].(T)
	return t, ok
}
