PREFIX = /usr

all: build

build: *.go
	go build -o dbus-generator

install: build
	install -Dm755 dbus-generator ${DESTDIR}${PREFIX}/bin/dbus-generator

test-golang-max-match-rules: build
	cd testdata/issue_golang_max_match_rules/; rm -rf gen; ../../dbus-generator -target GoLang; go test
