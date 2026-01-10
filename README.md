# CyberPower UPS (pwrstat) Prometheus Exporter

This is a simple server that read current status from CLI `pwrstat` and exports them via HTTP for Prometheus consumption.

# Installation

> [!NOTE]
> The pwrstat CLI and the pwrstatd daemon are responsible for monitoring and controlling the UPS and cannot be separated. Installing pwrstat in a Docker container is a bad decision, because it allows you to manage the UPS, while the exporter should only be responsible for receiving the UPS status. In any case, pwrstat must be installed on the host machine to manage the UPS. Therefore, only the file containing the output of pwrstat -status should be read.

### Write UPS status

> [!NOTE]
> I use systemd to write the output of pwrstat -status to a file. If you don't want to use systemd, you can do it in any way that is convenient for you.

1. Create directory to store status output
```sh
sudo mkdir /var/lib/pwrstat_status
```

2. Create service and timer

```sh
# /etc/systemd/system/pwrstat_status.service
[Unit]
Description=UPS status reader

[Service]
Type=oneshot
StandardOutput=file:/var/lib/pwrstat_status/status
ExecStart=/usr/sbin/pwrstat -status
```

```sh
# /etc/systemd/system/pwrstat_status.timer
[Unit]
Description=Export pwrstat data every 5 sec

[Timer]
OnBootSec=5s
OnUnitActiveSec=5s

AccuracySec=1s

[Install]
WantedBy=timers.target
```

3. Reload daemon and enable timer

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now pwrstat_status.timer

# check status
sudo systemctl status pwrstat_status.service
```

### Run docker compose

```yaml
services:
  pwrstat_exporter:
    image: ghcr.io/milden6/pwrstat-exporter:latest
    container_name: pwrstat-exporter
    restart: unless-stopped
    user: 1000:1000
    ports:
      - 9101:9101
    volumes:
      - /var/lib/pwrstat_status/status:/var/lib/pwrstat_status/status:ro
```

### Grafana dashboard (example)
![](/static/grafana_dashboard.png)