install:
	@go build -o cconv
	@mv cconv $(GOPATH)/bin/cconv

run: install
	@cconv