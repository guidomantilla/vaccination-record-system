.PHONY: phony
phony-goal: ; @echo $@

build: validate
	docker compose -f docker/docker-compose.yml -p vaccination-record-system  up --detach

validate: generate sort-import format vet lint coverage

generate:
	go generate ./...

sort-import:
	goimports-reviser -rm-unused -set-alias -format -recursive cmd
	goimports-reviser -rm-unused -set-alias -format -recursive pkg

format:
	go fmt ./...

vet:
	go vet ./cmd/... ./pkg/...

lint:
	golangci-lint run ./cmd/... ./pkg/...

test:
	go test -covermode count -coverprofile coverage.out.tmp.01 ./pkg/...
	cat coverage.out.tmp.01 | grep -v "mocks.go" > coverage.out

coverage: test
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

sonarqube: coverage
	sonar-scanner

update-dependencies:
	go get -u ./...
	go get -t -u ./...
	go mod tidy

prepare:
	go install github.com/kisielk/godepgraph@latest
	go install github.com/incu6us/goimports-reviser/v3@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/cweill/gotests/gotests@latest
	go mod download
	go mod tidy
