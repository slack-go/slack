EXAMPLES := $(shell find examples/ -maxdepth 1 -type d -exec sh -c 'echo $(basename {})' \;)
EXLIST := $(subst examples/,,$(EXAMPLES))

ifeq ($(TRAVIS_BUILD_DIR),)
	GOPATH := $(GOPATH)
else
	GOPATH := $(GOPATH):$(TRAVIS_BUILD_DIR)
endif

all: chmod clean lint test coverage $(EXLIST)

# this chmod is a side effect of some windows-development related issues
chmod:
	chmod +x script/*

lint:
	@script/lint

test:
	@script/test

coverage:
	@script/coverage

$(EXLIST):
	@echo $@
	@go test -v ./examples/$@
	@gocov test ./examples/$@ | gocov report

clean:
	@rm -rf bin/ pkg/

.PHONY: all chmod clean lint test coverage $(EXLIST)
