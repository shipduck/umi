package main

import (
	"time"
	"regexp"
	"strconv"
	"strings"
	"sort"
	"log"
)

type Article struct {
	Date  time.Time
	Slug  string
	Tags  []Tag
	Title string

	Origin string
	Reference string

	Media string
	MediaType string

	// twitter video card
	VideoMp4 string
	VideoWebM string
	VideoOgv string
	VideoJpg string
	VideoWidth int
	VideoHeight int

	// where is article from?
	Filepath string
}

func (a *Article) Url() string {
	return postsDir + a.Slug + "/"
}
func (a *Article) MediaUrl() string {
	return postsDir + a.Slug + "/" + a.Media
}

func (a *Article) PlayerCardUrl() string {
	return postsDir + a.Slug + "/player_card.html"
}

func (a *Article) VideoJpgUrl() string {
	return postsDir + a.Slug + "/" + a.VideoJpg
}
func (a *Article) VideoMp4Url() string {
	return postsDir + a.Slug + "/" + a.VideoMp4
}

func (a *Article) VideoWebmUrl() string {
	return postsDir + a.Slug + "/" + a.VideoWebM
}
func (a *Article) VideoOgvUrl() string {
	return postsDir + a.Slug + "/" + a.VideoOgv
}


func (a *Article) Keywords() string {
	tags := make([]string, len(a.Tags))
	for i, tag := range a.Tags {
		tags[i] = tag.Value
	}
	return strings.Join(tags, ",")
}

func (a *Article) HasTag(tag *Tag) bool {
	for _, t := range a.Tags {
		if t.Value == tag.Value {
			return true
		}
	}
	return false
}


func SortArticles(articles []Article) {
	byTitle := func (a1, a2 *Article) bool {
		return a1.Title < a2.Title
	}
	ps := &articleSorter{articles, byTitle}
	sort.Sort(ps)
}

type articleSorter struct {
	articles []Article
	by func(a1, a2 *Article) bool
}
func (s *articleSorter) Len() int {
	return len(s.articles)
}
func (s *articleSorter) Swap(i, j int) {
	s.articles[i], s.articles[j] = s.articles[j], s.articles[i]
}
func (s *articleSorter) Less(i, j int) bool {
	return s.by(&s.articles[i], &s.articles[j])
}

var dateFormatYMD = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
var dateFormatYMDHMS = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})$`)

func ParseDate(line string) time.Time {
	year, month, day, hour, min, sec := 0, 0, 0, 0, 0, 0
	matchYmdHms := dateFormatYMDHMS.FindStringSubmatch(line)
	if len(matchYmdHms) != 0 {
		year, _ = strconv.Atoi(matchYmdHms[1])
		month, _ = strconv.Atoi(matchYmdHms[2])
		day, _ = strconv.Atoi(matchYmdHms[3])
		hour, _ = strconv.Atoi(matchYmdHms[4])
		min, _ = strconv.Atoi(matchYmdHms[5])
		sec, _ = strconv.Atoi(matchYmdHms[6])
	}

	matchYmd := dateFormatYMD.FindStringSubmatch(line)
	if len(matchYmd) != 0 {
		year, _ = strconv.Atoi(matchYmd[1])
		month, _ = strconv.Atoi(matchYmd[2])
		day, _ = strconv.Atoi(matchYmd[3])
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
}

// markdown 기반으로 작성된것 파싱
func ParseArticleMarkdown(text string) Article {
	article := new(Article)
	article.MediaType = "image"

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}

		k, v := SplitKeyValue(line)
		k = strings.Trim(k, " ")
		v = strings.Trim(v, " ")
		switch k {
		case "date":
			article.Date = ParseDate(v)
		case "tags":
			tags := strings.Split(v, ",")
			article.Tags = make([]Tag, len(tags))
			for i, tag := range tags {
				article.Tags[i] = Tag{strings.Trim(tag, " ")}
			}
		case "slug":
			article.Slug = v
		case "title":
			article.Title = v
		case "media":
			article.Media = v
		case "image_file":
			article.Media = v
		case "origin":
			article.Origin = v
		case "ref":
			article.Reference = v
		case "media_type":
			article.MediaType = v
		case "video_mp4":
			article.VideoMp4 = v
		case "video_webm":
			article.VideoWebM = v
		case "video_ogv":
			article.VideoOgv = v
		case "video_jpg":
			article.VideoJpg = v
		case "video_width":
			article.VideoWidth, _ = strconv.Atoi(v)
		case "video_height":
			article.VideoHeight, _ = strconv.Atoi(v)
		default:
			log.Fatalf("Unknown key, value: %s, %s", k, v)
		}
	}
	return *article
}

var keyValueRe = regexp.MustCompile(`^([^:]+):(.*)$`)

func SplitKeyValue(line string) (key, value string) {
	m := keyValueRe.FindStringSubmatch(line)
	key = strings.Trim(m[1], " ")
	value = strings.Trim(m[2], " ")
	return
}
