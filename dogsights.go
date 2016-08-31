package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/PagerDuty/godspeed"
	"github.com/codeskyblue/go-sh"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Insights struct {
	Metadata struct {
		BeginTime       string `json:"beginTime"`
		BeginTimeMillis int    `json:"beginTimeMillis"`
		Contents        []struct {
			Function  string `json:"function"`
			OpenEnded bool   `json:"openEnded"`
			Simple    bool   `json:"simple"`
		} `json:"contents"`
		EndTime        string        `json:"endTime"`
		EndTimeMillis  int           `json:"endTimeMillis"`
		EventType      string        `json:"eventType"`
		EventTypes     []string      `json:"eventTypes"`
		GUID           string        `json:"guid"`
		Messages       []interface{} `json:"messages"`
		OpenEnded      bool          `json:"openEnded"`
		RawCompareWith string        `json:"rawCompareWith"`
		RawSince       string        `json:"rawSince"`
		RawUntil       string        `json:"rawUntil"`
		RouterGUID     string        `json:"routerGuid"`
	} `json:"metadata"`
	PerformanceStats struct {
		CacheMisses                int `json:"cacheMisses"`
		CacheSkipped               int `json:"cacheSkipped"`
		DecompressedBytes          int `json:"decompressedBytes"`
		DecompressionTime          int `json:"decompressionTime"`
		FileProcessingTime         int `json:"fileProcessingTime"`
		FileReadCount              int `json:"fileReadCount"`
		FullCacheHits              int `json:"fullCacheHits"`
		InspectedCount             int `json:"inspectedCount"`
		IoTime                     int `json:"ioTime"`
		MatchCount                 int `json:"matchCount"`
		MaxInspectedCount          int `json:"maxInspectedCount"`
		MergeTime                  int `json:"mergeTime"`
		MinInspectedCount          int `json:"minInspectedCount"`
		OmittedCount               int `json:"omittedCount"`
		PartialCacheHits           int `json:"partialCacheHits"`
		ProcessCount               int `json:"processCount"`
		RawBytes                   int `json:"rawBytes"`
		ResponseBodyBytes          int `json:"responseBodyBytes"`
		RunningQueriesTotal        int `json:"runningQueriesTotal"`
		SlowLaneFileProcessingTime int `json:"slowLaneFileProcessingTime"`
		SlowLaneFiles              int `json:"slowLaneFiles"`
		SlowLaneWaitTime           int `json:"slowLaneWaitTime"`
		WallClockTime              int `json:"wallClockTime"`
	} `json:"performanceStats"`
	Results []struct {
		Count int64 `json:"count"`
	} `json:"results"`
}

func main() {
	interval, _ := strconv.ParseInt(os.Getenv("DOGSIGHTS_INTERVAL"), 10, 64)
	api_key := os.Getenv("INSIGHTS_API_KEY")

	for {
		// Get Response from Insights API
		queries := map[string]string{
			"not_wifi":      "SELECT filter(count(*), WHERE sleep_code='notWifi') FROM AdServerEvents",
			"pub_not_found": "SELECT filter(count(*), WHERE sleep_code='pubNotFound') FROM AdServerEvents",
			"pub_dvl":       "SELECT filter(count(*), WHERE sleep_code='pubDvl' ) FROM AdServerEvents",
			"SELECT filter(count(*), WHERE sleep_code='filtersRemovedAllCampaigns' ) FROM AdServerEvents",
			"SELECT filter(count(*), WHERE sleep_code='exchangeRemovedAllCampaigns' ) FROM AdServerEvents",
			"SELECT filter(count(*), WHERE sleep_code='serverError' ) FROM AdServerEvents",
			"SELECT filter(count(*), WHERE sleep_code='inactivePub') FROM AdServerEvents",
			"SELECT filter(count(*), WHERE sleep_code='tooBusy') FROM AdServerEvents",
		}
		for _, nrql := range queries {
			account := "752957"
			domain := "https://insights-api.newrelic.com"
			path := "/v1/accounts/" + account + "/query?"
			q := &url.URL{Path: nrql}
			nrql := q.String()
			params := "nrql=" + nrql
			full_url := domain + path + params
			impressions, _ := sh.Command("curl", "-s", "-H", "Accept: application/json", "-H", "X-Query-Key: "+api_key, full_url).Output()
			// Store Response as Struct
			var insights Insights
			json.Unmarshal(impressions, &insights)
			v := insights.Results
			fmt.Printf("%v", v)
			// Submit Count as Metric
			g, err := godspeed.NewDefault()

			if err != nil {
				// handle error
			}

			defer g.Conn.Close()

			err = g.Gauge(fmt.Sprintf("dogsights.%v"), 1, nil)

			if err != nil {
				// handle error
			}
			//OLD SUBMIT
			client, _ := statsd.NewClient(statsd_host, "impressions")
			stat := "impressions"
			fmt.Println(fmt.Sprint(stat, ":", v, "|c"))
			errr := client.Inc(stat, v, 1)
			if errr != nil {
				fmt.Println("Error sending metric: %+v", errr)
			}
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
