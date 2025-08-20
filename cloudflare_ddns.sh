#!/bin/bash

#EMAIL="doraemon_hates_mice@protonmail.com"
API_KEY="fszF17L6VwWoLp3u7-gdKe5FFHZkoxBC6VvvylCh"
ZONE_ID="e70c1a8458c0ecd393004ed7ff32e93d"
DOMAIN="gohome.homes"
DNS_RECORD_ID="862957b4835aafb0f06fbfa2985ed33c"

IP=$(curl --request GET --url https://api.ipify.org 2> /dev/null)
#IP="113.92.159.93"
echo IP:$IP

url="https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records/$DNS_RECORD_ID"
echo url:$url

curl --request PATCH \
    --url  $url\
    --header "Content-Type: application/json" \
    --header "Authorization: Bearer $API_KEY" \
    --data "{
    \"content\": \"$IP\",
    \"name\": \"$DOMAIN\",
    \"proxied\": false,
    \"type\": \"A\",
    \"comment\": \"Domain verification record\",
    \"id\": \"$ZONE_ID\",
    \"ttl\": 60
}" 

#curl -X PUT "https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records/$DOMAIN" \
#     -H "Authorization: Bearer $API_KEY" \
#     -H "Content-Type: application/json" \
#     --data "{\"type\":\"A\",\"name\":\"$DOMAIN\",\"content\":\"$IP\",\"ttl\":120,\"proxied\":false}"