ifeq ($(OS),Windows_NT)
	COMPILE_TIME := $(shell echo %date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%)
else
	COMPILE_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
endif

LDFLAGS := -s -w -X "main.Version=1.3.4" -X "main.Build=$(COMPILE_TIME)"

.PHONY: default all

default: all

go-bindata:
	go get -u github.com/go-bindata/go-bindata/...

.PHONY: assets
assets: go-bindata
	go-bindata -nomemcopy -pkg=assets -o=assets/assets.go \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-ignore=assets.go assets/...

.PHONY: payoutgenerator
payoutgenerator:
	go build -ldflags '$(LDFLAGS)' ./cmd/payoutgenerator

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

OS-ARCHS=linux:386 linux:amd64
TARGETS=generator duplicate

cross-compile:
	@$(foreach n, $(OS-ARCHS),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		$(foreach target, $(TARGETS),\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} \
			go build -trimpath -ldflags '$(LDFLAGS)' \
				-o ./release/$(target)_$${target_suffix} ./cmd/$(target);\
		)\
		echo "Build $${os}-$${arch} done";\
	)

all: fmt payoutgenerator tidy
