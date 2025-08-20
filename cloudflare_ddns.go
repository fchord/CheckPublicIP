package main

import (
	//"fmt"
	"os/exec"
	"encoding/json"
	"errors"
)

type CfDnsRespResult struct {
	Content string `json:"content"`
	Name string `json:"name"`
	Type string `json:"type"`
	Id string `json:"id"`
	ZoneId string `json:"zone_id"`
	ZoneName string `json:"zone_name"`
	Ttl int `json:"ttl"`
	Comment string `json:"comment"`
	CreatedOn string `json:"created_on"`
	ModifiedOn string `json:"modified_on"`
}

type CfDnsResp struct  {
	Result CfDnsRespResult `json:"result"`
	Success bool `json:"success"`
	Message []string `json:"message"`
	Errors []string `json:"errors"`
}

var API_KEY = "fszF17L6VwWoLp3u7-gdKe5FFHZkoxBC6VvvylCh"
var ZONE_ID = "e70c1a8458c0ecd393004ed7ff32e93d"
var DOMAIN = "gohome.homes"
var DNS_RECORD_ID = "862957b4835aafb0f06fbfa2985ed33c"


func getDns() (string, string, error) {
	url := "https://api.cloudflare.com/client/v4/zones/" + ZONE_ID + "/dns_records/" + DNS_RECORD_ID

    curl := exec.Command(
		"curl", "--request", "GET", 
		"--url", url,
		"--header", "Content-Type: application/json",
		"--header", "Authorization: Bearer " + API_KEY)
    out, err := curl.Output()
	// fmt.Println("Get DNS: ", string(out))
    if err != nil {
        // fmt.Println("erorr", err)
        return "", "", err
    }
	var resp CfDnsResp
	json.Unmarshal([]byte(string(out)), &resp)
/* 	fmt.Println("Success: ", resp.Success)
	fmt.Println("Message: ", resp.Message)
	fmt.Println("Errors: ", resp.Errors)
	fmt.Println("Content: ", resp.Result.Content)
	fmt.Println("Name: ", resp.Result.Name)
	fmt.Println("Type: ", resp.Result.Type)
	fmt.Println("Id: ", resp.Result.Id)
	fmt.Println("ZoneId: ", resp.Result.ZoneId)
	fmt.Println("ZoneName: ", resp.Result.ZoneName)
	fmt.Println("Ttl: ", resp.Result.Ttl)
	fmt.Println("Comment: ", resp.Result.Comment)
	fmt.Println("CreatedOn: ", resp.Result.CreatedOn)
	fmt.Println("ModifiedOn: ", resp.Result.ModifiedOn) */

	if true == resp.Success {
		return resp.Result.Content, string(out), nil
	} else {
		return "", string(out), errors.New("Failed")
	}
}

func updateDns(IP string) (string, error){
	url := "https://api.cloudflare.com/client/v4/zones/" + ZONE_ID + "/dns_records/" + DNS_RECORD_ID
	//fmt.Println("url: ", url)
	data := "{\"content\": \"" + IP + "\"," +
		 "\"name\": \"" + DOMAIN + "\"," +
		 "\"proxied\": false," +
		 "\"type\": \"A\"," +
		 "\"comment\": \"Domain verification record\"," +
		 "\"id\": \"" + ZONE_ID + "\"," +
		 "\"ttl\": 60}"

	//fmt.Println("data: ", data)
    curl := exec.Command(
		"curl", "--request", "PATCH", 
		"--url", url,
		"--header", "Content-Type: application/json",
		"--header", "Authorization: Bearer " + API_KEY,
		"--data", data)
    out, err := curl.Output()
	//fmt.Println(string(out))
    if err != nil {
        //fmt.Println("erorr", err)
        return string(out), errors.New("Failed")
    }
    return string(out), nil
}

/* func main() {
	ip, err := getDns()
	if err == nil {
		fmt.Println("ip: ", ip)
	} 
	updateDns("113.92.159.93")
} */