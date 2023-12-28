.PHONY: build
build: build-cli build-gui

.PHONY: build-cli
build-cli:
	go build -v -o sonata ./cli

.PHONY: build-gui
build-gui:
	go build -v -o sonata-gui ./gui

# (build but with a smaller binary)
.PHONY: dist-gui
dist-gui:
	go build -gcflags=all=-l -v -ldflags="-w -s" -o sonata-gui ./gui/sonata

.PHONY: dist-cli
dist-cli:
	go build -gcflags=all=-l -v -ldflags="-w -s" -o sonata ./cli/sonata

.PHONY: cross-build
cross-build:
	docker build -f ./build/win/Dockerfile -t sonata-win .
	docker rm -f sonata-win || true
	docker run --rm --name sonata-win -d sonata-win:latest tail -f /dev/null
	rm -rf ./dist
	docker cp sonata-win:/root/go/src/sonata/dist ./dist
	docker rm -f sonata-win

.PHONY: run
run: build
	./sonata-gui

.PHONY: test
test: 
	go test -cover -parallel 5 -failfast -count=1 ./... 

# human readable test output
.PHONY: love
love:
ifeq ($(filter watch,$(MAKECMDGOALS)),watch)
	gotestsum --watch ./...
else
	gotestsum ./...
endif

.PHONY: tidy
tidy:
	go mod tidy

# auto restart
.PHONY: dev
dev:
	air

.PHONY: lint
lint:
	revive -formatter friendly -config revive.toml ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: gosec
gosec:
	gosec -tests ./... 

.PHONY: inspect
inspect: lint gosec staticcheck

.PHONY: install-inspect-tools
install-inspect-tools:
	go install github.com/mgechev/revive@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
