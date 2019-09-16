# go-openrpc

.PHONY: deps
deps:
	GO111MODULE=off go get -u github.com/gobuffalo/packr/v2/packr2

.PHONY: install
install:
	packr2 install

.PHONY: clean
clean:
	packr2 clean