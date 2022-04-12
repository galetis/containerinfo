package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var port int
var tlsPort int
var hello bool

var start time.Time

//go:embed cert.pem
var cert []byte

//go:embed cert-key.pem
var key []byte

var hostname string

func init() {
	flag.IntVar(&port, "port", lookupEnvOrInt("PORT", 80), "port")
	flag.IntVar(&tlsPort, "tlsport", lookupEnvOrInt("TLS_PORT", 443), "tlsport")
	flag.BoolVar(&hello, "hello", false, "hello")
	flag.Parse()

	hostname, _ = os.Hostname()
}

func handler(ctx *fasthttp.RequestCtx) {

	w := ctx.Response.BodyWriter()

	if hello {

		//fmt.Fprintln(w, hostname)
		ctx.WriteString(hostname)

	} else {

		cinfo, _ := cpu.Info()
		vm, _ := mem.VirtualMemory()

		host, _ := host.Info()

		fmt.Fprintln(w, "Environ \t "+strings.Join(os.Environ(), "-"))
		fmt.Fprintln(w, "Hostname \t "+hostname)

		fmt.Fprintln(w, "Headers")

		ctx.Request.Header.VisitAll(func(key, value []byte) {
			fmt.Fprintln(w, "\t \t", string(key), string(value))
		})

		fmt.Fprintln(w, "Request.Addr \t "+ctx.RemoteAddr().String())
		fmt.Fprintln(w, "RequestURI \t "+ctx.URI().String())

		fmt.Fprintln(w, "Uptime \t\t "+time.Now().Sub(start).String())

		fmt.Fprintf(w, "Memory \t\t total %s available %s used %s free %s percent %%%f \t\t \n",
			byteCountSI(vm.Total),
			byteCountSI(vm.Available),
			byteCountSI(vm.Used),
			byteCountSI(vm.Free),
			vm.UsedPercent)

		fmt.Fprintf(w, "Swap \t\t total %s free %s cached %s\n",
			byteCountSI(vm.SwapTotal),
			byteCountSI(vm.SwapFree),
			byteCountSI(vm.SwapCached))

		for _, cpu := range cinfo {
			fmt.Fprintf(w, "Cpu \t\t model %s core %d mhz %f\n",
				cpu.ModelName,
				cpu.Cores,
				cpu.Mhz)
		}

		fmt.Fprintf(w, "Host \t\t id %s name %s os %s vrole %s vsystem %s\n",
			host.HostID,
			host.Hostname,
			host.OS,
			host.VirtualizationRole,
			host.VirtualizationSystem)
	}

}

func main() {

	start = time.Now()

	i, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		port = i
	}

	log.Printf("Listening on: %d", port)
	log.Printf("Listening tls on: %d", tlsPort)

	go func() {
		log.Fatal(fasthttp.ListenAndServeTLSEmbed(fmt.Sprintf(":%d", tlsPort), cert, key, handler))
	}()

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func byteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
