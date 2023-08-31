package projects

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	memcache "github.com/patrickmn/go-cache"
)

const (
	DefaultTimezone = "Asia/Jakarta"
)

type Author struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Budget struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Project struct {
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	Description   string    `json:"description"`
	PublishedDate time.Time `json:"published_date"`
	Tags          []string  `json:"tags"`
	Author        Author    `json:"author"`
	Budget        Budget    `json:"-"`
}

var (
	cache *memcache.Cache
)

func init() {
	cache = memcache.New(5*time.Minute, 10*time.Hour)
}

func GetProjects(page uint, tag string) (projects []Project, err error) {
	fullURL := fmt.Sprintf("https://projects.co.id/public/browse_projects/listing?search=%s&page=%d&ajax=1", tag, page)

	cacheKey := fmt.Sprintf("page.%x", md5.Sum([]byte(fullURL)))
	if cached, ok := cache.Get(cacheKey); ok {
		return cached.([]Project), nil
	}

	spaces := regexp.MustCompile(`\s{2,}`)

	c := colly.NewCollector()
	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		var author Author

		// author's info
		e.DOM.Find("div.col-md-2").Each(func(_ int, a *goquery.Selection) {
			pub := a.Find("a.short-username").First()

			author.Name = strings.TrimSpace(pub.Text())
			author.URL, _ = pub.Attr("href")
		})

		// project's info
		e.DOM.Find("div.col-md-10").Each(func(_ int, p *goquery.Selection) {
			project := p.Find("h2").First()
			url, _ := p.Find("h2 > a").Attr("href")
			desc := strings.TrimSpace(p.Find("p").Text())

			// parse categories
			var tags []string
			p.Find("p > span.tag > a").Each(func(_ int, t *goquery.Selection) {
				tags = append(tags, strings.TrimSpace(t.Text()))
			})

			var publishedDate time.Time
			p.Find(".col-md-6.align-left").Contents().EachWithBreak(func(i int, s *goquery.Selection) bool {
				if s.Is("strong") && strings.TrimSpace(s.Text()) == "Published Date:" {
					nextNode := s.Get(0).NextSibling
					if nextNode != nil && nextNode.Data != "strong" {
						re := regexp.MustCompile(`(\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}) \w+`)
						match := re.FindStringSubmatch(nextNode.Data)
						if len(match) > 1 {
							loc, _ := time.LoadLocation(DefaultTimezone)
							publishedDate, _ = time.ParseInLocation("02/01/2006 15:04:05", match[1], loc)
						}
					}
					return false
				}
				return true
			})

			projects = append(projects, Project{
				Author:        author,
				Title:         strings.TrimSpace(project.Text()),
				URL:           url,
				Description:   spaces.ReplaceAllString(desc, ""),
				PublishedDate: publishedDate,
				Tags:          tags,
			})
		})
	})

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Requesting to:", req.URL.String())
	})

	err = c.Visit(fullURL)

	defer cache.Set(cacheKey, projects, 5*time.Minute)

	return
}
