package goatcounter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type GoatMetrics struct {
	Configured bool   `json:"configured"`
	Verified   bool   `json:"verified"`
	BaseURL    string `json:"base_url,omitempty"`
	Views      int    `json:"views"`
	Visitors   int    `json:"visitors"`
	Error      string `json:"error,omitempty"`
}

func Collect() (GoatMetrics, []string) {
	gc := GoatMetrics{}
	notes := []string{}
	gcBase := os.Getenv("GOATCOUNTER_URL")
	gcToken := os.Getenv("GOATCOUNTER_TOKEN")
	if gcBase != "" {
		gc.Configured = true
		gc.BaseURL = gcBase
		v, vis, err := fetchGoatCounter(gcBase, gcToken)
		if err != nil {
			gc.Error = err.Error()
			notes = append(notes, "goatcounter fetch failed: "+err.Error())
		} else {
			gc.Verified = true
			gc.Views = v
			gc.Visitors = vis
		}
	} else {
		notes = append(notes, "GOATCOUNTER_URL not set; traffic unavailable")
	}
	return gc, notes
}

func CollectJSON() ([]byte, error) {
	gc, notes := Collect()
	wrapper := map[string]any{
		"metrics": gc,
		"notes":   notes,
	}
	return json.MarshalIndent(wrapper, "", "  ")
}

func Verify() error {
	base := os.Getenv("GOATCOUNTER_URL")
	token := os.Getenv("GOATCOUNTER_TOKEN")
	if base == "" {
		return fmt.Errorf("GOATCOUNTER_URL not set")
	}
	_, _, err := fetchGoatCounter(base, token)
	return err
}

func fetchGoatCounter(base, token string) (int, int, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	url := strings.TrimRight(base, "/") + "/api/v0/stats/total"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("http %d", resp.StatusCode)
	}
	var parsed map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, 0, err
	}
	views := 0
	visitors := 0
	if v, ok := parsed["total"]; ok {
		views = toInt(v)
	}
	if v, ok := parsed["views"]; ok {
		views = toInt(v)
	}
	if v, ok := parsed["visitors"]; ok {
		visitors = toInt(v)
	}
	return views, visitors, nil
}

func toInt(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	case string:
		var n int
		fmt.Sscanf(t, "%d", &n)
		return n
	}
	return 0
}