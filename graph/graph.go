package graph

// Package is a Go package.
type Package struct {
	// Fullpath of the package, fully qualified.
	Fullpath string
	// IsModule states whether the package is the top-level package of a module.
	IsModule bool
	// FileList is the list of .go files in the package
	FileList []string
	// Children of the package in the package tree
	Children map[string]*Package
	// Imports (first-party) of the package
	Imports map[string]*Package
}

type Universe map[string]*Package

type Graph struct {
	Root *Package
	Univ Universe
}

type WalkGraphFunc func(pkg *Package) error

func (g *Graph) Walk(fn WalkGraphFunc) error {
	err := fn(g.Root)
	if err != nil {
		return err
	}

	for _, child := range g.Root.Children {
		gc := &Graph{Root: child, Univ: g.Univ}
		err := gc.Walk(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func newPackage(fullpath string) *Package {
	return &Package{
		Fullpath: fullpath,
		FileList: []string{},
		Children: map[string]*Package{},
		Imports:  map[string]*Package{},
	}
}
