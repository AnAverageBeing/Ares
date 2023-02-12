export GOARCH=amd64
export GOOS=linux
go build -o targets/Ares-linux
export GOOS=windows
go build -o targets/Ares-windows.exe