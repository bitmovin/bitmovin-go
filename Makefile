COVERAGEDIR = coverage

all: test cover
test:
	if [ ! -d coverage ]; then mkdir coverage; fi
	go test -v ./services -race -cover -coverprofile=$(COVERAGEDIR)/services.coverprofile
cover:
	go tool cover -html=$(COVERAGEDIR)/services.coverprofile -o $(COVERAGEDIR)/services.html
