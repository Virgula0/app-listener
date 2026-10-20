// Command codecheck reports cleanup candidates golangci-lint misses, as a `make lint` gate:
//
//	unused   package-level func/var/const/type with no use under any tag set (exported too)
//	samebody functions whose bodies are identical once locals and params are renamed
//	thin     package-level funcs that only forward their params, in order, to one call
//	reterr   `if err := f(); err != nil { return err }; return nil` (just `return f()`)
//	deadcode funcs x/tools' deadcode finds unreachable from every main and test, under every tag set
//
// Findings listed in allow.txt (`check<TAB>file:Name`, see there for why each stays) pass. Exits 1
// on any other finding. Own module so x/tools stays out of the product's go.mod.
//
//	go run -C tools/codecheck . -root ../.. -deadcode v0.52.0
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

// tagSets are the build configurations a symbol must be unused (or unreachable) in all of to be
// reported: one only another tag set reaches is live.
var tagSets = []string{"", "gui", "ci", "gui,e2e"}

type finding struct{ check, key, detail string }

func main() {
	root := flag.String("root", ".", "module root to check")
	allowPath := flag.String("allow", "allow.txt", "allowlist")
	deadcodeVersion := flag.String("deadcode", "", "golang.org/x/tools/cmd/deadcode version to run (empty: skip)")
	flag.Parse()
	abs, err := filepath.Abs(*root)
	if err != nil {
		fail(err)
	}
	allow, err := readAllow(*allowPath)
	if err != nil {
		fail(err)
	}
	files, err := parseSources(abs)
	if err != nil {
		fail(err)
	}
	found := slices.Concat(sameBody(files), thin(files), retErr(files))
	unused, err := unusedObjects(abs)
	if err != nil {
		fail(err)
	}
	found = append(found, unused...)
	if *deadcodeVersion != "" {
		dead, err := deadFuncs(abs, *deadcodeVersion)
		if err != nil {
			fail(err)
		}
		found = append(found, dead...)
	}

	bad := 0
	for _, f := range found {
		if allow[f.check+"\t"+f.key] {
			delete(allow, f.check+"\t"+f.key)
			continue
		}
		fmt.Printf("%s: %s\n", f.check, f.detail)
		bad++
	}
	stale := slices.Sorted(func(yield func(string) bool) {
		for k := range allow {
			if !yield(k) {
				return
			}
		}
	})
	for _, k := range stale {
		fmt.Printf("allow.txt: stale entry %q (no longer reported, remove it)\n", k)
		bad++
	}
	if bad > 0 {
		os.Exit(1)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "codecheck:", err)
	os.Exit(2)
}

func readAllow(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	allow := map[string]bool{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		check, key, ok := strings.Cut(line, "\t")
		if !ok {
			return nil, fmt.Errorf("%s: %q: want check<TAB>file:Name", path, line)
		}
		allow[check+"\t"+strings.TrimSpace(key)] = true
	}
	return allow, sc.Err()
}

type srcFile struct {
	rel  string
	file *ast.File
	fset *token.FileSet
}

// parseSources parses every non-test, non-generated .go file of the module at root.
func parseSources(root string) ([]srcFile, error) {
	fset := token.NewFileSet()
	var out []srcFile
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p == root {
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") || d.Name() == "graphify-out" {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(p, "go.mod")); err == nil {
				return filepath.SkipDir // another module
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") || strings.HasSuffix(p, "_bpf.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		out = append(out, srcFile{rel: rel, file: f, fset: fset})
		return nil
	})
	return out, err
}

func funcName(fd *ast.FuncDecl) string {
	if fd.Recv == nil || len(fd.Recv.List) == 0 {
		return fd.Name.Name
	}
	t := fd.Recv.List[0].Type
	if s, ok := t.(*ast.StarExpr); ok {
		t = s.X
	}
	if ix, ok := t.(*ast.IndexExpr); ok {
		t = ix.X
	}
	if id, ok := t.(*ast.Ident); ok {
		return id.Name + "." + fd.Name.Name
	}
	return fd.Name.Name
}

func render(n any) string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, token.NewFileSet(), n)
	return b.String()
}

func eachFunc(files []srcFile, fn func(f srcFile, fd *ast.FuncDecl)) {
	for _, f := range files {
		for _, d := range f.file.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
				fn(f, fd)
			}
		}
	}
}

// sameBody groups functions by signature + body with locals renamed positionally, so a renamed
// copy still matches. Bodies under 3 printed lines are ignored.
func sameBody(files []srcFile) []finding {
	groups := map[string][]string{}
	eachFunc(files, func(f srcFile, fd *ast.FuncDecl) {
		if len(fd.Body.List) < 2 {
			return
		}
		names := map[string]string{}
		rename := func(id *ast.Ident) {
			if _, ok := names[id.Name]; !ok {
				names[id.Name] = fmt.Sprintf("v%d", len(names))
			}
		}
		for _, fl := range fieldLists(fd) {
			for _, field := range fl.List {
				for _, id := range field.Names {
					rename(id)
				}
			}
		}
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.AssignStmt:
				if x.Tok == token.DEFINE {
					for _, l := range x.Lhs {
						if id, ok := l.(*ast.Ident); ok {
							rename(id)
						}
					}
				}
			case *ast.ValueSpec:
				for _, id := range x.Names {
					rename(id)
				}
			case *ast.RangeStmt:
				for _, e := range []ast.Expr{x.Key, x.Value} {
					if id, ok := e.(*ast.Ident); ok {
						rename(id)
					}
				}
			}
			return true
		})
		body := render(renamed(fd.Body, names))
		if strings.Count(body, "\n") < 3 {
			return
		}
		key := render(fd.Type) + "\n" + body
		groups[key] = append(groups[key], f.rel+":"+funcName(fd))
	})
	var out []finding
	for _, g := range groups {
		if len(g) < 2 {
			continue
		}
		slices.Sort(g)
		for _, m := range g {
			out = append(out, finding{"samebody", m, m + " has the same body as " + strings.Join(without(g, m), ", ")})
		}
	}
	return out
}

func fieldLists(fd *ast.FuncDecl) []*ast.FieldList {
	lists := []*ast.FieldList{fd.Type.Params}
	if fd.Recv != nil {
		lists = append(lists, fd.Recv)
	}
	if fd.Type.Results != nil {
		lists = append(lists, fd.Type.Results)
	}
	return lists
}

// renamed returns a copy of body (re-parsed from its printed form) with idents mapped through names.
func renamed(body *ast.BlockStmt, names map[string]string) ast.Node {
	src := "package p\nfunc _() " + render(body)
	f, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		return body
	}
	clone := f.Decls[0].(*ast.FuncDecl).Body
	ast.Inspect(clone, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			if nn, ok := names[id.Name]; ok {
				id.Name = nn
			}
		}
		return true
	})
	return clone
}

func without(s []string, x string) []string {
	return slices.DeleteFunc(slices.Clone(s), func(v string) bool { return v == x })
}

// thin flags package-level funcs whose whole body is one call taking exactly their params.
func thin(files []srcFile) []finding {
	var out []finding
	eachFunc(files, func(f srcFile, fd *ast.FuncDecl) {
		if fd.Recv != nil || len(fd.Body.List) != 1 {
			return
		}
		var call *ast.CallExpr
		switch s := fd.Body.List[0].(type) {
		case *ast.ReturnStmt:
			if len(s.Results) == 1 {
				call, _ = s.Results[0].(*ast.CallExpr)
			}
		case *ast.ExprStmt:
			call, _ = s.X.(*ast.CallExpr)
		}
		if call == nil || call.Ellipsis.IsValid() {
			return
		}
		var params []string
		for _, field := range fd.Type.Params.List {
			for _, n := range field.Names {
				params = append(params, n.Name)
			}
		}
		if len(call.Args) != len(params) {
			return
		}
		for i, a := range call.Args {
			if id, ok := a.(*ast.Ident); !ok || id.Name != params[i] {
				return
			}
		}
		key := f.rel + ":" + fd.Name.Name
		out = append(out, finding{"thin", key, key + " only forwards to " + render(call.Fun)})
	})
	return out
}

func isNil(e ast.Expr) bool {
	id, ok := e.(*ast.Ident)
	return ok && id.Name == "nil"
}

// retErr flags a trailing `if v := call(); v != nil { return v }; return nil` in a func with
// one unnamed result.
func retErr(files []srcFile) []finding {
	var out []finding
	for _, f := range files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			var body *ast.BlockStmt
			var ft *ast.FuncType
			switch x := n.(type) {
			case *ast.FuncDecl:
				body, ft = x.Body, x.Type
			case *ast.FuncLit:
				body, ft = x.Body, x.Type
			default:
				return true
			}
			if body == nil || len(body.List) < 2 || ft.Results == nil || len(ft.Results.List) != 1 ||
				len(ft.Results.List[0].Names) > 0 {
				return true
			}
			if ifs, ok := body.List[len(body.List)-2].(*ast.IfStmt); ok && isRetErr(ifs, body.List[len(body.List)-1]) {
				pos := f.fset.Position(ifs.Pos())
				key := fmt.Sprintf("%s:%d", f.rel, pos.Line)
				out = append(out, finding{"reterr", key, key + " can return the call directly"})
			}
			return true
		})
	}
	return out
}

func isRetErr(ifs *ast.IfStmt, last ast.Stmt) bool {
	ret, ok := last.(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 || !isNil(ret.Results[0]) || ifs.Else != nil || len(ifs.Body.List) != 1 {
		return false
	}
	as, ok := ifs.Init.(*ast.AssignStmt)
	if !ok || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
		return false
	}
	id, ok := as.Lhs[0].(*ast.Ident)
	if _, isCall := as.Rhs[0].(*ast.CallExpr); !ok || !isCall {
		return false
	}
	be, ok := ifs.Cond.(*ast.BinaryExpr)
	if !ok || be.Op != token.NEQ || !isNil(be.Y) {
		return false
	}
	if x, ok := be.X.(*ast.Ident); !ok || x.Name != id.Name {
		return false
	}
	r, ok := ifs.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(r.Results) != 1 {
		return false
	}
	x, ok := r.Results[0].(*ast.Ident)
	return ok && x.Name == id.Name
}

// unusedObjects lists package-level objects of this module never referenced under any tag set,
// tests included. Methods and fields are skipped: interfaces and layouts use them implicitly.
func unusedObjects(root string) ([]finding, error) {
	defined := map[string]string{}
	used := map[string]bool{}
	for _, tags := range tagSets {
		cfg := &packages.Config{
			Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule,
			Dir:   root,
			Tests: true,
			Env:   append(os.Environ(), "CGO_ENABLED=1"),
		}
		if tags != "" {
			cfg.BuildFlags = []string{"-tags", tags}
		}
		pkgs, err := packages.Load(cfg, "./...")
		if err != nil {
			return nil, err
		}
		if n := packages.PrintErrors(pkgs); n > 0 {
			return nil, fmt.Errorf("tags %q: %d package errors", tags, n)
		}
		for _, p := range pkgs {
			collect(root, p, defined, used)
		}
	}
	var out []finding
	for pos, name := range defined {
		if !used[pos] {
			file, _, _ := strings.Cut(pos, "#")
			key := file + ":" + name
			out = append(out, finding{"unused", key, key + " is never used (any tag set, tests included)"})
		}
	}
	slices.SortFunc(out, func(a, b finding) int { return strings.Compare(a.key, b.key) })
	return out, nil
}

func collect(root string, p *packages.Package, defined map[string]string, used map[string]bool) {
	if p.TypesInfo == nil {
		return
	}
	pos := func(o types.Object) string {
		ps := p.Fset.Position(o.Pos())
		rel, err := filepath.Rel(root, ps.Filename)
		if err != nil || strings.HasPrefix(rel, "..") {
			return ""
		}
		return fmt.Sprintf("%s#%d", rel, ps.Offset)
	}
	for id, o := range p.TypesInfo.Defs {
		if o == nil || id.Name == "_" || o.Pkg() == nil || o.Parent() != o.Pkg().Scope() {
			continue
		}
		if fn, ok := o.(*types.Func); ok && isEntryPoint(fn.Name()) {
			continue
		}
		at := pos(o)
		file, _, _ := strings.Cut(at, "#")
		if at == "" || strings.HasSuffix(file, "_bpf.go") || strings.HasPrefix(file, "graphify-out/") {
			continue
		}
		defined[at] = o.Name()
	}
	mark := func(o types.Object) {
		if o == nil || o.Pkg() == nil {
			return
		}
		used[pos(o)] = true
		switch x := o.(type) {
		case *types.Func:
			used[pos(x.Origin())] = true
		case *types.Var:
			used[pos(x.Origin())] = true
		}
	}
	for _, o := range p.TypesInfo.Uses {
		mark(o)
	}
	for _, sel := range p.TypesInfo.Selections {
		mark(sel.Obj())
	}
}

func isEntryPoint(name string) bool {
	if name == "main" || name == "init" {
		return true
	}
	for _, prefix := range []string{"Test", "Benchmark", "Fuzz", "Example"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// deadFuncs runs deadcode -test once per tag set and reports the funcs every run calls dead. A run
// that fails is an error, never an empty (passing) result.
func deadFuncs(root, version string) ([]finding, error) {
	runs := make([]map[string]bool, len(tagSets))
	errs := make([]error, len(tagSets))
	var wg sync.WaitGroup
	for i, tags := range tagSets {
		wg.Go(func() { runs[i], errs[i] = deadcodeRun(root, version, tags) })
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			return nil, fmt.Errorf("deadcode tags %q: %w", tagSets[i], err)
		}
	}
	var out []finding
	for key := range runs[0] {
		if !slices.ContainsFunc(runs[1:], func(r map[string]bool) bool { return !r[key] }) {
			out = append(out, finding{"deadcode", key, key + " is unreachable (every tag set, tests included)"})
		}
	}
	slices.SortFunc(out, func(a, b finding) int { return strings.Compare(a.key, b.key) })
	return out, nil
}

func deadcodeRun(root, version, tags string) (map[string]bool, error) {
	cmd := exec.Command("go", "run", "golang.org/x/tools/cmd/deadcode@"+version, "-test", "-tags="+tags,
		"-f={{range .Funcs}}{{printf \"%s\\t%s\\n\" .Position.File .Name}}{{end}}", "./...")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	dead := map[string]bool{}
	for line := range strings.Lines(string(stdout)) {
		file, name, ok := strings.Cut(strings.TrimSpace(line), "\t")
		if !ok {
			return nil, fmt.Errorf("unexpected output %q", line)
		}
		if rel, err := filepath.Rel(root, file); err == nil {
			file = rel
		}
		dead[file+":"+name] = true
	}
	return dead, nil
}
