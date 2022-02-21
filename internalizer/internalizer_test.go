package internalizer

import (
	"testing"

	"github.com/mx-psi/internalizer/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalizer(t *testing.T) {
	graph, err := graph.FromFolder("/Users/pablo.baeyens/Source/internalizer/internal/testdata/module")
	require.NoError(t, err)

	remaps, err := Internalize(graph)
	require.NoError(t, err)

	assert.Equal(t, remaps, map[string]string{
		"example.com/a":   "example.com/internal/a",
		"example.com/b/e": "example.com/b/internal/e",
		"example.com/g":   "example.com/internal/g",
	})
}

func TestLT(t *testing.T) {
	tests := []struct {
		name     string
		fst      []string
		snd      []string
		expected bool
	}{
		{
			name:     "fst is prefix of snd",
			fst:      []string{"example.com", "a"},
			snd:      []string{"example.com", "a", "b"},
			expected: true,
		},
		{
			name:     "snd is prefix of fst",
			fst:      []string{"example.com", "a", "b"},
			snd:      []string{"example.com", "a"},
			expected: false,
		},
		{
			name:     "fst == snd",
			fst:      []string{"example.com", "a"},
			snd:      []string{"example.com", "a"},
			expected: false,
		},
		{
			name:     "fst is prefix of snd (fst is empty)",
			fst:      []string{},
			snd:      []string{"example.com", "a"},
			expected: true,
		},
		{
			name:     "snd is prefix of fst (snd is empty)",
			fst:      []string{"example.com", "a"},
			snd:      []string{},
			expected: false,
		},
		{
			name:     "fst == snd (empty)",
			fst:      []string{},
			snd:      []string{},
			expected: false,
		},
		{
			name:     "common prefix, fst < snd",
			fst:      []string{"example.com", "a"},
			snd:      []string{"example.com", "b"},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, lt(test.fst, test.snd))
		})
	}
}

func TestLCP(t *testing.T) {
	tests := []struct {
		name string
		l    [][]string
		lcp  []string
	}{
		{
			name: "empty",
			l:    nil,
			lcp:  nil,
		},
		{
			name: "only one",
			l:    [][]string{{"example.com", "a"}},
			lcp:  []string{"example.com", "a"},
		},
		{
			name: "example.com",
			l: [][]string{
				{"example.com", "a"},
				{"example.com", "b"},
			},
			lcp: []string{"example.com"},
		},
		{
			name: "example.com",
			l: [][]string{
				{"example.com", "b"},
				{"example.com", "b", "c", "d"},
				{"example.com", "b", "c", "e"},
			},
			lcp: []string{"example.com", "b"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.lcp, lcp(test.l))
		})
	}
}
