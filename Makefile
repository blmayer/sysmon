PREFIX=${HOME}/.local

build:
	go build .

install: build
	mkdir -p ${PREFIX}/bin
	install -m755 sysmon ${PREFIX}/bin/sysmon

uninstall:
	rm -f ${PREFIX}/bin/sysmon

