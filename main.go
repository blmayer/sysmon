package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

const help = `sysmon v1.1.1 A SYStem MONitor for your system bar, designed for DWM.

Usage:
  sysmon [OPTIONS]

Available options:
  -h
  --help		show this help
  -t
  --time		set time refresh interval
  -c
  --cpu			set cpu refresh interval
  -m
  --mem			set memmory refresh interval
  -s
  --swap		set swap refresh interval
  -n
  --net			set network refresh interval
  -B
  --brightness		set brightness refresh interval, needs -d
  -b
  --battery		set battery refresh interval, needs -N
  -N
  --battery-name	set battery name to get information
  -d
  --display-name	set display name to get information
  -w
  --weather		set weather refresh interval
  -f
  --format		define output format, each module is defined
  			using $CPU, $MEM, $SWAP, $BAT, $NET, $BRI, $TIME,
			$WTR or $CHAR. Default:
			'NET I/O $NIN $NOUT | CPU $CPU% | MEM $MEM% | SWAP $SWAP% | $TIME'
  -T
  --time-format		set time format, default "2006-01-02 15:04:05", any
  			format can be passed using go's time format string

Examples:
  sysmon
All default values
 
  sysmon -s 10
Uses 10 seconds of interval for swap

License:
MIT Copyright (c) 2022-24 Brian Mayer

Report bugs to: bleemayer@gmail.com
Or open an issue at https://github.com/blmayer/sysmon`

func getWeather() string {
	resp, err := http.Get("https://wttr.in?format=%x+%w+%h+%t")
	if err != nil {
		println(err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return ""
	}

	return string(body)
}

func getCPU() (int, int) {
	statFile, err := os.Open("/proc/stat")
	if err != nil {
		println(err.Error())
		return -1, -1
	}
	defer statFile.Close()

	var usr, ni, sys, idl, io, irq, soft, steal, guest, gni int
	var pre string
	fmt.Fscanf(
		statFile,
		"%s %d %d %d %d %d %d %d %d %d %d\n",
		&pre, &usr, &ni, &sys, &idl, &io, &irq, &soft, &steal, &guest, &gni,
	)

	return idl + io, usr + ni + sys + irq + soft + steal + guest + gni
}

func getMem() float32 {
	memFile, err := os.Open("/proc/meminfo")
	if err != nil {
		println(err)
		return -1.0
	}
	defer memFile.Close()

	var total, avail float32
	var pre string
	fmt.Fscanf(memFile, "%s %f %s\n", &pre, &total, &pre)
	fmt.Fscanf(memFile, "%s\n", &pre)
	fmt.Fscanf(memFile, "%s %f %s\n", &pre, &avail, &pre)

	return (total-avail) / total
}

func getBat() (int, bool) {
	batFile, err := os.Open("/sys/class/power_supply/" + batName + "/capacity")
	if err != nil {
		println(err)
		return -1.0, false
	}
	defer batFile.Close()

	var cap int
	fmt.Fscanf(batFile, "%d\n", &cap)

	chargeFile, err := os.Open("/sys/class/power_supply/" + batName + "/status")
	if err != nil {
		println(err)
		return -1.0, false
	}
	defer chargeFile.Close()
	var status string
	fmt.Fscanf(chargeFile, "%s\n", &status)

	return cap, status == "Charging"
}

func getBrightness() int {
	briFile, err := os.Open("/sys/class/backlight/" + displayName + "/brightness")
	if err != nil {
		println(err)
		return -1.0
	}
	defer briFile.Close()

	var bri float32
	fmt.Fscanf(briFile, "%f\n", &bri)

	maxFile, err := os.Open("/sys/class/backlight/" + displayName + "/max_brightness")
	if err != nil {
		println(err)
		return -1.0
	}
	defer maxFile.Close()
	var max float32
	fmt.Fscanf(maxFile, "%f\n", &max)

	return int(100 * bri / max)
}

func getSwap() float32 {
	memFile, err := os.Open("/proc/swaps")
	if err != nil {
		println(err)
		return -1.0
	}
	defer memFile.Close()

	var total, used float32
	var pre string
	fmt.Fscanf(memFile, "%s %s %s %s %s\n", &pre, &pre, &pre, &pre, &pre)
	fmt.Fscanf(memFile, "%s %s %f %f %s\n", &pre, &pre, &total, &used, &pre)
	return used / total
}

func getNet() (int, int) {
	netFile, err := os.Open("/proc/net/dev")
	if err != nil {
		println(err)
		return -1, -1
	}
	defer netFile.Close()

	var e string
	var in, out, rc, tx, z int

	// Skip first two lines (table header)
	fmt.Fscanln(netFile, &e, &e)
	fmt.Fscan(netFile,
		&e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e, &e)
	_, err = fmt.Fscanln(netFile,
		&e, &rc, &z, &z, &z, &z, &z, &z, &z, &tx, &z, &z, &z, &z, &z, &z, &z)
	for err == nil {
		in += rc
		out += tx
		_, err = fmt.Fscanln(netFile, &e,
			&rc, &z, &z, &z, &z, &z, &z, &z, &tx, &z, &z, &z, &z, &z, &z, &z)
	}

	return in, out
}

var (
	batName     string
	displayName string
)

type values struct {
	cpu      float32
	cpuidl   int
	cpubusy  int
	mem      float32
	bri      int
	swap     float32
	wtr      string
	bat      int
	charging string
	netin    int
	netout   int
	netind   int
	netoutd  int
}

func main() {
	// dafault values
	format := "NET I/O $NIN $NOUT | CPU $CPU% | MEM $MEM% | SWAP $SWAP% | $TIME"
	timeFormat := time.DateTime
	tickers := map[string]*time.Ticker{
		"time": time.NewTicker(time.Second),
		"cpu":  time.NewTicker(2 * time.Second),
		"net":  time.NewTicker(2 * time.Second),
		"mem":  time.NewTicker(2 * time.Second),
		"swap": time.NewTicker(3 * time.Second),
		"bat":  {},
		"bri":  {},
		"wtr":  {},
	}

	// overrides
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-t", "--time":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["time"].Stop()
			tickers["time"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-T", "--time-format":
			i++
			timeFormat = os.Args[i]
		case "-b", "--battery":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["bat"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-N", "battery-name":
			i++
			batName = os.Args[i]
		case "-c", "--cpu":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["cpu"].Stop()
			tickers["cpu"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-n", "--net":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["net"].Stop()
			tickers["net"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-m", "--mem":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["mem"].Stop()
			tickers["mem"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-s", "--swap":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["swap"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-B", "--brightness":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["bri"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-d", "--display-name":
			i++
			displayName = os.Args[i]
		case "-w", "--weather":
			i++
			val, err := strconv.Atoi(os.Args[i])
			if err != nil {
				println(err)
				os.Exit(-1)
			}
			tickers["wtr"] = time.NewTicker(time.Duration(val) * time.Second)
		case "-f", "--format":
			i++
			format = os.Args[i]
		case "-h", "--help":
			fmt.Println(help)
			os.Exit(0)
		default:
			println("unreckognized argument")
			os.Exit(-1)
		}
	}

	x, err := xgb.NewConn() // connect to X
	if err != nil {
		println("cannot connect to X:", err.Error())
		return
	}
	defer x.Close()

	root := xproto.Setup(x).DefaultScreen(x).Root
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// TODO: only refresh what is needed. i.e. is on the format string
	vals := values{}
	for {
		select {
		case t := <-tickers["time"].C:
			rep := strings.NewReplacer(
				"$CPU", fmt.Sprintf("%.2f", vals.cpu*100),
				"$MEM", fmt.Sprintf("%.2f", vals.mem*100),
				"$SWAP", fmt.Sprintf("%.2f", vals.swap*100),
				"$TIME", t.Format(timeFormat),
				"$BRI", fmt.Sprintf("%d", vals.bri),
				"$WTR", vals.wtr,
				"$BAT", fmt.Sprintf("%d", vals.bat),
				"$CHAR", vals.charging,
				"$NIN", fmt.Sprintf("%d", vals.netind/1024),
				"$NOUT", fmt.Sprintf("%d", vals.netoutd/1024),
			)
			out := rep.Replace(format)

			xproto.ChangeProperty(
				x,
				xproto.PropModeReplace,
				root,
				xproto.AtomWmName,
				xproto.AtomString,
				8,
				uint32(len(out)),
				[]byte(out),
			)
		case <-tickers["bat"].C:
			var charging bool
			vals.bat, charging = getBat()
			if charging {
				vals.charging = "Charging"
			}
		case <-tickers["cpu"].C:
			prevIdl, prevBusy := vals.cpuidl, vals.cpubusy
			prevTotal := prevIdl + prevBusy
			vals.cpuidl, vals.cpubusy = getCPU()
			total := vals.cpuidl + vals.cpubusy
			vals.cpu = float32(vals.cpubusy-prevBusy) / float32(total-prevTotal)
		case <-tickers["mem"].C:
			vals.mem = getMem()
		case <-tickers["bri"].C:
			vals.bri = getBrightness()
		case <-tickers["swap"].C:
			vals.swap = getSwap()
		case <-tickers["wtr"].C:
			vals.wtr = getWeather()
		case <-tickers["net"].C:
			vals.netind, vals.netoutd = -vals.netin, -vals.netout
			vals.netin, vals.netout = getNet()
			vals.netind, vals.netoutd = vals.netind+vals.netin, vals.netoutd+vals.netout
		case <-sigs:
			return
		}
	}
}
