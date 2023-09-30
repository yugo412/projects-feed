package project

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gookit/slog"
	"github.com/oklog/ulid/v2"
)

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

type detailResponse struct {
	Data struct {
		JobFindOne struct {
			Status  string      `json:"status"`
			Message interface{} `json:"message"`
			Error   interface{} `json:"error"`
			Data    struct {
				Header struct {
					Id                      string `json:"id"`
					FlowStatus              string `json:"flow_status"`
					JobType                 int    `json:"job_type"`
					IsNoHired               int    `json:"is_no_hired"`
					PackageWorkspaceCreated int    `json:"package_workspace_created"`
					Title                   string `json:"title"`
					Mmember                 struct {
						Id            string `json:"id"`
						Name          string `json:"name"`
						UserId        string `json:"user_id"`
						Language      string `json:"language"`
						UrlAvatarFull string `json:"url_avatar_full"`
						Country       string `json:"country"`
						MmemberStatus struct {
							ClAvgRating        float64 `json:"cl_avg_rating"`
							ClWorkspacesClosed int     `json:"cl_workspaces_closed"`
							ClJobsPosted       int     `json:"cl_jobs_posted"`
							ClReviewed         int     `json:"cl_reviewed"`
						} `json:"mmember_status"`
					} `json:"mmember"`
					CreatedAt        time.Time `json:"created_at"`
					AppliedNumber    int       `json:"applied_number"`
					HiredNumber      int       `json:"hired_number"`
					Description      string    `json:"description"`
					CurrencyCode     string    `json:"currency_code"`
					WritingLanguage  string    `json:"writing_language"`
					WritingStyle     string    `json:"writing_style"`
					TargetAudience   string    `json:"target_audience"`
					FreelancerBudget int       `json:"freelancer_budget"`
					Msubcategory     struct {
						Id               string `json:"id"`
						Cname            string `json:"cname"`
						MainAssetUrl     string `json:"main_asset_url"`
						FallbackAssetUrl string `json:"fallback_asset_url"`
						Name             string `json:"name"`
						NameEn           string `json:"name_en"`
						MinimumBudget    int    `json:"minimum_budget"`
						MinimumBudgetUsd int    `json:"minimum_budget_usd"`
						BriefType        int    `json:"brief_type"`
						Mcategory        struct {
							Id               string `json:"id"`
							Name             string `json:"name"`
							NameEn           string `json:"name_en"`
							MainAssetUrl     string `json:"main_asset_url"`
							FallbackAssetUrl string `json:"fallback_asset_url"`
							Mlandingpage     struct {
								Id  string `json:"id"`
								Url string `json:"url"`
							} `json:"mlandingpage"`
						} `json:"mcategory"`
						Mlandingpage struct {
							Id  string `json:"id"`
							Url string `json:"url"`
						} `json:"mlandingpage"`
					} `json:"msubcategory"`
					Amount                  int       `json:"amount"`
					FreelancerPayableAmount int       `json:"freelancer_payable_amount"`
					ComissionAmount         int       `json:"comission_amount"`
					Budget                  int       `json:"budget"`
					Deadline                string    `json:"deadline"`
					PostingTimeLimit        time.Time `json:"posting_time_limit"`
					CreatedBy               string    `json:"created_by"`
					Status                  int       `json:"status"`
					UpdatedAt               time.Time `json:"updated_at"`
					IndustryId              string    `json:"industry_id"`
					CompanyName             string    `json:"company_name"`
					CompanyWebsite          string    `json:"company_website"`
					AboutCompany            string    `json:"about_company"`
					Keywords                string    `json:"keywords"`
					ProjectTitle            string    `json:"project_title"`
					ContentLength           int       `json:"content_length"`
					Revision                int       `json:"revision"`
					Mindustry               struct {
						Id     string `json:"id"`
						Name   string `json:"name"`
						NameEn string `json:"name_en"`
					} `json:"mindustry"`
					Mpackage struct {
						GroupCname string `json:"group_cname"`
						Name       string `json:"name"`
						NameEn     string `json:"name_en"`
						Tier       int    `json:"tier"`
					} `json:"mpackage"`
				} `json:"header"`
				DetailSkills       []interface{} `json:"detail_skills"`
				DetailApplicants   []interface{} `json:"detail_applicants"`
				DetailHired        []interface{} `json:"detail_hired"`
				DetailCustomBriefs []interface{} `json:"detail_custom_briefs"`
				DetailAttachments  []struct {
					Id       string `json:"id"`
					FileName string `json:"file_name"`
					UrlFile  string `json:"url_file"`
				} `json:"detail_attachments"`
			} `json:"data"`
		} `json:"jobFindOne"`
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

	slog.Debugf(LogRequestFormat, s.Name(), req.URL.String())

	defer res.Body.Close()

	var body response
	err = json.NewDecoder(res.Body).Decode(&body)

	for _, j := range body.Data.JobsList.Data {
		p = append(p, Project{
			ID:          ulid.Make().String(),
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

func (s Sribu) GetDetail(URL string) (p Project, err error) {
	target, err := url.Parse(s.BaseURL)
	if err != nil {
		return
	}

	URLs := strings.Split(URL, "/")
	id := URLs[len(URLs)-1]

	target = target.JoinPath("public")
	method := "POST"

	params := strings.NewReader(`{"query": "query jobFindOnePublic($id: UUID!) {\n                  jobFindOne: jobFindOnePublic(id: $id) {\n                    status\n                    message\n                    error\n                    data {\n                      header {\n                        id\n                        flow_status\n                        job_type\n                        is_no_hired\n                        package_workspace_created\n                        title\n                        mmember {\n                          id\n                          name\n                          user_id\n                          language\n                          url_avatar_full\n                          country\n                          mmember_status {\n                            cl_avg_rating\n                            cl_workspaces_closed\n                            cl_jobs_posted\n                            cl_reviewed\n                          }\n                        }\n                        created_at\n                        applied_number\n                        hired_number\n                        description\n                        currency_code\n                        writing_language\n                        writing_style\n                        target_audience\n                        freelancer_budget\n                        msubcategory {\n                          id\n                          cname\n                          main_asset_url\n                          fallback_asset_url\n                          name\n                          name_en\n                          minimum_budget\n                          minimum_budget_usd\n                          brief_type\n                          mcategory {\n                            id\n                            name\n                            name_en\n                            main_asset_url\n                            fallback_asset_url\n                            mlandingpage {\n                              id\n                              url\n                            }\n                          }\n                          mlandingpage {\n                            id\n                            url\n                          }\n                        }\n                        amount\n                        freelancer_payable_amount\n                        comission_amount\n                        budget\n                        deadline\n                        posting_time_limit\n                        created_by\n                        status\n                        updated_at\n                        industry_id\n                        company_name\n                        company_website\n                        about_company\n                        keywords\n                        project_title\n                        content_length\n                        revision\n                        mindustry {\n                          id\n                          name\n                          name_en\n                        }\n                        mpackage {\n                          group_cname\n                          name\n                          name_en\n                          tier\n                        }\n                      }\n                      detail_skills {\n                        id\n                        mskill {\n                          id\n                          name\n                          name_en\n                        }\n                      }\n                      detail_applicants {\n                        id\n                        updated_at\n                        status\n                        mmember {\n                          id\n                          status\n                          user_id\n                          name\n                          url_avatar_full\n                          url_avatar_thumbnail\n                          mmember_status {\n                            fl_rating_avg\n                          }\n                        }\n                      }\n                      detail_hired {\n                        id\n                        updated_at\n                        mmember {\n                          id\n                          status\n                          user_id\n                          name\n                          is_freelancer\n                          url_avatar_full\n                          url_avatar_thumbnail\n                          mmember_status {\n                            fl_rating_avg\n                          }\n                        }\n                      }\n                      detail_custom_briefs {\n                        title\n                        type\n                        content_length\n                        keywords\n                        company_name\n                        about_company\n                        company_website\n                      }\n                      detail_attachments {\n                        id\n                        file_name\n                        url_file\n                      }\n                    }\n                  }\n                  typesList(\n                    page: 1\n                    per_page: 5\n                    status: 1\n                    is_thanos_hide: 0\n                    type: \"categories\"\n                  ) {\n                    status\n                    count\n                    message\n                    error\n                    data {\n                      id\n                      name\n                      name_en\n                    }\n                  }\n                }",
    "variables": {
        "id": "` + id + `"
    },
    "operationName": null}`)

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

	slog.Debugf(LogRequestFormat, s.Name(), req.URL.String())

	defer res.Body.Close()

	var body detailResponse
	err = json.NewDecoder(res.Body).Decode(&body)

	p.Description = body.Data.JobFindOne.Data.Header.Description
	p.Author.AvatarURL = body.Data.JobFindOne.Data.Header.Mmember.UrlAvatarFull

	return
}
