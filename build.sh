go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct
go build -ldflags="-s -w" -o build/ownsa.elf cmd/gen.go cmd/app.go

export PATH=/opt/armlinuxgcc7/bin:$PATH

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

export GOARM=7
export GOARCH="arm"
export GOOS="linux"
export CGO_ENABLED=1

export CC="arm-linux-gnueabihf-gcc"
export CXX="arm-linux-gnueabihf-g++"

export CC="/home/lanyang/red2/tmp/gcc-linaro-6.5.0-2018.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc"
export CXX="/home/lanyang/red2/tmp/gcc-linaro-6.5.0-2018.12-x86_64_arm-linux-gnueabihf/bin/arm-linux-gnueabihf-g++"

go build -ldflags="-s -w" -o build/ownsa.arm cmd/gen.go cmd/app.go
