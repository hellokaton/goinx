#!/bin/sh
# author: biezhi

unamestr=`uname`

SHA256='shasum -a 256'
if ! hash shasum 2> /dev/null
then
	SHA256='sha256sum.exe'
fi

VERSION=`date -u +%Y%m%d`
LDFLAGS="-s -w"
GCFLAGS=""

OSES=(linux windows)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
        cgo_enabled=0
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o goinx_${os}_${arch}${suffix} github.com/biezhi/goinx

    upx goinx_${os}_${arch}${suffix}
		tar -zcf goinx-${os}-${arch}-$VERSION.tar.gz goinx_${os}_${arch}${suffix}
		$SHA256 goinx-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags "$TAGS" -o goinx_linux_arm${v}  github.com/biezhi/goinx
done
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags "$TAGS" -o goinx_linux_arm64  github.com/biezhi/goinx
upx goinx_linux_arm*
tar -zcf goinx-linux-arm-$VERSION.tar.gz goinx_linux_arm*
$SHA256 goinx-linux-arm-$VERSION.tar.gz
