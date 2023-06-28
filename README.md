# sysmon

> My system tray status monitor, goes well with DWM.

The defaults look like this:

`NET I/O 561 27O | CPU 15.67% | MEM 28.23% | SWAP 10.65% | Thu, 04 Mar 2022 18:23:14 -03`


## Installation

To have it installed in your GOPATH run `go install github.com/blmayer/sysmon@latest`, or
run `make install`, the default installation directory is `~/local/bin` and it can be changed
setting the `PREFIX` variable, e.g.: `PREFIX=~/.bin make install`.


## Usage

Add `sysmon &` to your *.xinitrc* file, default configuration
is safe. The see all options that can be passed run `sysmon -h`.


### Defaults

The format is `NET I/O $NIN $NOUT | CPU $CPU% | MEM $MEM% | SWAP $SWAP% | $TIME`
and it gives you that example above.

Update interval is different for each component and can be changed by using
command line arguments, intervals are:

- Network in and out, in Kbps: 2 seconds
- CPU usage percentage: 2 seconds
- RAM usage percentage: 2 seconds
- SWAP usage in percent: 3 seconds
- Time: each second

By default brightness, battery and weather are not displayed, to enable them
use command line arguments. In order to display them you must pass the battery
or the display name.


## Roadmap

- [x] Clock
- [x] CPU Usage
- [x] MEM %
- [x] SWAP %
- [x] Network
- [x] Battery
- [x] Brightness
- [x] Weather (uses wttr.in)
- [ ] More than one swap file/partition
- [ ] Emoji?
- [ ] Configuration
  - [ ] Use yaml file
  - [X] Set update interval
  - [X] Define line format
  - [X] Select what components to use
- [ ] Wayland support (maybe another project)


## Meta

License: MIT


### See also

- [gods](https://github.com/schachmat/gods)
- [gocaudices](https://github.com/lordrusk/gocaudices)

For a complete list visit [dwm's page](https://dwm.suckless.org/status_monitor/).

