GIT_VER := $(shell git describe --tags)
DATE    := $(shell date +%Y-%m-%dT%H:%M:%S%z)
OWNER   := "mix3"
REPO    := "unity3d2png"

clean:
	rm -rf pkg/*

binary: clean
	cd cmd/unity3d2png && \
	CGO_ENABLED=0 gox -osarch="linux/amd64 darwin/amd64" \
		-output "../../pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" \
		-ldflags "-w -s -X main.version=${GIT_VER} -X main.buildDate=${DATE}"

package: binary
	cd ./pkg && find . -name "*${GIT_VER}*" -type f \
		-exec mkdir -p ${REPO} \; \
		-exec cp {} ${REPO}/${REPO} \; \
		-exec zip -r {}.zip ${REPO} \; \
		-exec rm -rf ${REPO} \;