# systemd

Production unit-файлы NeoSync.

```bash
sudo cp http-systemd-amd64.service /etc/systemd/system/
sudo cp tcp-systemd-amd64.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now http-systemd-amd64.service
sudo systemctl enable --now tcp-systemd-amd64.service
```

Логи:

```bash
journalctl -u http-systemd-amd64 -f
journalctl -u tcp-systemd-amd64 -f
```
