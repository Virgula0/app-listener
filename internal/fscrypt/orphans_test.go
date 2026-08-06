package fscrypt

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/fscrypt/metadata"
)

// protector helpers build synthetic raw-key and login protectors the same
// way the protobuf metadata is laid out on disk.
func appProtector(desc string) *metadata.ProtectorData {
	return &metadata.ProtectorData{
		ProtectorDescriptor: desc,
		Source:              metadata.SourceType_raw_key,
		Name:                appProtectorPrefix + "1730000000000000000",
	}
}

func foreignProtector(desc string) *metadata.ProtectorData {
	return &metadata.ProtectorData{
		ProtectorDescriptor: desc,
		Source:              metadata.SourceType_pam_passphrase,
		Name:                "login-protector",
	}
}

func policy(desc string, protectorDescs ...string) *metadata.PolicyData {
	p := &metadata.PolicyData{KeyDescriptor: desc}
	for _, d := range protectorDescs {
		p.WrappedPolicyKeys = append(p.WrappedPolicyKeys, &metadata.WrappedPolicyKey{ProtectorDescriptor: d})
	}
	return p
}

// TestSelectOrphansLiveKept exercises the core safety property: metadata of a
// directory that still exists is never flagged for deletion, and foreign
// protectors/policies are untouched no matter what.
func TestSelectOrphansLiveKept(t *testing.T) {
	liveProt := appProtector("p-live")
	livePol := policy("pol-live", "p-live")
	protectors := []*metadata.ProtectorData{liveProt, foreignProtector("p-foreign")}
	policies := []*metadata.PolicyData{
		livePol,
		policy("foreign-pol", "p-foreign"),
		// ~/.ssh style: encrypted, exists, but not part of the config.
	}
	orphanProts, orphanPols := SelectOrphans(protectors, policies, map[string]struct{}{
		"pol-live": {},
	})
	if len(orphanProts) != 0 || len(orphanPols) != 0 {
		t.Errorf("SelectOrphans deleted live metadata: protectors=%v policies=%v", orphanProts, orphanPols)
	}
}

// TestSelectOrphansDeletedDirs exercises the other important behavior: metadata
// of deleted directories is removed, and only the metadata of this installer.
func TestSelectOrphansDeletedDirs(t *testing.T) {
	orphanProt := appProtector("p-gone")
	orphanPol := policy("pol-gone", "p-gone")
	// A live dir also exists ("~/.ssh") and is not referenced by the config.
	protectors := []*metadata.ProtectorData{
		orphanProt,
		appProtector("p-ssh"),
		foreignProtector("foreign"),
	}
	policies := []*metadata.PolicyData{
		orphanPol,
		policy("pol-ssh", "p-ssh"),
		policy("pol-foreign", "foreign"),
	}
	orphanProts, orphanPols := SelectOrphans(protectors, policies, map[string]struct{}{
		"pol-ssh": {},
	})

	if !reflect.DeepEqual(orphanProts, []string{"p-gone"}) {
		t.Errorf("orphan protectors = %v, want [p-gone]", orphanProts)
	}
	if !reflect.DeepEqual(orphanPols, []string{"pol-gone"}) {
		t.Errorf("orphan policies = %v, want [pol-gone]", orphanPols)
	}
}

// TestSelectOrphansEncryptedButUnmanaged verifies a directory that is still
// on disk but no longer in the configuration keeps its metadata: only
// directories that no longer exist are cleaned.
func TestSelectOrphansEncryptedButUnmanaged(t *testing.T) {
	protectors := []*metadata.ProtectorData{appProtector("p-unmanaged")}
	policies := []*metadata.PolicyData{policy("pol-unmanaged", "p-unmanaged")}
	// The live set is built by probing the catalog and final config: an
	// existing encrypted directory outside both is still found by the probe.
	live := map[string]struct{}{"pol-unmanaged": {}}
	orphanProts, orphanPols := SelectOrphans(protectors, policies, live)
	if len(orphanProts) != 0 || len(orphanPols) != 0 {
		t.Errorf("SelectOrphans deleted unmanaged-but-existing dir metadata: protectors=%v policies=%v", orphanProts, orphanPols)
	}
}

// TestSelectOrphansForeignUntouched verifies login protectors and policies
// referencing them survive even when they are not live.
func TestSelectOrphansForeignUntouched(t *testing.T) {
	protectors := []*metadata.ProtectorData{foreignProtector("foreign")}
	policies := []*metadata.PolicyData{policy("pol-foreign", "foreign")}
	orphanProts, orphanPols := SelectOrphans(protectors, policies, nil)
	if len(orphanProts) != 0 || len(orphanPols) != 0 {
		t.Errorf("foreign fscrypt metadata must never be deleted: protectors=%v policies=%v", orphanProts, orphanPols)
	}
}

// TestSelectOrphansSortedAndDedup verifies the orphan lists come back sorted
// (the removal order does not depend on the metadata directory order).
func TestSelectOrphansSortedAndDedup(t *testing.T) {
	protectors := []*metadata.ProtectorData{
		appProtector("zzz"),
		appProtector("aaa"),
	}
	policies := []*metadata.PolicyData{
		policy("pol-aaa", "aaa"),
		policy("pol-zzz", "zzz"),
	}
	orphanProts, orphanPols := SelectOrphans(protectors, policies, nil)
	if !reflect.DeepEqual(orphanProts, []string{"aaa", "zzz"}) {
		t.Errorf("orphan protectors = %v, want sorted [aaa zzz]", orphanProts)
	}
	if !reflect.DeepEqual(orphanPols, []string{"pol-aaa", "pol-zzz"}) {
		t.Errorf("orphan policies = %v, want sorted [pol-aaa pol-zzz]", orphanPols)
	}
}

// TestLivePolicyDescriptorsTmpfs verifies the probe walks the given paths and
// returns an empty set on a filesystem that cannot carry fscrypt policies
// (tmpfs), skipping directories and non-existing paths without error.
func TestLivePolicyDescriptorsTmpfs(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	live, err := LivePolicyDescriptors([]string{dir, sub, filepath.Join(dir, "does-not-exist")})
	if err != nil {
		t.Fatalf("LivePolicyDescriptors: %v", err)
	}
	if len(live) != 0 {
		t.Errorf("tmpfs should yield no policies, got %v", live)
	}
}
