# sysmon

> My system tray status monitor, goes well with DWM.

The defaults look like this:

`NET I/O 561 27O | CPU 15.67% | MEM 28.23% | SWAP 10.65% | 2022-03-04 18:23:14`

Emoji can also be used by pasting it directly on the format string, e.g.:

`$CPU% | $MEM% $SWAP% | $BRI% | $BAT%$CHAR | $TIME`


## Dependencies

- Go
- Make


## Installation

You can download and install the latest version directly using go with:

```
go install github.com/blmayer/sysmon@latest
```

Make sure your `GOPATH` is correctly set. Another option is to run

```
make install
```

The default installation directory is `~/local/bin` and it can be changed
setting the `PREFIX` variable, e.g.: `PREFIX=~/.bin make install`.


## Usage

Add `sysmon &` to your *.xinitrc* file, default configuration
is safe. The see all options that can be passed run `sysmon -h`.

If you are not using other status monitors you can try your configurations
running `sysmon` in a terminal.


### Changing the format

To change the output of sysmon use the `-F` flag, and pass the format string
as argument, use single quotes to prevent your shell messing up the string.
For example, just memory, CPU and date:

`sysmon -F 'MEM $MEM% | CPU $CPU% | $TIME'`


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
- [x] Emoji
- [ ] Configuration
  - [ ] Use yaml file
  - [X] Set update interval
  - [X] Define line format
  - [X] Select what components to use
- [ ] Wayland support (maybe another project)


## Known users

- Erik Dubois from [Arco Linux](https://arcolinux.info): [video](https://youtu.be/8SeCAXymXgw)


## Meta

License: MIT


### See also

- [gods](https://github.com/schachmat/gods)
- [gocaudices](https://github.com/lordrusk/gocaudices)

For a complete list visit [dwm's page](https://dwm.suckless.org/status_monitor/).
