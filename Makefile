include .env

Author := yefriddavid
Version := $(shell git describe --abbrev=0 --tags | head -1) ##$(shell date "+%Y%m%d%H%M")
ReleaseDate := $(shell date "+%Y/%m/%d-%H:%M")
GitCommit := $(shell git rev-parse HEAD)
GitShortCommit := $(shell git rev-parse --short HEAD)
SysConfigFile := $(SYS_CONFIG_FILE)
DevConfigFile := $(DEV_CONFIG_FILE)


LDFLAGS := "-s -w -X main.SysConfigFile=$(SysConfigFile) -X main.Version=$(Version) -X main.GitCommit=$(GitCommit) -X main.Author=$(Author) -X main.GitShortCommit=$(GitShortCommit) -X main.ReleaseDate='$(ReleaseDate)'"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ".:*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run:
run: ## run dev local
	@go run -ldflags $(LDFLAGS) cmd/main.go --configFile=$(DevConfigFile)

local-publish: copy-local-config
local-publish:
	@go build -ldflags $(LDFLAGS) cmd/main.go
	sudo mv main /usr/local/bin/refreshSshConfigHosts

copy-local-config:
copy-local-config:
	cp ./config.yml $(SysConfigFile)

use:
use:
	gvm use go1.15.2

show:
show:
	@echo $(SysConfigFile)
	@echo $(DevConfigFile)
