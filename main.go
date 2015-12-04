package main

import (
	//"log"
)

func main() {
	articleChannel := make(chan Article)
	filelist := FindAllMarkdownFile()
	go LoadArticles(filelist, articleChannel)

	// convert articles to array
	articles := make([]Article, filelist.Len())
	for i := 0; i < len(articles); i++ {
		articles[i] = <-articleChannel
	}

	ctx := &Context{}
	ctx.Reset(articles)

	// load template
	templateLoader := TemplateLoader{}
	templateLoader.LoadAll()

	gen := &Generator{&templateLoader}

	wr := Writer{}
	wr.PrepareDirectory()
	wr.Write(ctx, gen)
}
