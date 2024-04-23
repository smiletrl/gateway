# start the app
app:
	go run service.payment/cmd/api.go
# test
test:
	- go clean -testcache
	- go test -race ./...
# code check
staticcheck:
	- staticcheck ./...
