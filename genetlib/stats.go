// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type cpustats struct {
	User       int64
	Nice       int64
	System     int64
	Idle       int64
	Iowait     int64
	Irq        int64
	Softirq    int64
	Steal      int64
	Guest      int64
	Guest_nice int64
}

type procstats struct {
	Cpu           cpustats
	Intr          string
	Ctxt          int64
	Btime         time.Time
	Processes     int64
	Procs_running int64
	Procs_blocked int64
	Softirq       string
}

func (p procstats) String() string {
	return fmt.Sprintf("Cpu: %v Ctxt: %v Btime: %v Processes: %v Procs_running: %v Procs_blocked: %v SoftIRQs: %v", p.Cpu, p.Ctxt, p.Btime, p.Processes, p.Procs_running, p.Procs_blocked, p.Softirq)

}

func stringtoint64(s string) int64 {
	result, _ := strconv.ParseInt(s, 10, 64)
	return result

}
func GetUptime() (time.Duration, error) {
	lines, err := ReadLines("/proc/uptime")
	if err != nil {
		return -1 * time.Second, err
	}
	uptime := strings.Fields(lines[0])[0]
	return time.ParseDuration(uptime + "s")
}

func GetStatsSys() (procstats, error) {
	var result procstats

	lines, err := ReadLines("/proc/stat")
	if err != nil {
		return result, err
	}
	for _, line := range lines {
		switch strings.Fields(line)[0] {
		case "cpu":
			result.Cpu = fillcpustats(line)
		case "intr":
			result.Intr = line[5:]
		case "ctxt":
			result.Ctxt = stringtoint64(line[5:])
		case "btime":
			result.Btime = time.Unix(stringtoint64(line[6:]), 0)
		case "processes":
			result.Processes = stringtoint64(line[10:])
		case "procs_running":
			result.Procs_running = stringtoint64(line[14:])
		case "procs_blocked":
			result.Procs_blocked = stringtoint64(line[14:])
		case "softirq":
			result.Softirq = line[8:]
		default:
			continue
		}
	}
	return result, err

}

func fillcpustats(str string) cpustats {
	var result cpustats
	result.User = stringtoint64(strings.Fields(str)[1])
	result.Nice = stringtoint64(strings.Fields(str)[2])
	result.System = stringtoint64(strings.Fields(str)[3])
	result.Idle = stringtoint64(strings.Fields(str)[4])
	result.Iowait = stringtoint64(strings.Fields(str)[5])
	result.Irq = stringtoint64(strings.Fields(str)[6])
	result.Softirq = stringtoint64(strings.Fields(str)[7])
	result.Steal = stringtoint64(strings.Fields(str)[8])
	result.Guest = stringtoint64(strings.Fields(str)[9])
	result.Guest_nice = stringtoint64(strings.Fields(str)[10])

	return result
}

func GetStatsNet() (string, error) {
	return "Not implemented", nil
}
func GetStatsAll() (string, error) {
	return "Not implemented", nil
}
func GetStatsDisk() (string, error) {
	return "Not implemented", nil
}
