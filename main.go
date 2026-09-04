package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jordan-wright/email"
)

func initLog(cfg Config) {
	path := cfg.LogFile
	if path == "" {
		exe, err := os.Executable()
		if err != nil {
			log.SetOutput(os.Stderr)
			return
		}
		path = filepath.Join(filepath.Dir(exe), "message.log")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && !os.IsExist(err) {
		log.SetOutput(os.Stderr)
		log.Printf("log dir error: %v", err)
		return
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Printf("log file error: %v", err)
		return
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func getPublicIP(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	for {
		resp, err := client.Get(url)
		if err != nil {
			time.Sleep(30 * time.Second)
			continue
		}
		ip, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			time.Sleep(30 * time.Second)
			continue
		}
		return strings.TrimSpace(string(ip)), nil
	}
}

func sendEmail(cfg Config, publicIP string, when time.Time) error {
	if !cfg.Email.Enabled {
		return nil
	}
	e := email.NewEmail()
	e.From = cfg.Email.From
	e.To = cfg.Email.To
	e.Subject = cfg.Email.SubjectPrefix + publicIP
	e.Text = []byte(e.Subject + "\n\nquery time: " + when.String())
	addr := fmt.Sprintf("%s:%d", cfg.Email.SMTPHost, cfg.Email.SMTPPort)
	auth := smtp.PlainAuth("", cfg.Email.Username, cfg.Email.Password, cfg.Email.SMTPHost)
	return e.Send(addr, auth)
}

func work(cfg Config) {
	log.Println("Start.")
	defer log.Println("Quit.")

	interval := time.Duration(cfg.CheckIntervalMinutes) * time.Minute
	lastIP := ""

	for {
		currentDNS, result, err := getDNS(cfg)
		if err != nil {
			log.Println("getDNS err:", err)
			log.Println("getDNS result:", result)
			time.Sleep(time.Duration(cfg.CheckIntervalMinutes) * time.Second)
			continue
		}
		log.Println("current DNS A record:", currentDNS)
		break
	}

	for {
		now := time.Now()
		ip, err := getPublicIP(cfg.PublicIPURL)
		if err != nil {
			log.Println("getPublicIP err:", err)
			time.Sleep(interval)
			continue
		}
		log.Println("getPublicIP:", ip)

		currentDNS, _, err := getDNS(cfg)
		if err != nil {
			log.Println("getDNS err:", err)
		} else if currentDNS != ip {
			log.Println("Update DNS to", ip)
			res, err := updateDNS(cfg, ip)
			if err != nil {
				log.Println("updateDNS err:", err, "result:", res)
			} else {
				log.Println("updateDNS ok")
			}
		}

		if lastIP == ip {
			time.Sleep(interval)
			continue
		}
		if lastIP != "" {
			log.Println("Public ip changed:", ip)
		} else {
			log.Println("Public ip:", ip)
		}

		if cfg.Email.Enabled {
			if err := sendEmail(cfg, ip, now); err != nil {
				log.Println("sendEmail err:", err)
			} else {
				log.Println("Send email successfully.")
			}
		} else {
			log.Println("Email disabled; skip notify.")
		}

		lastIP = ip
		time.Sleep(interval)
	}
}

func main() {
	configPath := resolveConfigPath()
	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config %s: %v\n", configPath, err)
		os.Exit(1)
	}
	initLog(cfg)
	log.Println("using config:", configPath)
	log.Println("DDNS domain:", cfg.Cloudflare.Domain)
	log.Println("email enabled:", cfg.Email.Enabled)
	work(cfg)
}
