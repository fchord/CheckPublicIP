package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type EmailConfig struct {
	Enabled       bool     `yaml:"enabled"`
	From          string   `yaml:"from"`
	To            []string `yaml:"to"`
	SubjectPrefix string   `yaml:"subject_prefix"`
	SMTPHost      string   `yaml:"smtp_host"`
	SMTPPort      int      `yaml:"smtp_port"`
	Username      string   `yaml:"username"`
	Password      string   `yaml:"password"`
}

type CloudflareConfig struct {
	APIToken    string `yaml:"api_token"`
	ZoneID      string `yaml:"zone_id"`
	DNSRecordID string `yaml:"dns_record_id"`
	Domain      string `yaml:"domain"`
	Proxied     bool   `yaml:"proxied"`
	TTL         int    `yaml:"ttl"`
}

type Config struct {
	CheckIntervalMinutes int              `yaml:"check_interval_minutes"`
	PublicIPURL          string           `yaml:"public_ip_url"`
	LogFile              string           `yaml:"log_file"`
	Email                EmailConfig      `yaml:"email"`
	Cloudflare           CloudflareConfig `yaml:"cloudflare"`
}

func defaultConfig() Config {
	return Config{
		CheckIntervalMinutes: 10,
		PublicIPURL:          "https://api64.ipify.org?format=text",
		Email: EmailConfig{
			Enabled:       false,
			SubjectPrefix: "Public IP: ",
			SMTPPort:      25,
		},
		Cloudflare: CloudflareConfig{
			Domain:  "home.19121122.xyz",
			Proxied: false,
			TTL:     60,
		},
	}
}

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.CheckIntervalMinutes <= 0 {
		cfg.CheckIntervalMinutes = 10
	}
	if cfg.PublicIPURL == "" {
		cfg.PublicIPURL = "https://api64.ipify.org?format=text"
	}
	if cfg.Cloudflare.Domain == "" {
		cfg.Cloudflare.Domain = "home.19121122.xyz"
	}
	if cfg.Cloudflare.TTL <= 0 {
		cfg.Cloudflare.TTL = 60
	}
	if cfg.Email.SMTPPort <= 0 {
		cfg.Email.SMTPPort = 25
	}
	if cfg.Cloudflare.APIToken == "" || cfg.Cloudflare.ZoneID == "" || cfg.Cloudflare.DNSRecordID == "" {
		return cfg, fmt.Errorf("cloudflare.api_token, zone_id, and dns_record_id are required in %s", path)
	}
	return cfg, nil
}

func resolveConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config.yaml")
	flag.Parse()
	if path != "" {
		return path
	}
	if _, err := os.Stat("/etc/check-public-ip/config.yaml"); err == nil {
		return "/etc/check-public-ip/config.yaml"
	}
	exe, err := os.Executable()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(filepath.Dir(exe), "config.yaml")
}
