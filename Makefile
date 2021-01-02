build:
	go build
	go vet
	go test

clean:
	rm -f search
