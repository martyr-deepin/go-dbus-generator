PREFIX = /usr

all: build

build: *.go
	go build -o dbus-generator

install: build
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp dbus-generator ${DESTDIR}${PREFIX}/bin

test-golang-max-match-rules: build
	cd testdata/issue_golang_max_match_rules/; rm -rf gen; ../../dbus-generator -target GoLang; go test
