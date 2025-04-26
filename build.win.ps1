$Env:Path = "C:\Tools\SDK\msys64\mingw64\bin;" + $Env:Path

$Env:GOARCH="amd64"
$Env:GOOS="windows"
$Env:CGO_ENABLED=1

$Env:CC="gcc"
$Env:CXX="g++"

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

go build -o build/ownsa.exe cmd/gen.go cmd/app.go
