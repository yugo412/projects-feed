package projects

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	memcache "github.com/patrickmn/go-cache"
	"github.com/xeonx/timeago"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	DefaultTimezone    = "Asia/Jakarta"
	DefaultCacheExpire = 15 // in minutes
)

type Author struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatar_url"`
}

type Budget struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Project struct {
	Vendor      string    `json:"vendor"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Tags        []string  `json:"tags"`
	Author      Author    `json:"author"`
	Budget      Budget    `json:"budget"`
}

type cacheResponse struct {
	StartedAt time.Time `json:"started_at"`
	FlushedAt time.Time `json:"flushed_at"`
}

type response struct {
	Data       []Project     `json:"data"`
	ServerTime time.Time     `json:"server_time"`
	Cache      cacheResponse `json:"cache"`
}

var (
	cache *memcache.Cache
)

func init() {
	cache = memcache.New(
		DefaultCacheExpire*time.Minute,
		(DefaultCacheExpire+10)*time.Hour,
	)
}

func (p Project) LimitedTags(limit int) []string {
	if len(p.Tags) > limit {
		return p.Tags[:limit]
	}

	return p.Tags
}

func (p Project) RemainTag() int {
	if len(p.Tags) > 3 {
		return len(p.Tags) - 3
	}

	return 0
}

func (p Project) Timeago() string {
	return timeago.English.Format(p.PublishedAt)
}

func (p Project) GetBudget() string {
	printer := message.NewPrinter(language.Indonesian)

	return printer.Sprintf("Rp%.0f - Rp%0.f", p.Budget.Min, p.Budget.Max)
}

func GetProjects(page uint, tag string) (r response, err error) {
	now := time.Now().Local()

	fullURL := fmt.Sprintf("https://projects.co.id/public/browse_projects/listing?search=%s&page=%d&ajax=1", tag, page)

	cacheKey := fmt.Sprintf("page.%x", md5.Sum([]byte(fullURL)))
	if cached, ok := cache.Get(cacheKey); ok {
		r = cached.(response)
		r.ServerTime = now

		return
	}

	spaces := regexp.MustCompile(`\s{2,}`)

	c := colly.NewCollector()
	c.OnHTML("div.row", func(e *colly.HTMLElement) {
		var author Author

		// author's info
		e.DOM.Find("div.col-md-2").Each(func(_ int, a *goquery.Selection) {
			pub := a.Find("a.short-username").First()
			avatar, exists := a.Find("a > img.img-thumbnail").First().Attr("src")
			if exists {
				author.AvatarURL = avatar
			}

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

			var pubTime time.Time
			var budget Budget
			p.Find(".col-md-6.align-left").Contents().EachWithBreak(func(i int, s *goquery.Selection) bool {
				// parse published date
				// sometimes, date is not ordered correctly from origin source
				if s.Is("strong") && strings.TrimSpace(s.Text()) == "Published Date:" {
					nextNode := s.Get(0).NextSibling
					if nextNode != nil && nextNode.Data != "strong" {
						re := regexp.MustCompile(`(\d{2}/\d{2}/\d{4} \d{2}:\d{2}:\d{2}) \w+`)
						match := re.FindStringSubmatch(nextNode.Data)
						if len(match) > 1 {
							loc, _ := time.LoadLocation(DefaultTimezone)
							pubTime, _ = time.ParseInLocation("02/01/2006 15:04:05", match[1], loc)
						}
					}
					return false
				}

				p.Find(".col-md-6.align-left").Each(func(i int, s *goquery.Selection) {
					if strings.Contains(s.Text(), "Published Budget:") {
						text := s.Contents().Not("strong").Text()
						re := regexp.MustCompile(`([\d,]+ - [\d,]+)`)
						match := re.FindStringSubmatch(text)
						if len(match) >= 2 {
							// convert budget range to separate amount
							// as float64
							ranges := strings.Split(match[1], "-")

							san := func(n string) string {
								n = strings.ReplaceAll(n, ",", "")
								return strings.TrimSpace(n)
							}
							budget.Min, _ = strconv.ParseFloat(san(ranges[0]), 64)
							budget.Max, _ = strconv.ParseFloat(san(ranges[1]), 64)
						}
					}
				})

				return true
			})

			r.Data = append(r.Data, Project{
				Vendor:      "Projects.co.id",
				Author:      author,
				Title:       strings.TrimSpace(project.Text()),
				URL:         url,
				Description: spaces.ReplaceAllString(desc, " "),
				PublishedAt: pubTime,
				Budget:      budget,
				Tags:        tags,
			})
		})
	})

	c.OnRequest(func(req *colly.Request) {
		fmt.Println("Requesting to:", req.URL.String())
	})

	err = c.Visit(fullURL)

	// some responses are not order correctly
	sort.Slice(r.Data, func(i, j int) bool {
		return r.Data[j].PublishedAt.Before(r.Data[i].PublishedAt)
	})

	r.ServerTime = now
	r.Cache = cacheResponse{
		StartedAt: now,
		FlushedAt: now.Add(DefaultCacheExpire * time.Minute),
	}
	cache.SetDefault(cacheKey, r)

	return
}
