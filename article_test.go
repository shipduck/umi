package main

import (
	"testing"
	"time"
)

func TestSplitKeyValue(t *testing.T) {
	cases := []struct {
		in    string
		key   string
		value string
	}{
		{"key:value", "key", "value"},
		{"key: value", "key", "value"},
		{"key : value", "key", "value"},
		{"key:value1 value2", "key", "value1 value2"},
		{"key:value1:value2", "key", "value1:value2"},
		{" key  :  value1 value2   ", "key", "value1 value2"},
	}
	for _, c := range cases {
		got_key, got_value := SplitKeyValue(c.in)
		if got_key != c.key {
			t.Errorf("SplitKeyValue key - expected %q, got %q", c.key, got_key)
		}
		if got_value != c.value {
			t.Errorf("SplitKeyValue value - expected %q, got %q", c.value, got_value)
		}
	}
}

func TestParseDate(t *testing.T) {
	cases := []struct {
		txt  string
		time time.Time
	}{
		{"2010-12-09", time.Date(2010, time.Month(12), 9, 0, 0, 0, 0, time.UTC)},
		{"2010-12-09 01:02:03", time.Date(2010, time.Month(12), 9, 1, 2, 3, 0, time.UTC)},
	}
	for _, c := range cases {
		got := ParseDate(c.txt)
		if got != c.time {
			t.Errorf("ParseDate - expected %q, got %q", got, c.time)
		}
	}
}

func TestParseArticleMarkdown(t *testing.T) {
	text := `
date: 2010-10-09
tags: 김화백, 컴퓨터공학과, 연쇄살인범
slug: artist-kim-why-cs
title: 이 땅의 현실 때문에 컴퓨터 공학과로 진학했지만
media: why-cs.jpg
`
	actual := ParseArticleMarkdown(text)
	expected := Article{
		Date:  ParseDate("2010-10-09"),
		Slug:  "artist-kim-why-cs",
		Tags:  []Tag{{"김화백"}, {"컴퓨터공학과"}, {"연쇄살인범"}},
		Title: "이 땅의 현실 때문에 컴퓨터 공학과로 진학했지만",
		Media: "why-cs.jpg",
	}
	if actual.Date != expected.Date {
		t.Errorf("ParseArticleMarkdown Date - expected %q, got %q", expected.Date, actual.Date)
	}
	if actual.Slug != expected.Slug {
		t.Errorf("ParseArticleMarkdown Slug - expected %q, got %q", expected.Slug, actual.Slug)
	}
	if actual.Tags == nil {
		t.Errorf("ParseArticleMarkdown Tags - nil occur")
	}
	if len(actual.Tags) != len(expected.Tags) {
		t.Errorf("ParseArticleMarkdown Tags - expected %q, got %q", expected.Tags, actual.Tags)
	}
	for i, actualVal := range actual.Tags {
		expectedVal := expected.Tags[i]
		if actualVal != expectedVal {
			t.Errorf("ParseArticleMarkdown Tags - expected %q, got %q", expected.Tags, actual.Tags)
		}
	}
	if actual.Title != expected.Title {
		t.Errorf("ParseArticleMarkdown Title - expected %q, got %q", expected.Title, actual.Title)
	}
	if actual.Media != expected.Media {
		t.Errorf("ParseArticleMarkdown Media - expected %q, got %q", expected.Media, actual.Media)
	}
}

func TestArticleHasTag(t *testing.T) {
	article := Article{
		Tags: []Tag{{"김화백"}, {"컴퓨터공학과"}, {"연쇄살인범"}},
	}
	cases := []struct {
		tag   *Tag
		found bool
	}{
		{&Tag{"김화백"}, true},
		{&Tag{"not-exist"}, false},
	}
	for _, c := range cases {
		got := article.HasTag(c.tag)
		if got != c.found {
			t.Errorf("HasTag - expected %t, got %t", c.found, got)
		}
	}
}

func TestContextArticlesByTag(t *testing.T) {
	article1 := Article{Tags: []Tag{{"foo"}, {"bar"}}}
	article2 := Article{Tags: []Tag{{"bar"}, {"spam"}}}
	ctx := Context{
		Articles: []Article{article1, article2},
	}

	cases := []struct {
		tag    *Tag
		retval []Article
	}{
		{&Tag{"foo"}, []Article{article1}},
		{&Tag{"bar"}, []Article{article1, article2}},
		{&Tag{"not-exist"}, nil},
	}
	for _, c := range cases {
		got := ctx.ArticlesByTag(c.tag)
		// 내용비교까진 귀찮으니까 일단 요소 갯수만 비교
		if len(got) != len(c.retval) {
			t.Errorf("ArticlesByTag len - expected %q, got %q", got, c.retval)
		}
	}
}

func TestSortArticles(t *testing.T) {
	article1 := Article{Title: "a"}
	article2 := Article{Title: "b"}
	articles := []Article{article2, article1}

	SortArticles(articles)

	if articles[0].Title != article1.Title {
		t.Errorf("SortArticles - got %q", articles)
	}
	if articles[1].Title != article2.Title {
		t.Errorf("SortArticles -got %q", articles)
	}
}
