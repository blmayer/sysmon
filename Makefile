.PHONY: install uninstall clean

PREFIX=${HOME}/.local

sysmon: main.go go.sum go.mod
	go build .

install: sysmon
	mkdir -p ${PREFIX}/bin
	install -m755 sysmon ${PREFIX}/bin/sysmon

uninstall:
	rm -f ${PREFIX}/bin/sysmon

clean:
	rm sysmon
