@SET GOOS=linux
@SET GOARCH=amd64
@cd ..
del go2rtc
go build -ldflags "-s -w" -trimpath