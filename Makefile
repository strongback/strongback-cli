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
	$(info Cleaning 'out' directory)
	@rm -rf ${OUTPUT_DIR}

.PHONY: all
all: macos linux windows

#
# Packaging for each OS
#
.PHONY: macos
macos: out/macos package-macos

.PHONY: linux
linux: out/linux package-linux

.PHONY: windows
windows: out/windows package-windows

#
# Compile
#
out/macos: dependencies
	$(info macos:   building executable)
	@GOOS="darwin" GOARCH="amd64" go build ${LDFLAGS} -o ${OUTPUT_DIR}/macos/strongback ${WHAT}

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
package: package-macos package-linux package-windows

.PHONY: package-macos
package-macos: out/macos
	$(info macos:   packaging archive)
	@tar -czf ${OUTPUT_DIR}/strongback-cli-${VERSION}-macos.tar.gz -C ${OUTPUT_DIR}/macos strongback

.PHONY: package-linux
package-linux: out/linux
	$(info Linux:   packaging archive)
	@tar -czf ${OUTPUT_DIR}/strongback-cli-${VERSION}-linux.tar.gz -C ${OUTPUT_DIR}/linux strongback

.PHONY: package-windows
package-windows: out/windows
	$(info Windows: packaging archive)
	@zip -r -j ${OUTPUT_DIR}/strongback-cli-${VERSION}-windows.zip ${OUTPUT_DIR}/windows > /dev/null

#
# Dependencies 
#
dependencies: ../../github.com/rickar/props

../../github.com/rickar/props:
	# Used to read Java properties files
	go get -u github.com/rickar/props
