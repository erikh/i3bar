package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-sysinfo"
	"github.com/erikh/i3bar"
	"github.com/leekchan/timeutil"
	"github.com/vishvananda/netlink"
)

var lastRxBytes, lastTxBytes uint64

func netActivity(dev string) (uint64, uint64) {
	link, err := netlink.LinkByName(dev)
	if err != nil {
		return 0, 0
	}

	rx := link.Attrs().Statistics.RxBytes
	tx := link.Attrs().Statistics.TxBytes
	defer func() {
		lastRxBytes = rx
		lastTxBytes = tx
	}()

	return rx - lastRxBytes, tx - lastTxBytes
}

// available and total respectively.
func memoryUsage() (uint64, uint64) {
	host, err := sysinfo.Host()
	if err != nil {
		return 0, 0
	}

	memory, err := host.Memory()
	if err != nil {
		return 0, 0
	}

	return memory.Available, memory.Total
}

func loadAverage() float64 {
	content, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return 0.0
	}
	first, err := strconv.ParseFloat(strings.Split(string(content), " ")[0], 64)
	if err != nil {
		return 0.0
	}

	return first
}

var lastTime time.Duration

func cpuUsage() float64 {
	host, err := sysinfo.Host()
	if err != nil {
		return 0
	}

	times, err := host.CPUTime()
	if err != nil {
		return 0
	}

	defer func() { lastTime = times.Total() - times.Idle }()

	return float64(times.Total()-times.Idle-lastTime) / float64(time.Duration(runtime.NumCPU())*time.Second) * 100
}

func formatNowTime(fmt string) string {
	now := time.Now()
	return timeutil.Strftime(&now, fmt)
}

func spotifyTrack() string {
	out, err := exec.Command("playerctl", "metadata", "-p", "spotify", "-f", "{{ xesam:artist }} - {{ xesam:title }}").Output()
	if err != nil {
		return err.Error()
	}

	return strings.TrimSpace(string(out))
}

func makeBlock(text string) *i3bar.Block {
	return &i3bar.Block{FullText: text, Color: "#888888", Separator: true}
}

func main() {
	ch := make(chan i3bar.StatusLine)
	go func() {
		for {
			avail, total := memoryUsage()
			rx, tx := netActivity("wlp7s0")
			ch <- i3bar.StatusLine{
				makeBlock(fmt.Sprintf("Net: %3.2fK Rx, %3.2fK Tx", float64(rx)/1024, float64(tx)/1024)),
				makeBlock(fmt.Sprintf("Memory: %3.2f In-Use, %3.2f Total", float64(total-avail)/1024/1024/1024, float64(total)/1024/1024/1024)),
				makeBlock(fmt.Sprintf("Load: %3.2f", loadAverage())),
				makeBlock(fmt.Sprintf("CPU: %3.2f", cpuUsage())),
				makeBlock(spotifyTrack()),
				makeBlock(formatNowTime("%Y-%m-%d %H:%M")),
			}
			time.Sleep(time.Second)
		}
	}()
	if err := i3bar.Encode(os.Stdout, &i3bar.Header{Version: 1}, ch); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
