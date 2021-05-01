.PHONY: lint test vendor clean

export GO111MODULE=on

default: lint test

lint:
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor

coverage:
	rm profile.cov cover.html && go test -v -coverpkg=./... -coverprofile=profile.cov ./... && go tool cover -html=profile.cov -o cover.html