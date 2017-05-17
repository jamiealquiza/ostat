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
		data[m[0]]["usedp"] = uint64(float64(data[m[0]]["used"].(uint64)) / float64(data[m[0]]["total"].(uint64)) * 100)
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
