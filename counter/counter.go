package counter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	BungieURL         = "https://www.bungie.net/Platform/Destiny/Stats/ActivityHistory/%s/%s/%s/?mode=AllPvP&page=%d"
	UserAgent         = "cgc/0.1 (https://github.com/daneharrigan/cgc)"
	TimeFormat        = "2006-01-02"
	TrialsOfOsiris    = 14
	IronBanner        = 19
	Racing            = 29
	PrivateMatchesAll = 32
)

var (
	ErrBungieResponse = errors.New("could not parse bungie.net API")
)

type Results struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Total   int       `json:"total"`
	Periods []*Period `json:"periods"`
}

type Period struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ActivityHistory struct {
	Response struct {
		Data struct {
			Activities []struct {
				Period          string `json:"period"`
				ActivityDetails struct {
					Mode int `json:"mode"`
				} `json:"activityDetails"`
			} `json:"activities"`
		} `json:"data"`
	} `json:"Response"`
	ErrorStatus string `json:"ErrorStatus"`
	Message     string `json:"Message"`
}

type Counter struct {
	apiKey         string
	membershipType string
	membershipId   string
	characterId    string
}

func New(apiKey, membershipType, membershipId, characterId string) *Counter {
	return &Counter{
		apiKey:         apiKey,
		membershipType: membershipType,
		membershipId:   membershipId,
		characterId:    characterId,
	}
}

func (c *Counter) GetResults(from, to time.Time) (*Results, error) {
	results := &Results{
		From:  from.Format(TimeFormat),
		To:    to.Format(TimeFormat),
		Total: 0,
	}

	page := 0
	var currentKey string
	var per *Period

	for {
		page++
		url := fmt.Sprintf(BungieURL, c.membershipType, c.membershipId, c.characterId, page)
		response, err := c.get(url)
		if err != nil {
			return nil, err
		}

		activityHistory := &ActivityHistory{}
		if err := json.NewDecoder(response.Body).Decode(activityHistory); err != nil {
			return nil, ErrBungieResponse
		}

		if activityHistory.Message != "Ok" {
			return nil, ErrBungieResponse
		}

		for _, activity := range activityHistory.Response.Data.Activities {
			switch activity.ActivityDetails.Mode {
			case TrialsOfOsiris:
				fallthrough
			case IronBanner:
				fallthrough
			case PrivateMatchesAll:
				continue
			}

			period, err := time.Parse(time.RFC3339, activity.Period)
			if err != nil {
				return nil, ErrBungieResponse
			}

			if period.Before(from) {
				if per != nil {
					results.Periods = append(results.Periods, per)
				}
				return results, nil
			}

			key := period.Format(TimeFormat)
			if key != currentKey {
				if per != nil {
					results.Periods = append(results.Periods, per)
				}
				currentKey = key
				per = &Period{Date: key}
			}

			per.Count++
			results.Total++
		}
	}

	return results, nil
}

func (c *Counter) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	return http.DefaultClient.Do(req)
}
