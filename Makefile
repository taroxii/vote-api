run:
	go run main.go

build:
	go build -o app main.go


unit-test:
	go test -v -count=1 ./src/...

integration-test:
	INTEGRATION=true go test -v -count=1 -p=1 ./src/...

coverage-test:
	go test -coverprofile cover.out ./src/...

coverage-test-html: coverage-test
	go tool cover -html=cover.out


clean-test:
	go clean -testcache