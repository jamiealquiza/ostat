# cidrxpndr
Go library - convert CIDR notation to list of IPs

# Example

```go
package main

import (
	"fmt"

	"github.com/jamiealquiza/cidrxpndr"
)

func main() {
        ips, _ := cidrxpndr.Expand("10.0.1.248/28")
        for _, i := range ips {
        	fmt.Println(i)
        }
}
```

```
$ go run example.go
10.0.1.249
10.0.1.250
10.0.1.251
10.0.1.252
10.0.1.253
10.0.1.254
10.0.2.1
10.0.2.2
10.0.2.3
10.0.2.4
10.0.2.5
10.0.2.6
10.0.2.7
10.0.2.8
```
