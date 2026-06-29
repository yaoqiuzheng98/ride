package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ride/pkg/config"
)

const (
	baseURL = "https://restapi.amap.com/v5/place/text"
)

// POI is a single point of interest returned by amap.
type POI struct {
	Name     string `json:"name"`
	Location string `json:"location"` // "longitude,latitude"
	Address  string `json:"address"`
	Adname   string `json:"adname"` // 区
	Cityname string `json:"cityname"`
	Type     string `json:"type"`
	Typecode string `json:"typecode"`
}

// Result is a POI with computed distance (meters) from the origin.
type Result struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	Adname    string  `json:"adname"`
	Cityname  string  `json:"cityname"`
	Type      string  `json:"type"`
	Distance  float64 `json:"distance"` // meters
}

type v5Response struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Pois     []POI  `json:"pois"`
}

var client = &http.Client{Timeout: 10 * time.Second}

// SearchByDistance searches keyword around the given origin (lat, lng) and
// returns results sorted by distance ascending (nearest first).
// Uses the high-precision POI 2.0 API (restapi.amap.com/v5) which returns
// 6-decimal-place coordinates.
func SearchByDistance(lat, lng float64, keyword string) ([]Result, error) {
	cfg := config.Get().Amap
	if cfg.Key == "" {
		return nil, fmt.Errorf("amap: key not configured (set amap.key in config.toml)")
	}
	region := cfg.DefaultRegion
	if region == "" {
		region = "440100"
	}

	pois, err := searchAll(keyword, region, cfg.Key)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(pois))
	for _, p := range pois {
		lng2, lat2, ok := parseLocation(p.Location)
		if !ok {
			continue
		}
		results = append(results, Result{
			Name:      p.Name,
			Latitude:  lat2,
			Longitude: lng2,
			Address:   p.Address,
			Adname:    p.Adname,
			Cityname:  p.Cityname,
			Type:      p.Type,
			Distance:  haversine(lat, lng, lat2, lng2),
		})
	}
	sortByDistance(results)
	return results, nil
}

// Search returns all POIs matching the keyword (no distance sorting).
func Search(keyword string) ([]POI, error) {
	cfg := config.Get().Amap
	if cfg.Key == "" {
		return nil, fmt.Errorf("amap: key not configured (set amap.key in config.toml)")
	}
	region := cfg.DefaultRegion
	if region == "" {
		region = "440100"
	}
	return searchAll(keyword, region, cfg.Key)
}

// searchAll paginates through results (v5 API: max 200 total, page_size up to 25).
func searchAll(keyword, region, key string) ([]POI, error) {
	var pois []POI
	for pageNum := 1; pageNum <= 8; pageNum++ {
		batch, count, err := queryPage(keyword, region, key, pageNum, 25)
		if err != nil {
			return nil, err
		}
		pois = append(pois, batch...)
		if len(pois) >= count || len(batch) < 25 {
			break
		}
	}
	return pois, nil
}

func queryPage(keyword, region, key string, pageNum, pageSize int) ([]POI, int, error) {
	params := url.Values{
		"key":          {key},
		"keywords":     {keyword},
		"region":       {region},
		"region_limit": {"true"},
		"page_size":    {strconv.Itoa(pageSize)},
		"page_num":     {strconv.Itoa(pageNum)},
		"show_fields":  {"children,business,tel,photos"},
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var r v5Response
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, 0, fmt.Errorf("amap: parse response: %w (body=%s)", err, string(body))
	}
	if r.Status != "1" {
		return nil, 0, fmt.Errorf("amap: search failed status=%s info=%s code=%s", r.Status, r.Info, r.Infocode)
	}
	count, _ := strconv.Atoi(r.Count)
	return r.Pois, count, nil
}

func parseLocation(loc string) (lng, lat float64, ok bool) {
	parts := strings.Split(loc, ",")
	if len(parts) != 2 {
		return 0, 0, false
	}
	lng, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, false
	}
	lat, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, false
	}
	return lng, lat, true
}

// haversine returns the great-circle distance between two points in meters.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	dPhi := (lat2 - lat1) * math.Pi / 180
	dLambda := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dPhi/2)*math.Sin(dPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*math.Sin(dLambda/2)*math.Sin(dLambda/2)
	return 2 * R * math.Asin(math.Sqrt(a))
}

func sortByDistance(r []Result) {
	for i := 1; i < len(r); i++ {
		for j := i; j > 0 && r[j-1].Distance > r[j].Distance; j-- {
			r[j-1], r[j] = r[j], r[j-1]
		}
	}
}
