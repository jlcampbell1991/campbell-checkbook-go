# Setting Up systemd
## Create systemd service
`sudo nano /etc/systemd/system/campbell-checkbook-go.service`
```
[Unit]
Description=campbell-checkbook-go

[Service]
Type=simple
User=jcampbell
ExecStart=bash /home/jcampbell/go/campbell-checkbook-go/start.sh

[Install]
WantedBy=multi-user.target
```
`sudo systemctl enable --now campbell-checkbook-go`
`sudo journalctl -xfu campbell-checkbook-go`

alias campbell-checkbook-go-logs="sudo journalctl -xfu campbell-checkbook-go"
alias campbell-checkbook-go="sudo systemctl restart campbell-checkbook-go && campbell-checkbook-go-logs"