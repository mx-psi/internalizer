package graph

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// getModuleName from its go.mod file.
func getModuleName(goModPath string) (modName string, err error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod at %q: %w", goModPath, err)
	}

	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return "", fmt.Errorf("failed to parse go.mod at %q: %w", goModPath, err)
	}

	return f.Module.Mod.Path, nil
}

type walker struct {
	modName  string
	basePath string
	pkg      *Package
	pkgMap   Universe
}

func (w *walker) getOrCreatePackage(path string) *Package {
	cur := w.pkg
	paths := strings.Split(strings.TrimPrefix(path, w.pkg.Fullpath), "/")
	for _, curFolder := range paths {
		if curFolder == "" {
			continue
		}
		curpath := filepath.Join(cur.Fullpath, curFolder)
		if _, ok := cur.Children[curpath]; !ok {
			cur.Children[curpath] = newPackage(curpath)
			w.pkgMap[curpath] = cur.Children[curpath]
		}
		cur = cur.Children[curpath]
	}
	return cur
}

func (w *walker) walkGoFile(relPath string, d fs.DirEntry) error {
	dir := filepath.Dir(relPath)
	pkg := w.getOrCreatePackage(w.modName + "/" + dir)
	fullPath := filepath.Join(w.basePath, relPath)
	pkg.FileList = append(pkg.FileList, fullPath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fullPath, nil, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", fullPath, err)
	}

	for _, importStmt := range f.Imports {
		importPath := strings.Trim(importStmt.Path.Value, "\"")
		if !strings.HasPrefix(importPath, w.modName) {
			continue // skip imports outside of module
		}
		if _, ok := pkg.Imports[importPath]; !ok {
			pkg.Imports[importPath] = w.getOrCreatePackage(importPath)
		}
	}

	return nil
}

func (w *walker) WalkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() {
		return nil
	}
	if strings.HasSuffix(path, ".go") {
		err = w.walkGoFile(path, d)
	}
	return err
}

// FromFolder returns a package and a universe from a given folder.
func FromFolder(path string) (*Graph, error) {
	modName, err := getModuleName(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("failed to get module name: %w", err)
	}

	pkg := newPackage(modName)
	walker := &walker{
		modName:  modName,
		basePath: path,
		pkg:      pkg,
		pkgMap: Universe{
			pkg.Fullpath: pkg,
		},
	}
	filesys := os.DirFS(path)
	err = fs.WalkDir(filesys, ".", walker.WalkDir)
	if err != nil {
		return nil, fmt.Errorf("failed to walk %q: %w", modName, err)
	}

	return &Graph{walker.pkg, walker.pkgMap}, nil
}
