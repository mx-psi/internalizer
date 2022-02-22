package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModulePath(t *testing.T) {
	tests := []struct {
		path    string
		modName string
		err     string
	}{
		{
			path:    "./..",
			modName: "github.com/mx-psi/internalizer",
		},
		{
			path:    "./../internal/testdata/module",
			modName: "example.com",
		},
		{
			path: "somethingthatdoesntexist",
			err:  "failed to read go.mod at \"somethingthatdoesntexist/go.mod\": open somethingthatdoesntexist/go.mod: no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			modName, err := getModuleName(test.path)
			if err != nil {
				assert.EqualError(t, err, test.err)
			} else {
				assert.Equal(t, test.modName, modName)
			}
		})
	}
}

func TestGetModuleStructure(t *testing.T) {
	pkg, pkgMap, err := fromFolderWithUniverse("/Users/pablo.baeyens/Source/internalizer/internal/testdata/module")
	require.NoError(t, err)

	var paths []string
	for importPath := range pkgMap {
		paths = append(paths, importPath)
	}
	assert.ElementsMatch(t, paths, []string{
		"example.com",
		"example.com/a",
		"example.com/b",
		"example.com/b/c",
		"example.com/b/c/d",
		"example.com/b/c/f",
		"example.com/b/e",
		"example.com/g",
	})

	// Children
	assert.Contains(t, pkg.Children, "example.com/a")
	assert.Empty(t, pkgMap["example.com/a"].Children)
	assert.Contains(t, pkgMap["example.com/b"].Children, "example.com/b/c")
	assert.Contains(t, pkgMap["example.com/b/c"].Children, "example.com/b/c/d")
	assert.Empty(t, pkgMap["example.com/b/c/d"].Children)
	assert.Contains(t, pkgMap["example.com/b/c"].Children, "example.com/b/c/f")
	assert.Empty(t, pkgMap["example.com/b/c/f"].Children)
	assert.Contains(t, pkgMap["example.com/b"].Children, "example.com/b/c")

	// Imports
	assert.Empty(t, pkgMap["example.com/a"].Imports)
	assert.Contains(t, pkgMap["example.com/b"].Imports, "example.com/a")
	assert.Contains(t, pkgMap["example.com/b/c"].Imports, "example.com/b")
	assert.Contains(t, pkgMap["example.com/b/c"].Imports, "example.com/b/e")
}
