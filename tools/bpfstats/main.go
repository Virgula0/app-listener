// Command bpfstats loads each BPF program of an object file on its own and reports the verifier's
// "processed N insns" count — the figure the 1M limit applies to, which program size only loosely
// predicts. Development aid for keeping the guard programs inside the budget; needs root.
//
//	sudo ./bpfstats new.o [base.o ...]
package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
)

var processedRe = regexp.MustCompile(`processed (\d+) insns`)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: bpfstats <object.o> [object.o ...]")
		os.Exit(2)
	}
	if err := rlimit.RemoveMemlock(); err != nil {
		fmt.Fprintf(os.Stderr, "removing memlock rlimit (run as root): %v\n", err)
		os.Exit(1)
	}

	results := make([]map[string]int, 0, len(os.Args)-1)
	for _, path := range os.Args[1:] {
		fmt.Printf("== %s\n", path)
		results = append(results, measure(path))
	}
	if len(results) == 2 {
		compare(results[0], results[1])
	}
}

// measure loads every program in the object separately, so one rejection does not hide the rest.
func measure(path string) map[string]int {
	spec, err := cilium.LoadCollectionSpec(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading %s: %v\n", path, err)
		os.Exit(1)
	}

	out := make(map[string]int)
	names := make([]string, 0, len(spec.Programs))
	for name := range spec.Programs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		insns, loadErr := loadOne(spec, name)
		switch {
		case loadErr != nil:
			out[name] = -1
			fmt.Printf("  %-32s REJECTED  %s\n", name, short(loadErr))
		default:
			out[name] = insns
			fmt.Printf("  %-32s %9d insns\n", name, insns)
		}
	}
	return out
}

func loadOne(spec *cilium.CollectionSpec, name string) (int, error) {
	one := &cilium.CollectionSpec{
		Maps:     spec.Maps,
		Programs: map[string]*cilium.ProgramSpec{name: spec.Programs[name]},
	}
	coll, err := cilium.NewCollectionWithOptions(one, cilium.CollectionOptions{
		Programs: cilium.ProgramOptions{LogLevel: cilium.LogLevelStats},
	})
	if err != nil {
		return 0, err
	}
	defer coll.Close()

	m := processedRe.FindStringSubmatch(coll.Programs[name].VerifierLog)
	if m == nil {
		return 0, nil
	}
	n, _ := strconv.Atoi(m[1])
	return n, nil
}

// compare prints the delta between the first two objects, worst regression first.
func compare(newer, older map[string]int) {
	type row struct {
		name           string
		newVal, oldVal int
		delta          int
	}
	var rows []row
	for name, nv := range newer {
		ov, ok := older[name]
		if !ok {
			continue
		}
		rows = append(rows, row{name, nv, ov, nv - ov})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].delta > rows[j].delta })

	fmt.Printf("\n== delta (arg1 vs arg2), worst first\n")
	for _, r := range rows {
		fmt.Printf("  %-32s %9d -> %9d  %+d\n", r.name, r.oldVal, r.newVal, r.delta)
	}
}

func short(err error) string {
	var ve *cilium.VerifierError
	if errors.As(err, &ve) {
		return fmt.Sprintf("%v", ve)
	}
	s := err.Error()
	if len(s) > 120 {
		s = s[:120]
	}
	return s
}
