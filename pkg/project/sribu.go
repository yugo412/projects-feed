package project

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type query struct {
	Language string `json:"language"`
}

type payload struct {
	Query         string `json:"query"`
	Variables     string `json:"variables"`
	OperationName string `json:"operationName"`
}

type response struct {
	Data struct {
		JobsList struct {
			Status  string      `json:"status"`
			Count   int         `json:"count"`
			Message interface{} `json:"message"`
			Error   interface{} `json:"error"`
			Data    []struct {
				Id            string    `json:"id"`
				Title         string    `json:"title"`
				CreatedAt     time.Time `json:"created_at"`
				AppliedNumber int       `json:"applied_number"`
				HiredNumber   int       `json:"hired_number"`
				TdjobsHireds  []struct {
					Member struct {
						Id     string `json:"id"`
						UserId string `json:"user_id"`
					} `json:"mmember"`
				} `json:"tdjobs_hireds"`
				Status           int    `json:"status"`
				FlowStatus       string `json:"flow_status"`
				Amount           int    `json:"amount"`
				TdjobsApplicants []struct {
					Id       string `json:"id"`
					MemberId string `json:"member_id"`
				} `json:"tdjobs_applicants"`
				FreelancerBudget int `json:"freelancer_budget"`
				Member           struct {
					Id     string `json:"id"`
					Name   string `json:"name"`
					UserID string `json:"user_id"`
				} `json:"mmember"`
				IsNoHired        int       `json:"is_no_hired"`
				Budget           float64   `json:"budget"`
				CurrencyCode     string    `json:"currency_code"`
				Deadline         string    `json:"deadline"`
				PostingTimeLimit time.Time `json:"posting_time_limit"`
				SubCategory      struct {
					Id               string `json:"id"`
					Name             string `json:"name"`
					NameInEn         string `json:"name_en"`
					MainAssetUrl     string `json:"main_asset_url"`
					FallbackAssetUrl string `json:"fallback_asset_url"`
					Mcategory        struct {
						Id               string      `json:"id"`
						MainAssetUrl     string      `json:"main_asset_url"`
						FallbackAssetUrl string      `json:"fallback_asset_url"`
						Name             string      `json:"name"`
						NameEn           string      `json:"name_en"`
						Mlandingpage     interface{} `json:"mlandingpage"`
					} `json:"mcategory"`
					Mlandingpage interface{} `json:"mlandingpage"`
				} `json:"msubcategory"`
			} `json:"data"`
		} `json:"jobsList"`
		TypesList struct {
			Status  string      `json:"status"`
			Count   int         `json:"count"`
			Message interface{} `json:"message"`
			Error   interface{} `json:"error"`
			Data    []struct {
				Id     string `json:"id"`
				Name   string `json:"name"`
				NameEn string `json:"name_en"`
			} `json:"data"`
		} `json:"typesList"`
		Subcategories struct {
			Status  string      `json:"status"`
			Message interface{} `json:"message"`
			Error   interface{} `json:"error"`
			Data    []struct {
				Id        string `json:"id"`
				Name      string `json:"name"`
				NameEn    string `json:"name_en"`
				Sort      int    `json:"sort"`
				Mcategory struct {
					Id string `json:"id"`
				} `json:"mcategory"`
			} `json:"data"`
		} `json:"subcategories"`
	} `json:"data"`
}

type Sribu struct {
	BaseURL string
}

func (s Sribu) Name() string {
	return "Sribu.com"
}

func (s Sribu) GetProjects(page int, tag string) (p []Project, err error) {
	// sribu doesn't support tag filtering right now
	if tag != "" {
		return
	}

	target, err := url.Parse(s.BaseURL)
	if err != nil {
		return
	}

	target = target.JoinPath("public")
	method := "POST"

	params := strings.NewReader(`{
		"query": "{\n                jobsList:jobsListPublic(page: 1,job_type:0, per_page:10, status: 1, posting_time_limit_start:\"2023-09-05T14:25:39+07:00\" ,language: \"\"\"id\"\"\"    ,order_by:[{name:\"\"\"created_at\"\"\",order:-1}] )\n                {status count message error data{\n                  id title created_at\n                  applied_number\n                  hired_number\n                  tdjobs_hireds{mmember{id user_id}}\n                  status flow_status amount tdjobs_applicants{id member_id}\n                  freelancer_budget\n                  mmember{id name user_id}\n                  is_no_hired\n                  budget currency_code deadline\n                  posting_time_limit msubcategory{id name name_en main_asset_url fallback_asset_url mcategory{id main_asset_url fallback_asset_url name name_en mlandingpage {id url}} mlandingpage {id url}} }}\n                  typesList(page: 0, per_page:0, type:\"categories\", is_thanos_hide: 0, status: 1, order_by: [{name: \"navbar_sort\", order: 1}])\n                  {status count message error data\n                    {id name name_en}\n                  }\n                  subcategories:subcategoriesList(page: 0, per_page:0, status: 1, order_by: [{name: \"sort\", order: 1}])\n                    {status message error data{id name name_en sort mcategory{id}}}\n                }",
		"variables": null,
		"operationName": null
	}`)

	req, err := http.NewRequest(method, target.String(), params)
	if err != nil {
		return
	}

	baseURL, _ := url.Parse("https://www.sribu.com")

	req.Header.Set("Authority", target.Host)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Origin", baseURL.String())
	req.Header.Add("Referer", baseURL.String())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	log.Printf(LogRequestFormat, s.Name(), req.URL.String())

	defer res.Body.Close()

	var body response
	err = json.NewDecoder(res.Body).Decode(&body)

	for _, j := range body.Data.JobsList.Data {
		p = append(p, Project{
			Vendor:      s.Name(),
			Title:       j.Title,
			URL:         baseURL.JoinPath("id", "jobs", j.Id).String(),
			PublishedAt: j.CreatedAt.Local(),
			Budget: Budget{
				Min: j.Budget,
				Max: j.Budget,
			},
			Author: Author{
				Name:      j.Member.Name,
				Username:  j.Member.UserID,
				URL:       baseURL.JoinPath("id", "users", j.Member.UserID).String(),
				AvatarURL: "/public/img/sribu-small.png",
			},
			Tags: []string{j.SubCategory.NameInEn},
		})
	}

	return
}
