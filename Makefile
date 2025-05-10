Commit ?= $(shell git rev-parse HEAD)
Branch ?= $(shell git branch --show-current)
Source ?= $(shell git remote get-url --push origin)
Built  ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
Tag    ?= $(shell git describe --tags --exact-match 2>/dev/null || echo 'none')

all:
	go build -o lutefisk --ldflags "-X main.GitCommit=${Commit} -X main.GitBranch=${Branch} -X main.GitTag=${Tag} -X main.BuildTime=${Built} -X main.GitSource=${Source}"

.phoney: all
