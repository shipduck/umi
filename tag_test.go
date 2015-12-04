package main

import (
	"testing"
)

func TestSortTags(t *testing.T) {
	cases := []struct {
		tags []Tag
		expected []Tag
	}{
		{ []Tag{}, []Tag{} },
		{ []Tag{{"a"},{"b"}}, []Tag{{"a"},{"b"}} },
		{ []Tag{{"b"},{"a"}}, []Tag{{"a"},{"b"}} },
	}
	for _, c := range cases {
		SortTags(c.tags)
		for i := 0; i < len(c.tags); i++ {
			if c.tags[i].Value != c.expected[i].Value {
				t.Errorf("SortTag  - expected %q, got %q", c.tags, c.expected)
			}
		}
	}
}
