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
	"log"
	"os"
	"regexp"
	"strings"
	"syscall"
)

var (
	blkdevr *regexp.Regexp
)

func init() {
	metrics.registerInput("storage", storageMetrics)
	// Find block devices with a '/' in the name. Assumes they're interesting.
	blkdevr, _ = regexp.Compile("/")
}

func storageMetrics() interface{} {
	var data = make(map[string]map[string]interface{})

	// Gets ["mount", "fstype"] pairs.
	fsMounts := getFsMounts()

	// For each mount, call statfs.
	for _, m := range fsMounts {
		data[m[0]] = make(map[string]interface{})

		s := &syscall.Statfs_t{}
		if err := syscall.Statfs(m[0], s); err != nil {
			log.Println(err)
		}

		data[m[0]]["type"] = m[1]
		// Needs some method to define units. KB for now.
		data[m[0]]["total"] = uint64(s.Bsize) * s.Blocks / 1024
		data[m[0]]["free"] = uint64(s.Bsize) * s.Bfree / 1024
		data[m[0]]["used"] = data[m[0]]["total"].(uint64) - data[m[0]]["free"].(uint64)
		data[m[0]]["usedp"] = uint64(float64(data[m[0]]["used"].(uint64)) / float64(data[m[0]]["free"].(uint64)) * 100)
		data[m[0]]["inodestotal"] = s.Files
		data[m[0]]["inodesfree"] = s.Ffree
		data[m[0]]["inodesused"] = data[m[0]]["inodestotal"].(uint64) - data[m[0]]["inodesfree"].(uint64)
	}

	return data
}

// getFsMounts finds all entries in /proc/mounts, were the device name
// includes a '/'. This assumes it's a real, or 'interesting' device that we care about.
// Returns the device's mount and fstype as a [2]{"mount", "fstype"} pair.
func getFsMounts() [][2]string {
	fsMounts := [][2]string{}

	f, _ := os.Open("/proc/mounts")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := strings.Fields(scanner.Text())
		if match := blkdevr.Match([]byte(l[0])); match {
			mount := [2]string{}
			mount[0], mount[1] = l[1], l[2]
			fsMounts = append(fsMounts, mount)
		}
	}

	return fsMounts
}
