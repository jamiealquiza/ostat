
# Overview
Ofetch is an extremely fast scatter/gather tool for fetching metrics from a network of nodes running ostat. Ofetch takes a CIDR notated network and parallelizes requests to the entire range.

# Usage

```bash
Usage of ./ofetch:
  -c int
        request concurrency (default 256)
  -filter string
        regex filter by key
  -net string
        network CIDR range (default "192.168.1.100/32")
  -port string
        ostat listening port (default "8080")
  -t int
        request timeout in ms (default 50)
```

The `-filter` directive takes a "key:regex" format parameter used for filtering metrics. Use "hostname:regex" to filter hostnames, or "general.cpu.model:regex" references for other keys.

```bash
$ ./ofetch -net="192.168.239.10/28" | jq '.'
[
  {
    "somenode01": {
      "general": {
        "uptime": 4313705,
        "cpu": {
          "Model": "Intel(R) Xeon(R) CPU E5-2670 v2 @ 2.50GHz",
          "cores": 1
        },
        "load": {
          "short": 0.01,
          "mid": 0.03,
          "long": 0.05,
          "procs": 135
        },
        "mem": {
          "total": 2050172,
          "free": 263028,
          "used": 1787144,
          "usedp": 87,
          "shared": 0,
          "buffer": 256208,
          "swaptotal": 0,
          "swapfree": 0
        }
      },
      "storage": {
        "/": {
          "free": 79435428,
          "inodesfree": 5079300,
          "inodestotal": 5242880,
          "inodesused": 163580,
          "total": 82559188,
          "type": "ext4",
          "used": 3123760,
          "usedp": 3
        }
      }
    }
  },
  {
    "somenode02": {
      "general": {
        "uptime": 5200439,
        "cpu": {
          "Model": "Intel(R) Xeon(R) CPU E5-2670 v2 @ 2.50GHz",
          "cores": 2
        },
        "load": {
          "short": 0,
          "mid": 0.03,
          "long": 0.05,
          "procs": 176
        },
        "mem": {
          "total": 7629404,
          "free": 3953864,
          "used": 3675540,
          "usedp": 48,
          "shared": 0,
          "buffer": 293628,
          "swaptotal": 0,
          "swapfree": 0
        }
      },
      "storage": {
        "/": {
          "free": 79461584,
          "inodesfree": 5079816,
          "inodestotal": 5242880,
          "inodesused": 163064,
          "total": 82569904,
          "type": "ext4",
          "used": 3108320,
          "usedp": 3
        },
        "/mnt": {
          "free": 30779832,
          "inodesfree": 1966069,
          "inodestotal": 1966080,
          "inodesused": 11,
          "total": 30956028,
          "type": "ext3",
          "used": 176196,
          "usedp": 0
        }
      }
    }
  },
  {
    "somenode03": {
      "general": {
        "uptime": 3119959,
        "cpu": {
          "Model": "Intel(R) Xeon(R) CPU E5-2680 v2 @ 2.80GHz",
          "cores": 8
        },
        "load": {
          "short": 1.02,
          "mid": 1.02,
          "long": 1.08,
          "procs": 165
        },
        "mem": {
          "total": 15339148,
          "free": 9875264,
          "used": 5463884,
          "usedp": 35,
          "shared": 0,
          "buffer": 287380,
          "swaptotal": 0,
          "swapfree": 0
        }
      },
      "storage": {
        "/": {
          "free": 79646532,
          "inodesfree": 5075398,
          "inodestotal": 5242880,
          "inodesused": 167482,
          "total": 82569904,
          "type": "ext4",
          "used": 2923372,
          "usedp": 3
        },
        "/mnt": {
          "free": 82370348,
          "inodesfree": 5242869,
          "inodestotal": 5242880,
          "inodesused": 11,
          "total": 82558640,
          "type": "ext3",
          "used": 188292,
          "usedp": 0
        }
      }
    }
  },
  {
    "somenode04": {
      "general": {
        "uptime": 5200403,
        "cpu": {
          "Model": "Intel(R) Xeon(R) CPU E5-2676 v3 @ 2.40GHz",
          "cores": 1
        },
        "load": {
          "short": 0,
          "mid": 0.01,
          "long": 0.05,
          "procs": 121
        },
        "mem": {
          "total": 1017984,
          "free": 286312,
          "used": 731672,
          "usedp": 71,
          "shared": 0,
          "buffer": 167728,
          "swaptotal": 0,
          "swapfree": 0
        }
      },
      "storage": {
        "/": {
          "free": 80557684,
          "inodesfree": 5082708,
          "inodestotal": 5242880,
          "inodesused": 160172,
          "total": 82559188,
          "type": "ext4",
          "used": 2001504,
          "usedp": 2
        }
      }
    }
  }
]
```

# Setup
- ```go get github.com/jamiealquiza/ostat```
- ```go install go install github.com/jamiealquiza/ostat/utils/ofetch```
- binary will be found at ```$GOPATH/bin/ofetch```
