package main

import (
	"github.com/rainycape/unidecode"
	"sort"
	"strings"
)

type Tag struct {
	Value string
}

func (tag *Tag) Url() string {
	return tagsDir + tag.Slug()
}
func (tag *Tag) Slug() string {
	slug := unidecode.Unidecode(tag.Value)
	slug = strings.Replace(slug, " ", "-", -1)
	return slug
}

func SortTags(tags []Tag) {
	byValue := func(t1, t2 *Tag) bool {
		return t1.Value < t2.Value
	}
	ps := &tagSorter{
		tags: tags,
		by:   byValue,
	}
	sort.Sort(ps)
}

type tagSorter struct {
	tags []Tag
	by   func(t1, t2 *Tag) bool
}

func (s *tagSorter) Len() int {
	return len(s.tags)
}

func (s *tagSorter) Swap(i, j int) {
	s.tags[i], s.tags[j] = s.tags[j], s.tags[i]
}

func (s *tagSorter) Less(i, j int) bool {
	return s.by(&s.tags[i], &s.tags[j])
}
