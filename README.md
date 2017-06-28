# OpenVPN DNS sync

Tool that parse a openvpn status file and store to hosts name/ip pairs as hosts-like file for dnsmasq or push to [a dynamic DNS](https://github.com/muka/ddns)

## Run with Docker

```bash
docker run -v `pwd`/data:/data -v /tmp/openvpn-status.log:/tmp/openvpn-status.log raptorbox/ovpndns-amd64 -s /tmp/openvpn-status.log -o /data/hosts
```

## Options

- `--src` (`OVPN_STATUS_FILE`) Set the openvpn status file source
- `--out` (`OUT_FILE`) Set the output file of a hosts-like formatted list of clients, let empty to disable
- `--domain` (`DOMAIN`) Set the default domain to append to each host name
- `--ddns` (`DDNS`) Enable DDNS sync
- `--ddns-host` (`DDNS_HOST`) ddns gRPC port in format `host:port`, let empty to disable
- `--debug` (`DEBUG`) Enable debugging logs

## Development Setup

`go get ./...`

## License

MIT License
