
Pour Builder
```shell
set GOOS=linux
set GOARCH=arm
set GOARM=5
go build -o out/gotemp
```

ou àvec just
`just build`



## Installation du service

 Installation
```shell
sudo cp /usr/rep/gotemp.service /etc/systemd/system/gotemp.service
sudo systemctl enable gotemp.service
sudo systemctl start gotemp.service
```

Voir l'état
```shell
sudo systemctl status gotemp.service
```

Consulter les logs
```shell
journalctl -u gotemp.service
```



pour éteindre la lumiere du raspberry pi zero 2:
```shell
echo 0 | sudo tee /sys/class/leds/ACT/brightness
```
