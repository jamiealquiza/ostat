package main

import (
        "bufio"
        "flag"
        "fmt"
        "net"
        "sync"
        "time"

        "github.com/jamiealquiza/cidrxpndr"
)

var Settings struct {
        c    int
        t    int
        net  string
        port string
}

func init() {
        flag.IntVar(&Settings.c, "c", 256, "request concurrency")
        flag.IntVar(&Settings.t, "t", 25, "request timeout in ms")
        flag.StringVar(&Settings.net, "net", "192.168.1.100/32", "network CIDR range")
        flag.StringVar(&Settings.port, "port", "8080", "ostat listening port")
        flag.Parse()
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

        output := []byte{91}
        for i := range response {
                output = append(output, i[:len(i)-1]...)
                output = append(output, 44)
        }
        output[len(output)-1] = 93

        fmt.Println(string(output))
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
