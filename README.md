# sysmon

> My system tray status monitor, goes well with DWM.

The defaults look like this:

`NET 561 IN 27O OUT | CPU 15.67% | MEM 28.23% | SWAP 10.65% | Thu, 04 Mar 2022 18:23:14 -03`


## Installation

Run `make`, the default installation directory is `~/local/bin` and it can be changed
setting the `PREFIX` variable, e.g.: `PREFIX=~/.bin make`.


## Usage

Add `sysmon &` to your *.xinitrc* file, default configuration
is safe.


### Defaults

Update interval is 1 second, components are:

- Network in and out, in Kbps
- CPU usage percentage
- RAM usage percentage
- SWAP usage in percent
- Date


## Roadmap

- [x] Clock
- [x] CPU Usage
- [x] MEM
- [x] SWAP
- [x] Network
- [ ] More than one swap file/partition
- [ ] Configuration
  - [ ] Use yaml file
  - [ ] Set update interval
  - [ ] Define line format
  - [ ] Select what components to use
- [ ] Wayland support (maybe another project)


## Meta

License: MIT


### See also

- [gods](https://github.com/schachmat/gods)
- [gocaudices](https://github.com/lordrusk/gocaudices)

For a complete list visit [dwm's page](https://dwm.suckless.org/status_monitor/).
