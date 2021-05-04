.PHONY: lint test vendor clean coverage demo

export GO111MODULE=on

default: lint test

lint:
	golangci-lint run

test:
	rm -fr .cache | true
	mkdir .cache
	go test -v -cover ./...

yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor

coverage:
	rm profile.cov cover.html && go test -v -coverpkg=./... -coverprofile=profile.cov ./... && go tool cover -html=profile.cov -o cover.html

demo-init:
	docker network create traefik-net | true

demo-up: demo-init
	docker-compose --env-file=demo/.env  -f demo/gateway/docker-compose.yml -f demo/frontend/docker-compose.yml -f demo/imaginary/docker-compose.yml up -d

demo-restart: demo-init
	docker-compose --env-file=demo/.env  -f demo/gateway/docker-compose.yml -f demo/frontend/docker-compose.yml -f demo/imaginary/docker-compose.yml restart

demo-logs: demo-up
	docker-compose --env-file=demo/.env  -f demo/gateway/docker-compose.yml -f demo/frontend/docker-compose.yml -f demo/imaginary/docker-compose.yml logs -f gateway imaginary