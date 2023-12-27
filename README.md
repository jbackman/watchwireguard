# watchwireguard

Tool for ensuring wireguard can reconnect to a dynamic server and bring up the tunnel

To compile:

go build watchwg.go

To compile arm64:

env GOOS=linux GOARCH=arm64 go build -o watchwg-arm64 watchwg.go
