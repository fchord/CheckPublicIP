package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type cfDNSResult struct {
	Content    string `json:"content"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ID         string `json:"id"`
	ZoneID     string `json:"zone_id"`
	ZoneName   string `json:"zone_name"`
	TTL        int    `json:"ttl"`
	Comment    string `json:"comment"`
	CreatedOn  string `json:"created_on"`
	ModifiedOn string `json:"modified_on"`
}

type cfDNSResp struct {
	Result  cfDNSResult    `json:"result"`
	Success bool           `json:"success"`
	Errors  []interface{}  `json:"errors"`
	Messages []interface{} `json:"messages"`
}

func cfClient() *http.Client {
	return &http.Client{Timeout: 20 * time.Second}
}

func getDNS(cfg Config) (string, string, error) {
	url := fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s",
		cfg.Cloudflare.ZoneID,
		cfg.Cloudflare.DNSRecordID,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Cloudflare.APIToken)

	resp, err := cfClient().Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var parsed cfDNSResp
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", string(body), err
	}
	if !parsed.Success {
		return "", string(body), fmt.Errorf("cloudflare get dns failed")
	}
	return parsed.Result.Content, string(body), nil
}

func updateDNS(cfg Config, ip string) (string, error) {
	url := fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s",
		cfg.Cloudflare.ZoneID,
		cfg.Cloudflare.DNSRecordID,
	)
	payload := map[string]interface{}{
		"content": ip,
		"name":    cfg.Cloudflare.Domain,
		"proxied": cfg.Cloudflare.Proxied,
		"type":    "A",
		"comment": "DDNS update",
		"ttl":     cfg.Cloudflare.TTL,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Cloudflare.APIToken)

	resp, err := cfClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(body), err
	}

	var parsed cfDNSResp
	if err := json.Unmarshal(body, &parsed); err != nil {
		return string(body), err
	}
	if !parsed.Success {
		return string(body), fmt.Errorf("cloudflare update dns failed")
	}
	return string(body), nil
}
