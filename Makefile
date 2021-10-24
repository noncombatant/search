build:
	go build
	go vet
	go test
	./test.sh

clean:
	-rm -f search
	-rm -f test_files
