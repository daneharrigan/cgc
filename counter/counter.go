package counter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	BungieURL  = "https://www.bungie.net/Platform/Destiny/Stats/ActivityHistory/%s/%s/%s/?mode=%s&page=%d"
	UserAgent  = "cgc/0.1 (https://github.com/daneharrigan/cgc)"
	TimeFormat = "2006-01-02"
)

var Modes = strings.Join([]string{
	"Lockdown",
	"ThreeVsThree",
	"Control",
	"FreeForAll",
	"Doubles",
	"Elimination",
	"Rift",
	"AllMayhem",
	"ZoneControl",
	"Supremacy",
}, ",")

var (
	ErrBungieResponse = errors.New("could not parse bungie.net API")
)

type Results struct {
	From      string         `json:"from"`
	To        string         `json:"to"`
	Total     int            `json:"total"`
	Breakdown map[string]int `json:"breakdown"`
}

type ActivityHistory struct {
	Response struct {
		Data struct {
			Activities []struct {
				Period string `json:"period"`
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
		From:      from.Format(TimeFormat),
		To:        to.Format(TimeFormat),
		Breakdown: make(map[string]int),
		Total:     0,
	}

	page := 0
	for {
		page++
		url := fmt.Sprintf(BungieURL, c.membershipType, c.membershipId, c.characterId, Modes, page)
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

		var counted int
		for _, activity := range activityHistory.Response.Data.Activities {
			period, err := time.Parse(time.RFC3339, activity.Period)
			if err != nil {
				return nil, ErrBungieResponse
			}

			key := period.Format(TimeFormat)
			results.Breakdown[key]++
			results.Total++
			counted++
		}

		if counted == 0 {
			break
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
