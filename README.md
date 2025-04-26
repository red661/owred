# OWNSA web service


### Hardware
* SoC: SigmaStar SSD201 ARMv7 Processor rev 5 (v7l)

    `half thumb fastmult vfp edsp thumbee neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm`

* Memory: 47812 kB

* Flash disk: 32 MiB


### Software
* Programming language: go1.22.3 windows/amd64

* CGO Compiler: Linaro GCC 7.5-2019.12

* vue-pure-admin: 5.6.0
    - nodejs: v20.13.1
    - pnpm: 9.1.1


### Build
* Windows
    - download msys2 and install mingw64-gcc
        - download msys2: https://mirrors.ustc.edu.cn/msys2/distrib/x86_64/
        - sed -i "s#mirror.msys2.org/#mirrors.ustc.edu.cn/msys2/#g" /etc/pacman.d/mirrorlist*
        - pacman -Syu
        - pacman -Su
        - pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-gdb
        - add "D:\Tools\SDK\msys64\mingw64\bin" to Path
    - .\build.win.ps1
    - .\build\ownsa.exe

* SSD201
    - download ARM GCC: https://releases.linaro.org/components/toolchain/binaries/7.5-2019.12/arm-linux-gnueabihf/gcc-linaro-7.5.0-2019.12-i686-mingw32_arm-linux-gnueabihf.tar.xz
    - .\build.ps1
# owred
