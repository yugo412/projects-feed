package project

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	DefaultTimezone = "Asia/Jakarta"
)

type Projects struct {
	BaseURL string
}

func (p Projects) Name() string {
	return "Projects.co.id"
}

func (p Projects) GetProjects(page int, tag string) (r []Project, err error) {
	target, err := url.Parse(p.BaseURL)
	if err != nil {
		return r, fmt.Errorf("invalid base url: %s", p.BaseURL)
	}

	target = target.JoinPath("public", "browse_projects", "listing")
	q := target.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("search", tag)
	q.Set("ajax", "1")
	target.RawQuery = q.Encode()

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
			author.Username = author.Name
			author.URL, _ = pub.Attr("href")
		})

		// project's info
		e.DOM.Find("div.col-md-10").Each(func(_ int, pr *goquery.Selection) {
			project := pr.Find("h2").First()
			url, _ := pr.Find("h2 > a").Attr("href")
			desc := strings.TrimSpace(pr.Find("p").Text())

			// parse categories
			var tags []string
			pr.Find("p > span.tag > a").Each(func(_ int, t *goquery.Selection) {
				tags = append(tags, strings.TrimSpace(t.Text()))
			})

			var pubTime time.Time
			var budget Budget
			pr.Find(".col-md-6.align-left").Contents().EachWithBreak(func(i int, s *goquery.Selection) bool {
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

				pr.Find(".col-md-6.align-left").Each(func(i int, s *goquery.Selection) {
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

			r = append(r, Project{
				Vendor:      p.Name(),
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
		log.Printf(LogRequestFormat, p.Name(), req.URL.String())
	})

	err = c.Visit(target.String())

	// some responses are not order correctly
	sort.Slice(r, func(i, j int) bool {
		return r[j].PublishedAt.Before(r[i].PublishedAt)
	})

	return
}
