package main

import (
	"html/template"
	"bytes"
	"strings"
)

type Generator struct {
	loader *TemplateLoader
}

func (gen *Generator) GenerateHtml(tpl *template.Template, ctx *Context) string {
	var b bytes.Buffer
	err := tpl.Execute(&b, ctx)
	if err != nil { panic(err) }

	lines := strings.Split(b.String(), "\n")
	buf := ""
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		buf += strings.Trim(line, " ") + "\n"
	}
	return buf
}

func (gen *Generator) GenerateIndexHtml(ctx *Context) string {
	return gen.GenerateHtml(gen.loader.index, ctx)
}

func (gen *Generator) GenerateSearchHtml(ctx *Context) string {
	return gen.GenerateHtml(gen.loader.search, ctx)
}

func (gen *Generator) GenerateArticleHtml(ctx *Context, article *Article) string {
	ctx.Title = article.Title
	return gen.GenerateHtml(gen.loader.article, ctx)
}

func (gen *Generator) GeneratePlayerCardHtml(ctx *Context, article *Article) string {
	ctx.Title = article.Title
	return gen.GenerateHtml(gen.loader.playerCard, ctx)
}

func (gen *Generator) GenerateArticlesHtml(ctx *Context) string {
	return gen.GenerateHtml(gen.loader.articles, ctx)
}

func (gen *Generator) GenerateArticlesJson(ctx *Context) string {
	return ""
}

func (gen *Generator) GenerateTagHtml(ctx *Context, tag *Tag) string {
	ctx.Title = tag.Value
	return gen.GenerateHtml(gen.loader.tag, ctx)
}

func (gen *Generator) GenerateTagsHtml(ctx *Context) string {
	return gen.GenerateHtml(gen.loader.tags, ctx)
}
