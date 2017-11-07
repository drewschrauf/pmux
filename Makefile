default:
	go build

lint:
	golint ./...

test:
	gocoverutil -coverprofile=cover.out test -v -covermode=count ./...
	go tool cover -html=cover.out -o coverage.html
