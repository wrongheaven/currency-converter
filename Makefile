install:
	@go build -o bin/cconv
	@cp bin/cconv $(GOPATH)/bin/cconv

run: install
	@cconv