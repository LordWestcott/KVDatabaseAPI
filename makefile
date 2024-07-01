run:
	go run .

test:
	go test -race -coverprofile=cover.out ./...

bench:
	go test -race -bench=. ./...

coverage: test
	go tool cover -html=cover.out
