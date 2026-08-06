package fscrypt

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/filesystem"
	"github.com/google/fscrypt/metadata"
)

// appProtectorPrefix is the name given to the raw-key protector created for
// every directory migrated by this installer (see applyRawKeyPolicy). It is
// the signature that marks fscrypt metadata as ours: only metadata carrying
// this name (and a raw-key source) is ever touched by the orphan cleanup, so
// login protectors and policies created by other tools survive it.
const appProtectorPrefix = "app-listener-key-"

// LivePolicyDescriptors returns the set of policy descriptors attached to the
// given paths. Only paths that exist right now can contribute: a directory
// that was deleted is absent, so its descriptor is not returned and its
// metadata is treated as orphaned. The paths are the catalog's watch
// directories for every local user plus the final config resources — the
// installer encrypts at exactly these levels and never anywhere else, so
// this set is complete for everything the installer ever created.
func LivePolicyDescriptors(paths []string) (map[string]struct{}, error) {
	live := make(map[string]struct{})
	for _, p := range paths {
		data, err := metadata.GetPolicy(p)
		switch {
		case err == nil:
			live[data.GetKeyDescriptor()] = struct{}{}
		case errors.Is(err, metadata.ErrEncryptionNotSupported),
			errors.Is(err, metadata.ErrEncryptionNotEnabled):
			// The filesystem cannot carry fscrypt policies at all.
		default:
			var notEncrypted *metadata.ErrNotEncrypted
			var locked *metadata.ErrLockedRegularFile
			switch {
			case errors.As(err, &notEncrypted):
				// Existing directory without a policy.
			case errors.As(err, &locked):
				// A locked regular file, which the installer never manages.
			case os.IsNotExist(err):
				// Deleted since the catalog probe: not live.
			default:
				return nil, fmt.Errorf("get fscrypt policy of %s: %w", p, err)
			}
		}
	}
	return live, nil
}

// isAppRawRawKeyProtector reports whether a protector belongs to this
// installer: a raw-key protector (never a login/PAM one) carrying the
// app-listener-key-* name prefix.
func isAppProtector(p *metadata.ProtectorData) bool {
	return p.GetSource() == metadata.SourceType_raw_key &&
		strings.HasPrefix(p.GetName(), appProtectorPrefix)
}

// SelectOrphans partitions the fscrypt metadata of one mount into entries to
// keep and entries that may be deleted. An orphan protector is one that no
// policy attached to a still-existing directory references; an orphan policy
// is one whose directory no longer exists and that only uses protectors
// created by this installer. Deleted status differs from "not in live": a
// directory that still exists keeps its metadata even when it is not part of
// the current configuration (e.g. an independently encrypted ~/.ssh). The two
// returned slices are sorted and deduplicated.
func SelectOrphans(protectors []*metadata.ProtectorData, policies []*metadata.PolicyData, livePolicies map[string]struct{}) (orphanProtectors, orphanPolicies []string) {
	refs := make(map[string][]string, len(policies))
	for _, p := range policies {
		var desc []string
		for _, k := range p.GetWrappedPolicyKeys() {
			desc = append(desc, k.GetProtectorDescriptor())
		}
		refs[p.GetKeyDescriptor()] = desc
	}

	orphanProtectors = orphanedProtectors(protectors, policies, refs, livePolicies)
	orphanPolicies = orphanedPolicies(policies, refs, appProtectorSet(protectors), livePolicies)
	sort.Strings(orphanProtectors)
	sort.Strings(orphanPolicies)
	return orphanProtectors, orphanPolicies
}

// appProtectorSet returns the descriptors of every protector belonging to
// this installer (raw-key source, app-listener-key-* name).
func appProtectorSet(protectors []*metadata.ProtectorData) map[string]struct{} {
	set := make(map[string]struct{})
	for _, p := range protectors {
		if isAppProtector(p) {
			set[p.GetProtectorDescriptor()] = struct{}{}
		}
	}
	return set
}

// orphanedProtectors is the set of app protectors that no live policy
// references anymore, i.e. every directory that could unlock with them is
// gone.
func orphanedProtectors(protectors []*metadata.ProtectorData, policies []*metadata.PolicyData, refs map[string][]string, livePolicies map[string]struct{}) []string {
	live := make(map[string]struct{})
	for _, p := range policies {
		if _, ok := livePolicies[p.GetKeyDescriptor()]; !ok {
			continue
		}
		for _, d := range refs[p.GetKeyDescriptor()] {
			live[d] = struct{}{}
		}
	}
	var orphan []string
	for _, p := range protectors {
		if !isAppProtector(p) {
			continue
		}
		if _, ok := live[p.GetProtectorDescriptor()]; !ok {
			orphan = append(orphan, p.GetProtectorDescriptor())
		}
	}
	return orphan
}

// orphanedPolicies is the set of policies whose directory no longer exists
// and that reference only protectors created by this installer.
func orphanedPolicies(policies []*metadata.PolicyData, refs map[string][]string, appProtectors, livePolicies map[string]struct{}) []string {
	var orphan []string
	for _, p := range policies {
		if _, ok := livePolicies[p.GetKeyDescriptor()]; ok {
			continue
		}
		protectorDescs := refs[p.GetKeyDescriptor()]
		if len(protectorDescs) == 0 {
			continue
		}
		appOnly := true
		for _, d := range protectorDescs {
			if _, ok := appProtectors[d]; !ok {
				appOnly = false
				break
			}
		}
		if appOnly {
			orphan = append(orphan, p.GetKeyDescriptor())
		}
	}
	return orphan
}

// CleanOrphans deletes the fscrypt policy and protector metadata on the mount
// containing anchor that SelectOrphans declared orphaned. It only ever
// removes raw-key app-listener-key-* protector metadata and the policy
// metadata that references only it, never login protectors or foreign
// setups. A mount that carries no fscrypt metadata at all is a no-op. The
// returned counts are the protectors and policies actually removed.
func CleanOrphans(anchor string, livePolicies map[string]struct{}) (removedProtectors, removedPolicies int, err error) {
	ctx, ctxErr := actions.NewContextFromPath(anchor, nil)
	if ctxErr != nil {
		return 0, 0, fmt.Errorf("fscrypt context for %s: %w", anchor, ctxErr)
	}
	if _, statErr := os.Stat(ctx.Mount.ProtectorDir()); statErr != nil {
		if os.IsNotExist(statErr) {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("stat fscrypt metadata on %s: %w", ctx.Mount, statErr)
	}

	protectors, err := listProtectors(ctx)
	if err != nil {
		return 0, 0, err
	}
	policies, err := listPolicies(ctx)
	if err != nil {
		return 0, 0, err
	}

	orphanProtectors, orphanPolicies := SelectOrphans(protectors, policies, livePolicies)
	for _, d := range orphanPolicies {
		if err := ctx.Mount.RemovePolicy(d); err != nil {
			var notFound *filesystem.ErrPolicyNotFound
			if errors.As(err, &notFound) {
				continue
			}
			return removedProtectors, removedPolicies, fmt.Errorf("remove policy %s: %w", d, err)
		}
		removedPolicies++
	}
	for _, d := range orphanProtectors {
		if err := ctx.Mount.RemoveProtector(d); err != nil {
			var notFound *filesystem.ErrProtectorNotFound
			if errors.As(err, &notFound) {
				continue
			}
			return removedProtectors, removedPolicies, fmt.Errorf("remove protector %s: %w", d, err)
		}
		removedProtectors++
	}
	return removedProtectors, removedPolicies, nil
}

// listProtectors reads every protector's metadata on the context's mount.
func listProtectors(ctx *actions.Context) ([]*metadata.ProtectorData, error) {
	descriptors, err := ctx.Mount.ListProtectors(nil)
	if err != nil {
		return nil, fmt.Errorf("listing fscrypt protectors on %s: %w", ctx.Mount, err)
	}
	protectors := make([]*metadata.ProtectorData, 0, len(descriptors))
	for _, d := range descriptors {
		_, data, getErr := ctx.Mount.GetProtector(d, nil)
		if getErr != nil {
			var notFound *filesystem.ErrProtectorNotFound
			if errors.As(getErr, &notFound) {
				continue
			}
			return nil, fmt.Errorf("get protector %s: %w", d, getErr)
		}
		protectors = append(protectors, data)
	}
	return protectors, nil
}

// listPolicies loads every policy's listed metadata on the context's mount.
func listPolicies(ctx *actions.Context) ([]*metadata.PolicyData, error) {
	descriptors, err := ctx.Mount.ListPolicies(nil)
	if err != nil {
		return nil, fmt.Errorf("listing fscrypt policies on %s: %w", ctx.Mount, err)
	}
	policies := make([]*metadata.PolicyData, 0, len(descriptors))
	for _, d := range descriptors {
		data, getErr := ctx.Mount.GetPolicy(d, nil)
		if getErr != nil {
			var notFound *filesystem.ErrPolicyNotFound
			if errors.As(getErr, &notFound) {
				continue
			}
			return nil, fmt.Errorf("get policy %s: %w", d, getErr)
		}
		policies = append(policies, data)
	}
	return policies, nil
}
