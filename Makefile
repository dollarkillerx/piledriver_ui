build_agent:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -o piledriver_linux -ldflags "-s -w" .
	#GOARCH=amd64 GOOS=darwin CGO_ENABLED=1 go build -o piledriver_darwin -ldflags "-s -w" .
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o piledriver.exe -i -ldflags="-s -w -H windowsgui"
	upx piledriver_linux
	#upx piledriver_darwin
	upx piledriver.exe