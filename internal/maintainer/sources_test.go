package maintainer

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func writeTar(t *testing.T, p string, files map[string]string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var w io.Writer = f
	if filepath.Ext(p) == ".gz" {
		zw := gzip.NewWriter(f)
		defer zw.Close()
		w = zw
	}
	tw := tar.NewWriter(w)
	defer tw.Close()
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		if err := tw.WriteHeader(&tar.Header{Name: n, Mode: 0o644, Size: int64(len(files[n])), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(files[n])); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSourcesHoldFilesAndArchiveMembers(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.so"), "plain file content", 0o644)
	writeTar(t, filepath.Join(dir, "x.tar"), map[string]string{"x/in-tar": "member of a plain tar"})
	writeTar(t, filepath.Join(dir, "sub", "y.tar.gz"), map[string]string{"y/in-tgz": "member of a gzipped tar"})
	s := newSources(dir)
	for _, c := range []string{"plain file content", "member of a plain tar", "member of a gzipped tar"} {
		if ok, err := s.has(int64(len(c)), sum(c)); err != nil || !ok {
			t.Errorf("%q: %v %v", c, ok, err)
		}
	}
	if ok, err := s.has(5, sum("other")); err != nil || ok {
		t.Errorf("content the lib_dir does not hold: %v %v", ok, err)
	}
}

func TestArchiveMembersIgnoresNonArchives(t *testing.T) {
	p := filepath.Join(t.TempDir(), "f")
	write(t, p, string(make([]byte, 1024)), 0o644)
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if m, err := archiveMembers(f); err != nil || len(m) != 0 {
		t.Fatalf("%v %v", m, err)
	}
}
