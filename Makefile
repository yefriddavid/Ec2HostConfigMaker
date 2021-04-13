include .env

Author := yefriddavid
Version := $(shell git describe --abbrev=0 --tags | head -1)
ReleaseDate := $(shell date "+%Y/%m/%d-%H:%M")
GitCommit := $(shell git rev-parse HEAD)
GitShortCommit := $(shell git rev-parse --short HEAD)
SysConfigFile := $(SYS_DEFAULT_TARGET_CONFIG_FILE)
DevConfigFile := $(DEV_SOURCE_CONFIG_FILE)


LDFLAGS := "-o refreshSshConfigHosts -s -w -X main.SysConfigFile=$(SysConfigFile) -X main.Version=$(Version) -X main.GitCommit=$(GitCommit) -X main.Author=$(Author) -X main.GitShortCommit=$(GitShortCommit) -X main.ReleaseDate='$(ReleaseDate)'"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ".:*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run:
run: ## run dev local
	@go run -ldflags $(LDFLAGS) cmd/main.go --configFile=$(DevConfigFile)

freeze:
freeze:
	git add .
	git commit -m "freeze"
	git tag -d 0.0.1
	git tag 0.0.1

release:
release:
	goreleaser --skip-validate --skip-publish --rm-dist
#@go build -ldflags $(LDFLAGS) cmd/main.go


local-release: release copy-local-config
local-release:
	@echo "Success"

copy-local-app:
copy-local-app:
	sudo cp ./dist/Ec2SshConfigHostMaker_linux_amd64/Ec2SshConfigHostMaker /usr/local/bin/refreshSshConfigHosts
#sudo mv main /usr/local/bin/refreshSshConfigHosts
#@go build -ldflags $(LDFLAGS) cmd/main.go

copy-local-config:
copy-local-config:
	sudo cp ./configs/config.yml $(SysConfigFile)

use:
use:
	gvm use go1.15.2

show:
show:
	@echo $(SysConfigFile)
	@echo $(DevConfigFile)
