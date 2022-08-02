# sysmon

> My system tray status monitor

Very minimalistic, but this works for me. If you want a real
system tray there are better options:

- [gods](https://github.com/schachmat/gods)
- [gocaudices](https://github.com/lordrusk/gocaudices)

For a complete list visit [dwm's page](https://dwm.suckless.org/status_monitor/).


## Installation

`go install`


## Usage

Add `sysmon &` to your *.xinitrc* file, default configuration
is safe.


### Defaults

Update interval is 1 second, the result is like this:

`CPU 15.67% | MEM 28.23% | SWAP 10.65% | Thu, 04 Mar 2022 18:23:14 -03`

