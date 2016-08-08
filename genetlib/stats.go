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

type diskstats struct {
	Major             int64
	Minor             int64
	Name              string
	Reads_ok          int64
	Reads_m           int64
	Sectors_read      int64
	Time_reading      time.Duration
	Writes_ok         int64
	Writes_m          int64
	Sectors_written   int64
	Time_writting     time.Duration
	IOs_current       int64
	Time_IOs          time.Duration
	Weighted_time_IOs time.Duration
}

func (p procstats) String() string {
	return fmt.Sprintf("Cpu: %v Ctxt: %v Btime: %v Processes: %v Procs_running: %v Procs_blocked: %v SoftIRQs: %v", p.Cpu, p.Ctxt, p.Btime, p.Processes, p.Procs_running, p.Procs_blocked, p.Softirq)

}

func (d diskstats) String() string {
	return fmt.Sprintf("Name: %v Major: %v Minor: %v reads completed successfully: %v reads merged: %v sectors read: %v time spent reading: %v writes completed: %v writes merged: %v sectors written: %v time spent writing: %v I/Os currently in progress: %v time spent doing I/Os: %v weighted time spent doing I/Os: %v", d.Name, d.Major, d.Minor, d.Reads_ok, d.Reads_m, d.Sectors_read, d.Time_reading, d.Writes_ok, d.Writes_m, d.Sectors_written, d.Time_writting, d.IOs_current, d.Time_IOs, d.Weighted_time_IOs)
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
func GetStatsDisk() ([]diskstats, error) {
	var result []diskstats

	lines, err := ReadLines("/proc/diskstats")
	if err != nil {
		return result, err
	}

	for _, line := range lines {
		// filter out loop devices and ram
		if !strings.HasPrefix(strings.Fields(line)[2], "ram") && !strings.HasPrefix(strings.Fields(line)[2], "loop") {
			var res diskstats
			res.Major = stringtoint64(strings.Fields(line)[0])
			res.Minor = stringtoint64(strings.Fields(line)[1])
			res.Name = strings.Fields(line)[2]
			res.Reads_ok = stringtoint64(strings.Fields(line)[3])
			res.Reads_m = stringtoint64(strings.Fields(line)[4])
			res.Sectors_read = stringtoint64(strings.Fields(line)[5])
			res.Time_reading = time.Duration(stringtoint64(strings.Fields(line)[6]) * 1000000)
			res.Writes_ok = stringtoint64(strings.Fields(line)[7])
			res.Writes_m = stringtoint64(strings.Fields(line)[8])
			res.Sectors_written = stringtoint64(strings.Fields(line)[9])
			res.Time_writting = time.Duration(stringtoint64(strings.Fields(line)[10]) * 1000000)
			res.IOs_current = stringtoint64(strings.Fields(line)[11])
			res.Time_IOs = time.Duration(stringtoint64(strings.Fields(line)[12]) * 1000000)
			res.Weighted_time_IOs = time.Duration(stringtoint64(strings.Fields(line)[13]) * 100000)
			result = append(result, res)

		}
	}
	return result, nil
}
