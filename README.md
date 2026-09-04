# CheckPublicIP / Cloudflare DDNS

Poll public IP and update a Cloudflare A record. Optional email notify (off by default).

## Target

- Domain: `home.19121122.xyz`
- Workdir on LAN host: `/home/liuzhi/work/cloudflare_ddns`

## Config

```bash
cp config.example.yaml config.yaml
# edit config.yaml — put Cloudflare token / zone_id / dns_record_id there
# keep email.enabled: false unless you want mail
```

Never commit `config.yaml`.

## Build

```bash
cd /home/liuzhi/work/cloudflare_ddns
go mod tidy
go build -o check-public-ip .
```

## Run once (foreground)

```bash
./check-public-ip -config ./config.yaml
```

## systemd

```bash
sudo cp deploy/check-public-ip.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now check-public-ip
sudo systemctl status check-public-ip
journalctl -u check-public-ip -f
```

## Notes

- Old hardcoded secrets in earlier source must be rotated at Cloudflare / mailbox.
- `dns_record_id` is the Cloudflare DNS record id for the A record of `home.19121122.xyz`.
