package main

import (
	//"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	//"net/url"
    "log"
    "net/smtp"
	//"errors"
	"time"
	"os"
	"path/filepath"
    "github.com/jordan-wright/email"
)

// go build check_public_ip.go cloudflare_ddns.go

func initLog() {
		path, err := os.Executable()
		dir := filepath.Dir(path)		
        file := dir + "/" +"message"+ ".log"
        logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
        if err != nil {
                panic(err)
        }
        log.SetOutput(logFile) // 将文件设置为log输出的文件
        // log.SetPrefix("[qSkipTool]")
        log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
        return
}

func getPublicIP() (string, error) {
	
	for ;; {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    5 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Timeout: 10 * time.Second, Transport: tr}
		resp, err := client.Get("https://api64.ipify.org?format=text")
		if err != nil {
			// return "", err
			time.Sleep(time.Duration(30) * time.Second)
			continue					
		}
		defer resp.Body.Close()

		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//return "", err
			time.Sleep(time.Duration(30) * time.Second)
			continue					
		}
		return strings.TrimSpace(string(ip)), nil
	}
}


func sendEmail(public_ip_addr string, this_time time.Time) error {
    e := email.NewEmail()
    //设置发送方的邮箱
    e.From = "时代微尘 <hudaliuzhi@sohu.com>"
    // 设置接收方的邮箱
    e.To = []string{"hudaliuzhi@sohu.com"}
    //设置主题
    e.Subject = "公网地址： " + public_ip_addr
    //设置文件发送的内容
	content := e.Subject + "\n" + "\n"
	content += "查询时间： " + this_time.String()
    e.Text = []byte(content)
    //设置服务器相关的配置
	//密码是搜狐邮箱设置中的 “第三方客户端独立密码”
    err := e.Send("smtp.sohu.com:25", smtp.PlainAuth("", "hudaliuzhi@sohu.com", "GV1HDB53AQUK", "smtp.sohu.com"))
    if err != nil {
        log.Fatal(err)
		return err
    }
	return nil
}

func work() {
	log.Println("Start.")
	defer log.Println("Quit.")

	// 几分钟查一次
	interval := 10
	last_ip_addr := ""
	current_dns := ""
	for ;; {
		current_dns, result, err := getDns()
		if err != nil {
			log.Println("The first time, getDns err: ", err)
			log.Println("getDns result: ", result)
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			log.Println("The first time, getDns: ", current_dns)
			log.Println("getDns result: ", result)
			break
		}
	}
	for ;; {
		this_time := time.Now()
		ip_addr, err := getPublicIP()
		if err != nil {
			log.Println("Get public ip. err: ", err)
			time.Sleep(time.Duration(interval) * time.Minute)
			continue			
		}
		log.Println("getPublicIP: ", ip_addr)
		// Update DNS
		current_dns, _, err = getDns()
		if current_dns != ip_addr {
			log.Println("Update DNS to ", ip_addr)
			res, err := updateDns(ip_addr)
			if err != nil {
				log.Println("updateDns result: ", res)
			}
		}
		if last_ip_addr == ip_addr {
			time.Sleep(time.Duration(interval) * time.Minute)
			continue
		} else {
			if last_ip_addr != "" {
				log.Println("Public ip changed: ", ip_addr)
			} else {
				log.Println("Public ip: ", ip_addr)
			}
		}
		err = sendEmail(ip_addr, this_time)
		if err != nil {
			log.Println("err: ", err)
		} else {
			log.Println("Send successfully.")
		}

		last_ip_addr = ip_addr
	}
}

func main() {
	initLog()
	work()
	return 
}