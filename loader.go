package main

import (
	"io/ioutil"
	"reflect"
	"strings"
	"html/template"
	"container/list"
	"path/filepath"
)

func RemoveUTF8BOM(data []byte) string {
	// bom for utf-8
	bom := []byte{0xEF, 0xBB, 0xBF}
	isBOM := len(data) >= 3 && reflect.DeepEqual(data[0:3], bom)
	if isBOM {
		return string(data[3:])
	} else {
		return string(data)
	}
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

type TemplateLoader struct {
	rawLayout string
	rawIndex string
	rawSearch string
	rawArticle string
	rawArticles string
	rawTag string
	rawTags string
	rawPlayerCard string

	partialTwitterCard string

	index *template.Template
	search *template.Template
	article *template.Template
	articles *template.Template
	tag *template.Template
	tags *template.Template
	playerCard *template.Template
}

func (loader *TemplateLoader) GetTemplate(content string) string {
	text := loader.rawLayout
	text = strings.Replace(text, "{{{content}}}", content, -1)
	text = strings.Replace(text, "{{{partial:twitter_card}}}", loader.partialTwitterCard, -1)
	return text
}

func (loader *TemplateLoader) loadFile(filename string) string {
	const templateDir = "theme/templates/"
	data, err := ioutil.ReadFile(templateDir + filename)
	checkErr(err)
	return RemoveUTF8BOM(data)
}

func (loader *TemplateLoader) LoadAll() {
	loader.rawLayout = loader.loadFile("layout.html")
	loader.rawIndex = loader.loadFile("index.html")
	loader.rawSearch = loader.loadFile("search.html")
	loader.rawArticle = loader.loadFile("article.html")
	loader.rawArticles = loader.loadFile("articles.html")
	loader.rawTag = loader.loadFile("tag.html")
	loader.rawPlayerCard = loader.loadFile("player_card.html")
	loader.partialTwitterCard = loader.loadFile("partial_twitter_card.html")

	createTemplate := func (name string, tpl string) *template.Template {
		funcMap := template.FuncMap{
			"title": strings.Title,
		}
		return template.Must(template.New(name).Funcs(funcMap).Parse(tpl))
	}

	// tag 목록 보여주는 페이지 따로 필요없겠다. 인덱스 그냥 사용
	//loader.rawTags = loader.loadFile("tags.html")
	loader.rawTags = loader.loadFile("index.html")

	indexTpl := loader.GetTemplate(loader.rawIndex)
	loader.index = createTemplate("index", indexTpl)

	searchTpl := loader.GetTemplate(loader.rawSearch)
	loader.search = createTemplate("search", searchTpl)

	articleTpl := loader.GetTemplate(loader.rawArticle)
	loader.article = createTemplate("article", articleTpl)

	articlesTpl := loader.GetTemplate(loader.rawArticles)
	loader.articles = createTemplate("articles", articlesTpl)

	tagTpl := loader.GetTemplate(loader.rawTag)
	loader.tag = createTemplate("tag", tagTpl)

	tagsTpl := loader.GetTemplate(loader.rawTags)
	loader.tags = createTemplate("tags", tagsTpl)

	playerCardTpl := loader.rawPlayerCard
	playerCardTpl = strings.Replace(playerCardTpl, "{{{partial:twitter_card}}}", loader.partialTwitterCard, -1)
	loader.playerCard = createTemplate("playerCard", playerCardTpl)
}

func LoadArticleMarkdown(filepath string) Article {
	data, err := ioutil.ReadFile(filepath)
	checkErr(err)
	article := ParseArticleMarkdown(RemoveUTF8BOM(data))
	article.Filepath = filepath
	return article
}

func FindAllMarkdownFile() *list.List {
	// find all articles filename
	filelist := list.New()
	FindAllMarkdownFile_r("content/article/*", filelist)
	return filelist
}

func FindAllMarkdownFile_r(dirname string, filelist *list.List) {
	files, _ := filepath.Glob(dirname)
	for _, file := range files {
		if filepath.Ext(file) == ".md" {
			filelist.PushBack(file)
		}
		FindAllMarkdownFile_r(file + "/*", filelist)
	}
}

func LoadArticles(filelist *list.List, ch chan Article) {
	for e := filelist.Front(); e != nil; e = e.Next() {
		if file, ok := e.Value.(string); ok {
			//log.Printf("Start: %s\n", file)
			article := LoadArticleMarkdown(file)
			ch <- article
			//log.Printf("Finish: %s\n", file)
		}
	}
	close(ch)
}
