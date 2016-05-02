// The MIT License (MIT)
//
// Copyright (c) 2016 Jamie Alquiza
//
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
        "bufio"
        "encoding/json"
        "flag"
        "fmt"
        "net"
        "os"
        "regexp"
        "sync"
        "time"

        "github.com/jamiealquiza/cidrxpndr"
)

var Settings struct {
        c       int
        t       int
        net     string
        port    string
        filter  string
        filterr *regexp.Regexp
}

func init() {
        flag.IntVar(&Settings.c, "c", 256, "request concurrency")
        flag.IntVar(&Settings.t, "t", 25, "request timeout in ms")
        flag.StringVar(&Settings.net, "net", "192.168.1.100/32", "network CIDR range")
        flag.StringVar(&Settings.port, "port", "8080", "ostat listening port")
        flag.StringVar(&Settings.filter, "filter", "", "regex filter by key")
        flag.Parse()

        if Settings.filter != "" {
                var err error
                Settings.filterr, err = regexp.Compile(Settings.filter)
                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }
        }
}

func main() {
        ips, _ := cidrxpndr.Expand(Settings.net)

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

        if Settings.filterr != nil {
                output := parseFetch(response)
                for _, v := range output {
                        fmt.Println(v)
                }
        } else {
                output := rawFetch(response)
                fmt.Println(string(output))
        }

}

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

// parseFetch unmarshals respones and populates a map of
// metrics keyed by hostname.
func parseFetch(c chan []byte) map[string]Stat {
        metrics := make(map[string]Stat)

        for i := range c {
                m := Stat{}
                json.Unmarshal(i[:len(i)-1], &m)
                for k, _ := range m {
                        metrics[k] = m
                }
        }

        return metrics
}

// rawFetch returns all the metrics responses
// in a []byte without any parsing or inspection whatsoever.
func rawFetch(c chan []byte) []byte {
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

        return metrics
}
