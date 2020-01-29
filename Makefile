.PHONY all:
	env GOOS=linux GOARCH=arm GOARM=6 go build -trimpath -ldflags="-s -w"
	cp go-ble-mqtt /Users/space/Dropbox/
	md5sum /Users/space/Dropbox/go-ble-mqtt
