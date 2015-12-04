package main

import (
	"github.com/deckarep/golang-set"
	"os"
)

type Context struct {
	// global env
	SITENAME string
	SITEURL  string

	// shared
	Articles []Article
	Tags     []Tag

	// local
	Article *Article
	Tag     *Tag

	Title string
}

func (ctx *Context) Reset(articles []Article) {
	if len(os.Getenv("PUBLISH")) > 0 {
		ctx.SITEURL = "http://zzal.collapsed.me"
	} else {
		ctx.SITEURL = ""
	}

	ctx.SITENAME = "Project UMI"

	ctx.Articles = articles
	ctx.Title = ctx.SITENAME

	// 태그 목록 계산
	tagSet := mapset.NewSet()
	for _, article := range articles {
		for _, tag := range article.Tags {
			tagSet.Add(tag)
		}
	}

	ctx.Tags = make([]Tag, tagSet.Cardinality())
	for i, e := range tagSet.ToSlice() {
		if tag, ok := e.(Tag); ok {
			ctx.Tags[i] = tag
		}
	}
	SortTags(ctx.Tags)
}

func (ctx *Context) Clone() *Context {
	clone := new(Context)
	*clone = *ctx
	return clone
}

func (ctx *Context) ArticlesByTag(tag *Tag) []Article {
	found := 0
	articles := make([]Article, len(ctx.Articles))
	for _, article := range ctx.Articles {
		hasTag := article.HasTag(tag)
		if hasTag {
			articles[found] = article
			found += 1
		}
	}

	filtered := articles[:found]
	SortArticles(filtered)
	return filtered
}
