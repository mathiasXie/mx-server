package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IPInfo struct {
	IP          string  `json:"ip"`
	Success     bool    `json:"success"`
	Country     string  `json:"country"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CallingCode float64 `json:"calling_code"`
}

func GetIPLocation(ip string, lang string) (*IPInfo, error) {
	if lang == "" {
		lang = "zh-CN" // 默认为中文
	}
	url := fmt.Sprintf("https://ipwho.is/%s?lang=%s", ip, lang)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	if !info.Success {
		return nil, fmt.Errorf("failed to lookup IP: %s", ip)
	}
	return &info, nil
}
