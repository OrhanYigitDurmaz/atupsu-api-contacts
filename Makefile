.PHONY: build build-32 install-svc uninstall-svc

build:
	GOOS=windows GOARCH=amd64 go build -o atupsu-api-amd64.exe .

build-32:
	GOOS=windows GOARCH=386 go build -o atupsu-api-386.exe .

install-svc:
	sc create AtupsuAPI binPath= "C:\atupsu-api\atupsu-api-amd64.exe" start= auto
	sc description AtupsuAPI "Atupsu REST API for atilim.mdb"

uninstall-svc:
	sc stop AtupsuAPI
	sc delete AtupsuAPI
