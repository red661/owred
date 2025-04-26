$Env:Path = "D:\SDK\ARMLinuxGCC7\bin;" + $Env:Path

$Env:GOARM=7
$Env:GOARCH="arm"
$Env:GOOS="linux"
$Env:CGO_ENABLED=1

$Env:CC="arm-linux-gnueabihf-gcc"
$Env:CXX="arm-linux-gnueabihf-g++"

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

go build -ldflags="-s -w" -o build/ownsa cmd/gen.go cmd/app.go
