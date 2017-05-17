package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jamiealquiza/cidrxpndr"
)

// Settings holds oftech settings.
var Settings struct {
	c       int
	t       int
	net     string
	port    string
	filter  string
	filterk string
	filterr *regexp.Regexp
}

func init() {
	flag.IntVar(&Settings.c, "c", 256, "request concurrency")
	flag.IntVar(&Settings.t, "t", 50, "request timeout in ms")
	flag.StringVar(&Settings.net, "net", "192.168.1.100/32", "network CIDR range")
	flag.StringVar(&Settings.port, "port", "8080", "ostat listening port")
	flag.StringVar(&Settings.filter, "filter", "", "regex filter by key")
	flag.Parse()

	kr := strings.Split(Settings.filter, ":")
	if len(kr) == 2 {
		var err error
		Settings.filterk = kr[0]
		Settings.filterr, err = regexp.Compile(kr[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	ips, _ := cidrxpndr.Expand(Settings.net)

	// Scatter-gather stuff.
	var wg sync.WaitGroup
	wg.Add(Settings.c)

	nodes := make(chan string, Settings.c+256)
	response := make(chan []byte, len(ips))

	for i := 0; i < Settings.c; i++ {
		go requester(nodes, response, &wg)
	}

	for _, i := range ips {
		nodes <- i
	}

	close(nodes)
	wg.Wait()
	close(response)
	//

	if Settings.filterr != nil {
		// Fetch and filter could be merged into one op
		// and parallelized in the requesters.
		responses := fetch(response)
		metrics := filter(responses, Settings.filterk, Settings.filterr)

		metricsList := []Stat{}
		for _, v := range metrics {
			metricsList = append(metricsList, v)
		}

		results, _ := json.Marshal(metricsList)
		fmt.Println(string(results))
	} else {
		results := rawFetch(response)
		fmt.Println(results)
	}

}

// requester reads from a channel of nodes n and hits the ostat
// api on each node, returning the results on channel r.
func requester(n chan string, r chan []byte, wg *sync.WaitGroup) {
	for h := range n {
		c, err := net.DialTimeout("tcp",
			h+":"+Settings.port,
			time.Duration(time.Millisecond*time.Duration(Settings.t)))
		if err != nil {
			continue
		}

		fmt.Fprintf(c, "stats")
		resp, err := bufio.NewReader(c).ReadBytes(10)
		if err != nil {
			continue
		} else {
			r <- resp
		}
	}
	wg.Done()
}

// fetch unmarshals respones and populates a map of
// metrics keyed by hostname.
func fetch(c chan []byte) map[string]Stat {
	metrics := make(map[string]Stat)

	for i := range c {
		m := Stat{}
		json.Unmarshal(i[:len(i)-1], &m)
		for k := range m {
			metrics[k] = m
		}
	}

	return metrics
}

// filter returns a new map[string]Stat from filtering
// m by key k with regex r.
func filter(metrics map[string]Stat, key string, re *regexp.Regexp) map[string]Stat {
	f := make(map[string]Stat)

	// Storage needs special handling for the mount
	// reference.
	keyf := strings.Split(key, ".")
	if keyf[0] == "storage" && len(keyf) == 3 {
		for k, v := range metrics {
			// Does the mount path exist?
			if path, ok := v[k].Storage[keyf[1]]; ok {
				if re.Match([]byte(path.Type)) {
					f[k] = v
				}
			}
		}
		return f
	}

	// This should be cleaned up, but
	// written this way to avoid dynamic references
	// through reflection. Also some struct references need
	// special handling.
	switch key {
	case "hostname":
		for k, v := range metrics {
			if re.Match([]byte(k)) {
				f[k] = v
			}
		}
	case "general.uptime":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Uptime)}) {
				f[k] = v
			}
		}
	case "general.cpu.model":
		for k, v := range metrics {
			if re.Match([]byte(v[k].General.CPU.Model)) {
				f[k] = v
			}
		}
	case "general.cpu.cores":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.CPU.Cores)}) {
				f[k] = v
			}
		}
	case "general.load.short":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Load.Short)}) {
				f[k] = v
			}
		}
	case "general.load.mid":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Load.Mid)}) {
				f[k] = v
			}
		}
	case "general.load.long":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Load.Long)}) {
				f[k] = v
			}
		}
	case "general.mem.total":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Total)}) {
				f[k] = v
			}
		}
	case "general.mem.free":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Free)}) {
				f[k] = v
			}
		}
	case "general.mem.used":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Used)}) {
				f[k] = v
			}
		}
	case "general.mem.usedp":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Usedp)}) {
				f[k] = v
			}
		}
	case "general.mem.shared":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Shared)}) {
				f[k] = v
			}
		}
	case "general.mem.buffer":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Buffer)}) {
				f[k] = v
			}
		}
	case "general.mem.swaptotal":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Swaptotal)}) {
				f[k] = v
			}
		}
	case "general.mem.swapfree":
		for k, v := range metrics {
			if re.Match([]byte{byte(v[k].General.Mem.Swapfree)}) {
				f[k] = v
			}
		}
	}
	return f
}

// rawFetch returns all the metrics responses as the string form
// of an array of JSON objects, without any parsing or inspection whatsoever.
func rawFetch(c chan []byte) string {
	metrics := []byte{91}
	for i := range c {
		metrics = append(metrics, i[:len(i)-1]...)
		metrics = append(metrics, 44)
	}

	if len(metrics) > 1 {
		metrics[len(metrics)-1] = 93
	} else {
		metrics = []byte{91, 93}
	}

	return string(metrics)
}
