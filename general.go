// The MIT License (MIT)
//
// Copyright (c) 2016 Jamie Alquiza
//
// http://knowyourmeme.com/memes/deal-with-it.
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
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var (
	cpuModelr *regexp.Regexp
	cpuModel  string
)

func init() {
	metrics.registerInput("general", generalMetrics)
	// This will never change at runtime. Fetch CPU info once.
	cpuModelr, _ = regexp.Compile("model name")
	cpuModel = getCpuModel()
}

func generalMetrics() interface{} {
	var data struct {
		Uptime int64 `json:"uptime"`
		CPU    struct {
			Model string `json:"model"`
			Cores int    `json:"cores"`
		} `json:"cpu"`
		Load struct {
			Short float64 `json:"short"`
			Mid   float64 `json:"mid"`
			Long  float64 `json:"long"`
			Procs uint16  `json:"procs"`
		} `json:"load"`
		Mem struct {
			Total     uint64 `json:"total"`
			Free      uint64 `json:"free"`
			Used      uint64 `json:"used"`
			Usedp     uint64 `json:"usedp"`
			Shared    uint64 `json:"shared"`
			Buffer    uint64 `json:"buffer"`
			Swaptotal uint64 `json:"swaptotal"`
			Swapfree  uint64 `json:"swapfree"`
		} `json:"mem"`
	}

	// Fetch sysinfo and load.
	s := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(s); err != nil {
		log.Println(err)
	}

	data.Uptime = s.Uptime
	data.CPU.Model = cpuModel
	data.CPU.Cores = runtime.NumCPU()
	data.Load.Short = siLoadShift(s.Loads[0])
	data.Load.Mid = siLoadShift(s.Loads[1])
	data.Load.Long = siLoadShift(s.Loads[2])
	data.Load.Procs = s.Procs
	// Probably want some way to define what units we want. KB for now.
	data.Mem.Total = s.Totalram / 1024
	data.Mem.Free = s.Freeram / 1024
	data.Mem.Used = data.Mem.Total - data.Mem.Free
	data.Mem.Usedp = uint64(float64(data.Mem.Used) / float64(data.Mem.Total) * 100.00)
	data.Mem.Shared = s.Sharedram / 1024
	data.Mem.Buffer = s.Bufferram / 1024
	data.Mem.Swaptotal = s.Totalswap / 1024
	data.Mem.Swapfree = s.Freeswap / 1024

	return data
}

// siLoadShift takes uint64's from sysinfo and
// translates / formats into 2 decimal place floats.
// This needs to be made more efficient.
func siLoadShift(u uint64) float64 {
	n := fmt.Sprintf("%.2f", float64(u)/65536.0)
	f, _ := strconv.ParseFloat(n, 64)
	return f
}

// getCpuModel returns the CPU model string from /proc/cpuinfo.
func getCpuModel() string {
	f, _ := os.Open("/proc/cpuinfo")
	defer f.Close()

	var cpuInfo string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		if cpuModelr.MatchString(l) {
			cpuInfo = string(l)
			break
		}
	}

	model := strings.TrimSpace(strings.Split(cpuInfo, ":")[1])
	return model
}
