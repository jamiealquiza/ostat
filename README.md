# ostat

# Overview
Ostat is an extremely lightweight Linux daemon that provides basic systems metrics in json via a TCP endpoint.

Ostat includes the ofetch tool (utils/ofetch) for quickly gathering metrics from a whole network of hosts.

# Usage

```shell
$ ostat -h
Usage of ./ostat:
  -listen string
        Listen address:port (default "localhost:8080")
  -update-int int
        Metrics update interval (default 30)
```

Output with comments:

```shell
$ echo "stats" | nc localhost 8080 | jq '.'                                                                           
{
  "hostname01": {
    "general": {
      "uptime": 405602, # seconds
      "cpu": {
        "Model": "Intel(R) Xeon(R) CPU E5-2676 v3 @ 2.40GHz",
        "cores": 2
      },
      "load": {
        "short": 1.05,
        "mid": 0.01,
        "long": 0.05,
        "procs": 117
      },
      "mem": {
        "total": 8175632, # All mem values in KB
        "free": 6845488,
        "used": 1330144,
        "usedp": 16, # Memory used in percent
        "shared": 0, 
        "buffer": 102332,
        "swaptotal": 0,
        "swapfree": 0
      }
    },
    "storage": {
      "/": { # All mounted block-device filesystems are automatically discovered
        "free": 60062488, # All storage values in KB
        "inodesfree": 3838717,
        "inodestotal": 3932160,
        "inodesused": 93443,
        "total": 61784292,
        "type": "ext4",
        "used": 1721804,
        "usedp": 2 # Storage used in percent
      }
    }
  }
}
```

# Setup

- `go get github.com/jamiealquiza/ostat`
- `go install github.com/jamiealquiza/ostat`
- Binary will be found at `$GOPATH/bin/ostat`
