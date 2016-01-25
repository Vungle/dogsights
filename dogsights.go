package main

import (
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/codeskyblue/go-sh"
	//"net/url"
	"encoding/json"
	"fmt"
	"os"
	//"strconv"
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

	api_key := os.Getenv("INSIGHTS_API_KEY")

	// Get Response from Insights API
	url := "https://insights-api.newrelic.com/v1/accounts/752957/query?nrql=SELECT%20count(*)%20FROM%20AdServerEvents%20WHERE%20vungleType%3D%27reportAd%27%20and%20pub_app_id%3D%20%27com.cmplay.tiles2%27"
	impressions, _ := sh.Command("curl", "-H", "Accept: application/json", "-H", "X-Query-Key: "+api_key, url).Output()

	// Store Response as Struct
	var insights Insights
	json.Unmarshal(impressions, &insights)
	v := insights.Results[0].Count

	// Submit Count as Metric to StatsD
	client, _ := statsd.NewClient("localhost:8125", "impressions")
	stat := "impressions"
	fmt.Println(fmt.Sprint(stat, ":", v, "|c"))
	errr := client.Inc(stat, v, 1)
	if errr != nil {
		fmt.Println("Error sending metric: %+v", errr)
	}
}
