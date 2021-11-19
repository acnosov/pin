echo "begin build"
CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o dist/pin.exe main.go
echo "done"
