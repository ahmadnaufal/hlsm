GOOS	 ?= linux
GOARCH  = amd64
ODIR	 := _output

export GO111MODULE ?= on

all: test compile

start:
	go run main.go

bin:
	go build -o apply main.go

$(ODIR):
	@mkdir -p $(ODIR)

compile: $(ODIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(ODIR)/apply main.go

test:
	go test ./...
