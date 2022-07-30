package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// func getWeather() string {
// 	resp, err := http.Get("https://wttr.in?format=3")
// 	if err != nil {
// 		println(err.Error())
// 		return ""
// 	}
// 
// }

func getCPU() (float32, float32) {
	statFile, err := os.Open("/proc/stat")
	if err != nil {
		println(err.Error())
		return -1, -1
	}
	defer statFile.Close()

	var usr, ni, sys, idl, io, irq, soft, steal, guest, gni float32
	var pre string
	fmt.Fscanf(
		statFile, 
		"%s %f %f %f %f %f %f %f %f %f %f\n",
		&pre, &usr, &ni, &sys, &idl, &io, &irq, &soft, &steal, &guest, &gni,
	)
	
	return idl+io, usr+ni+sys+irq+soft+steal+guest+gni
}

func getMem() float32 {
	memFile, err := os.Open("/proc/meminfo")
	if err != nil {
		println(err)
		return -1.0
	}
	defer memFile.Close()

	var total, free float32
	var pre string
	fmt.Fscanf(memFile, "%s %f %s\n", &pre, &total, &pre)
	fmt.Fscanf(memFile, "%s %f %s\n", &pre, &free, &pre)
	
	return 1.0-(free/total)
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
	return used/total
}

func main() {
	x, err := xgb.NewConn() // connect to X
	if err != nil {
		println("Cannot connect to X:", err.Error())
		return
	}
	defer x.Close()

	root := xproto.Setup(x).DefaultScreen(x).Root
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	interval := 1*time.Second
	format := "CPU %.2f%% | MEM %.2f%% | SWAP %.2f%% | %s"
	datef := time.RFC1123
	// Parse CLA


	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var idl, busy float32
	for {
		select {
		case t := <-ticker.C:
			mem := getMem()
			swap := getSwap()

			idl2, busy2 := getCPU()
			total := idl + busy
			total2 := idl2 + busy2
			cpu := (busy2-busy)/(total2-total)

			out := fmt.Sprintf(format, cpu*100, mem*100, swap*100, t.Format(datef))
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
			idl, busy = idl2, busy2
		case <-sigs:
			return
		}
	}
}

