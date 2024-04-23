# start the app
app:
	go run service.payment/cmd/api.go
# test
test:
	- go clean -testcache
	- go test -race ./...
# code check
lintcheck:
	- golangci-lint run ./...
	- staticcheck ./...
# docker
docker:
	docker-compose up
