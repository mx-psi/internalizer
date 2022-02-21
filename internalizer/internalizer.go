package internalizer

import (
	"path/filepath"
	"strings"

	"github.com/mx-psi/internalizer/graph"
)

type set map[string]struct{}

type walker struct {
	// importedBy
	importedBy map[string]set
}

func (w *walker) Walk(pkg *graph.Package) error {
	for _, importPkg := range pkg.Imports {
		if _, ok := w.importedBy[importPkg.Fullpath]; !ok {
			w.importedBy[importPkg.Fullpath] = set{}
		}
		w.importedBy[importPkg.Fullpath][pkg.Fullpath] = struct{}{}
	}
	return nil
}

// lt says if a < b lexicographically.
func lt(fst []string, snd []string) bool {
	for i, a := range fst {
		if i > len(snd)-1 {
			// snd is a prefix of fst
			return false
		}
		b := snd[i]
		if a < b {
			return true
		} else if a > b {
			return false
		}
	}

	// if they are equal, then false,
	// otherwise, fst is a prefix of snd
	return len(fst) != len(snd)
}

// lcp gets the longest common prefix.
// Adapted From https://rosettacode.org/wiki/Longest_common_prefix#Go
func lcp(l [][]string) []string {
	// Special cases first
	switch len(l) {
	case 0:
		return nil
	case 1:
		return l[0]
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := l[0], l[0]
	for _, s := range l[1:] {
		switch {
		case lt(s, min):
			min = s
		case lt(max, s):
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}

func Internalize(g *graph.Graph) (map[string]string, error) {
	w := &walker{importedBy: map[string]set{}}
	err := g.Walk(w.Walk)
	if err != nil {
		return nil, err
	}

	moves := map[string]string{}

	for path, importedBySet := range w.importedBy {
		imports := make([][]string, 0, len(importedBySet)+1)
		imports = append(imports, strings.Split(path, "/"))
		for importPath := range importedBySet {
			imports = append(imports, strings.Split(importPath, "/"))
		}
		prefix := strings.Join(lcp(imports), "/")
		if prefix != path {
			moves[path] = filepath.Join(prefix, "internal", path[len(prefix):])
		}
	}

	return moves, nil
}
