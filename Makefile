PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

.PHONY: build install uninstall test clean

build:
	go build -o tuido

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 tuido $(DESTDIR)$(BINDIR)/tuido

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/tuido

test:
	go test -v ./...

clean:
	rm -f tuido
