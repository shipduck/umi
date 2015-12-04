package main

import (
	"os"
	"path/filepath"
	//"strings"
	//"fmt"
)

const outputDir = "output/"
const postsDir = "posts/"
const tagsDir = "tags/"
const staticDir = "static/"

type Writer struct {
	loader *TemplateLoader
}

func (wr *Writer) PrepareDirectory() {
	os.MkdirAll(outputDir, 0755)
	os.MkdirAll(outputDir+postsDir, 0755)
	os.MkdirAll(outputDir+tagsDir, 0755)

	// copy static
	dstStaticDir := outputDir + staticDir
	srcStaticDir := "theme/" + staticDir
	os.MkdirAll(outputDir+staticDir, 0755)
	files, _ := filepath.Glob(srcStaticDir + "*")
	for _, file := range files {
		_, filename := filepath.Split(file)
		Copy(file, dstStaticDir+filename)
	}

	// copy extra
	srcExtraDir := "content/extra/"
	files, _ = filepath.Glob(srcExtraDir + "*")
	for _, file := range files {
		_, filename := filepath.Split(file)
		Copy(file, outputDir+filename)
	}
}

func (wr *Writer) Write(ctx *Context, gen *Generator) {
	startCh := make(chan string)
	finishCh := make(chan string)

	go wr.WriteIndex(ctx.Clone(), gen, startCh, finishCh)
	go wr.WriteArticles(ctx.Clone(), gen, startCh, finishCh)
	go wr.WriteTags(ctx.Clone(), gen, startCh, finishCh)
	go wr.WriteSearch(ctx.Clone(), gen, startCh, finishCh)
	go wr.WriteArticlesJson(ctx.Clone(), gen, startCh, finishCh)

	for _, article := range ctx.Articles {
		cloneCtx := ctx.Clone()
		go wr.WriteArticle(cloneCtx, article, gen, startCh, finishCh)
	}

	for _, tag := range ctx.Tags {
		cloneCtx := ctx.Clone()
		go wr.WriteTag(cloneCtx, tag, gen, startCh, finishCh)
	}

	completeMap := make(map[string]bool)
Loop:
	for {
		select {
		case file := <-startCh:
			completeMap[file] = false

		case file := <-finishCh:
			completeMap[file] = true
			//fmt.Printf("write %d %s\n", counter, file)

			allChecked := true
			for _, checked := range completeMap {
				allChecked = allChecked && checked
			}
			if allChecked {
				//fmt.Println("Success")
				break Loop
			}
		}
	}
}

func (wr *Writer) WriteIndex(ctx *Context, gen *Generator, startCh chan string, finishCh chan string) {
	dst := outputDir + "index.html"
	startCh <- dst

	html := gen.GenerateIndexHtml(ctx)
	wr.WriteHtmlFile(html, dst)

	finishCh <- dst
}

func (wr *Writer) WriteSearch(ctx *Context, gen *Generator, startCh chan string, finishCh chan string) {
	dst := outputDir + "search.html"
	startCh <- dst

	html := gen.GenerateSearchHtml(ctx)
	wr.WriteHtmlFile(html, dst)

	finishCh <- dst
}

func (wr *Writer) WriteArticles(ctx *Context, gen *Generator, startCh chan string, finishCh chan string) {
	dst := outputDir + "articles.html"
	startCh <- dst

	html := gen.GenerateArticlesHtml(ctx)
	wr.WriteHtmlFile(html, dst)

	finishCh <- dst
}

func (wr *Writer) WriteArticlesJson(ctx *Context, gen *Generator, startCh chan string, finishCh chan string) {
	dst := outputDir + "articles.json"
	startCh <- dst

	txt := gen.GenerateArticlesJson(ctx)
	wr.WriteHtmlFile(txt, dst)

	finishCh <- dst
}

func (wr *Writer) WriteArticle(ctx *Context, article Article, gen *Generator, startCh chan string, finishCh chan string) {
	dstDir := outputDir + postsDir + article.Slug
	dst := dstDir + "/index.html"
	startCh <- dst

	ctx.Article = &article
	os.MkdirAll(dstDir, 0755)
	html := gen.GenerateArticleHtml(ctx, &article)
	wr.WriteHtmlFile(html, dst)

	srcDir, _ := filepath.Split(article.Filepath)
	srcImage := filepath.Join(srcDir, article.Media)
	dstImage := filepath.Join(dstDir, article.Media)
	Copy(srcImage, dstImage)

	if article.MediaType == "video" {
		playerCardHtml := gen.GeneratePlayerCardHtml(ctx, &article)
		playerCardDst := dstDir + "/player_card.html"
		wr.WriteHtmlFile(playerCardHtml, playerCardDst)

		files := []string{
			article.VideoMp4,
			article.VideoWebM,
			article.VideoOgv,
			article.VideoJpg,
		}
		for _, file := range files {
			srcFile := filepath.Join(srcDir, file)
			dstFile := filepath.Join(dstDir, file)
			Copy(srcFile, dstFile)
		}
	}

	finishCh <- dst
}

func (wr *Writer) WriteTags(ctx *Context, gen *Generator, startCh chan string, finishCh chan string) {
	dst := outputDir + "tags.html"
	startCh <- dst

	html := gen.GenerateTagsHtml(ctx)
	wr.WriteHtmlFile(html, dst)

	finishCh <- dst
}

func (wr *Writer) WriteTag(ctx *Context, tag Tag, gen *Generator, startCh chan string, finishCh chan string) {
	dstDir := outputDir + tagsDir + tag.Slug()
	dst := dstDir + "/index.html"
	startCh <- dst

	ctx.Tag = &tag
	html := gen.GenerateTagHtml(ctx, &tag)
	os.MkdirAll(dstDir, 0755)
	wr.WriteHtmlFile(html, dst)
	finishCh <- dst
}

func (wr *Writer) WriteHtmlFile(html string, filepath string) {
	f, err := os.Create(filepath)
	checkErr(err)
	_, err = f.WriteString(html)
	f.Close()
}
