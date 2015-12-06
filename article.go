package main

import (
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	Date  time.Time `key:"date"`
	Slug  string    `key:"slug"`
	Tags  []Tag     `key:"tags"`
	Title string    `key:"title"`

	Origin    string `key:"origin"`
	Reference string `key:"ref"`

	Media     string `key:"media"`
	MediaType string `key:"media_type"`

	// twitter video card
	VideoMp4    string `key:"video_mp4"`
	VideoWebM   string `key:"video_webm"`
	VideoOgv    string `key:"video_ogv"`
	VideoJpg    string `key:"video_jpg"`
	VideoWidth  int    `key:"video_width"`
	VideoHeight int    `key:"video_height"`

	// where is article from?
	Filepath string `key:-`
}

var fieldsMap map[string]reflect.StructField = map[string]reflect.StructField{}

func init() {
	t := reflect.TypeOf(Article{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag
		key := tag.Get("key")
		if key != "" {
			fieldsMap[key] = f
		}
	}
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
	byTitle := func(a1, a2 *Article) bool {
		return a1.Title < a2.Title
	}
	ps := &articleSorter{articles, byTitle}
	sort.Sort(ps)
}

type articleSorter struct {
	articles []Article
	by       func(a1, a2 *Article) bool
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

	val := reflect.Indirect(reflect.ValueOf(article))

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		k, v := SplitKeyValue(line)
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)

		if field, exist := fieldsMap[k]; exist {
			val_field := val.FieldByIndex(field.Index)
			switch val_field.Type().Kind() {
			case reflect.Struct:
				val_field.Set(reflect.ValueOf(ParseDate(v)))
			case reflect.Slice:
				l := val_field.Interface().([]Tag)
				tags := strings.Split(v, ",")
				for _, tag := range tags {
					l = append(l, Tag{strings.TrimSpace(tag)})
				}
				val_field.Set(reflect.ValueOf(l))
			case reflect.Int:
				i, _ := strconv.Atoi(v)
				val_field.SetInt(int64(i))
			case reflect.String:
				val_field.SetString(v)
			}
		} else {
			log.Fatalf("Unknown key, value: %s, %s", k, v)
		}
	}

	return *article
}

var keyValueRe = regexp.MustCompile(`^([^:]+):(.*)$`)

func SplitKeyValue(line string) (key, value string) {
	m := keyValueRe.FindStringSubmatch(line)
	key = strings.TrimSpace(m[1])
	value = strings.TrimSpace(m[2])
	return
}
