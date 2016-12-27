# This how we want to name the binary output
BINARY=strongback
OUTPUT_DIR=out
WHAT=strongback.org/cli

# Read the version information from our file
VERSION=$(shell cat VERSION)
# Or could get from most recent tag reachable from the latest commit
#VERSION=`git describe --tags`

# Set the build date
BUILD_TIMESTAMP=`date +%FT%T%z`
BUILD_DATE=`date +"%Y-%m-%d"`

# Setup the -ldflags option for go build here
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Date=${BUILD_DATE} -X main.ExecName=${BINARY} -X main.Build=${BUILD_TIMESTAMP}"

#
# Cleans and builds everything
#
.PHONY: clean
clean:
	rm -rf ${OUTPUT_DIR}

.PHONY: all
all: osx linux windows

#
# Packaging for each OS
#
.PHONY: osx
osx: out/osx package-osx
	$(info )
	$(info OSX)

.PHONY: linux
linux: out/linux package-linux
	$(info )
	$(info Linux)

.PHONY: windows
windows: out/windows package-windows
	$(info )
	$(info Windows)

#
# Compile
#
out/osx: dependencies
	$(info OSX:     building executable)
	@GOOS="darwin" GOARCH="amd64" go build ${LDFLAGS} -o ${OUTPUT_DIR}/osx/strongback ${WHAT}

out/linux: dependencies
	$(info Linux:   building executable)
	@GOOS="linux" GOARCH="amd64" go build ${LDFLAGS} -o ${OUTPUT_DIR}/linux/strongback ${WHAT}

out/windows: dependencies
	$(info Windows: building executable)
	@GOOS="windows" GOARCH="386" go build ${LDFLAGS} -o ${OUTPUT_DIR}/windows/strongback.exe ${WHAT}

.PHONY: test
test:
	echo ${VERSION}

#
# Packaging 
#
package: package-osx package-linux package-windows

.PHONY: package-osx
package-osx: out/osx
	$(info OSX:     packaging archive)
	@tar -czf ${OUTPUT_DIR}/strongback-cli-${VERSION}-osx.tar.gz -C ${OUTPUT_DIR}/osx strongback
	@zip -r -j ${OUTPUT_DIR}/strongback-cli-${VERSION}-osx.zip ${OUTPUT_DIR}/osx > /dev/null

.PHONY: package-linux
package-linux: out/linux
	$(info Linux:   packaging archive)
	@tar -czf ${OUTPUT_DIR}/strongback-cli-${VERSION}-linux.tar.gz -C ${OUTPUT_DIR}/linux strongback
	@zip -r -j ${OUTPUT_DIR}/strongback-cli-${VERSION}-linux.zip ${OUTPUT_DIR}/linux > /dev/null

.PHONY: package-windows
package-windows: out/windows
	$(info Windows: packaging archive)
	@tar -czf ${OUTPUT_DIR}/strongback-cli-${VERSION}-windows.tar.gz -C ${OUTPUT_DIR}/windows strongback.exe
	@zip -r -j ${OUTPUT_DIR}/strongback-cli-${VERSION}-windows.zip ${OUTPUT_DIR}/windows > /dev/null

#
# Dependencies 
#
dependencies: ../../github.com/rickar/props

../../github.com/rickar/props:
	# Used to read Java properties files
	go get -u github.com/rickar/props
