package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"sort"
	"strconv"
	"os"
	"runtime"
)

func main() {
	port := 80
		
	i, err := strconv.Atoi(os.Getenv("PORT"))
        if err == nil {
		port = i
        }
	
	log.Printf("Listening on: %d", port)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		v, _ := mem.VirtualMemory()
		hostinfo, _ := host.Info()

		minfo := fmt.Sprintf("Memory: Total: %v, Free:%v, UsedPercent:%f%%", ByteCountDecimal(v.Total), ByteCountDecimal(v.Free), v.UsedPercent)

		info := map[string]string{

			"Go.Version":     runtime.Version(),
			"Host.Id":        hostinfo.HostID,
			"Host.Os":        hostinfo.OS,
			"Host.Hostname":  hostinfo.Hostname,
			"Host.Uptime":    strconv.FormatUint(hostinfo.Uptime, 10),
			"Host.Procs":     strconv.FormatUint(hostinfo.Procs, 10),
			"Host.Platform":  fmt.Sprintf("%s %s", hostinfo.Platform, hostinfo.PlatformVersion),
			"Memory":         minfo,
			"Request.Host":   r.Host,
			"Request.Addr":   r.RemoteAddr,
		}

		var keys []string
		for k := range info {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			_, _ = fmt.Fprintln(w, k, " \t ", info[k])
		}

	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func ByteCountDecimal(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
