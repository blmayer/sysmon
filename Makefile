PREFIX=${HOME}/.local

build:
	go build .

install: build
	install -m400 sysmon ${PREFIX}/bin

