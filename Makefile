BINARY    := scribe
GO        := CGO_ENABLED=0 go

.PHONY: build install run test clean

build:
	$(GO) build -o $(BINARY) .

install: build
	sudo mv $(BINARY) /usr/local/bin/$(BINARY)

run: build
	./$(BINARY)

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY)
