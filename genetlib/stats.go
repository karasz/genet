// Copyright © 2016 Nagy Károly Gábriel <karasz@jpi.io>
// This file, part of genet, is free and unencumbered software
// released into the public domain.
// For more information, please refer to <http://unlicense.org/>
//

package genetlib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var tcpState = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}

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

type netstats struct {
	Prot       string
	User       string
	Name       string
	Pid        string
	State      string
	LocalIp    string
	LocalPort  string
	RemoteIp   string
	RemotePort string
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

type mntent struct {
	fsname  string
	dir     string
	mnttype string
	opts    string
	freq    int64
	passno  int64
}

type strucdf struct {
	Name  string
	Total int64
	Free  int64
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

func hextodec(s string) string {
	var result string

	d, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		fmt.Println(err)
		result = ""
	}

	result = strconv.FormatInt(d, 10)
	return result
}

func getIPfromHex(ipport string) string {
	hip := strings.Split(ipport, ":")[0]
	port := hextodec(strings.Split(ipport, ":")[1])
	ip := fmt.Sprintf("%v.%v.%v.%v", hextodec(hip[6:8]), hextodec(hip[4:6]), hextodec(hip[2:4]), hextodec(hip[0:2]))
	return fmt.Sprintf("%s:%s", ip, port)
}
func getNamefromPid(pid string) string {
	exe := fmt.Sprintf("/proc/%s/exe", pid)
	path, _ := os.Readlink(exe)
	n := strings.Split(path, "/")
	name := n[len(n)-1]
	return strings.Title(name)
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

func findPid(inode string) string {
	pid := "-"

	d, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		fmt.Println(err)
		return pid
	}

	re := regexp.MustCompile(inode)
	for _, item := range d {
		path, _ := os.Readlink(item)
		out := re.FindString(path)
		if len(out) != 0 {
			pid = strings.Split(item, "/")[2]
		}
	}
	return pid
}

func GetStatsNet() ([]netstats, error) {
	var result []netstats

	lines, err := ReadLinesNoFrills("/proc/net/tcp", 1, "")
	if err != nil {
		return result, err
	}

	for _, line := range lines {
		res := tcpudpFillLine(line, "tcp")
		result = append(result, res)
	}

	lines, err = ReadLinesNoFrills("/proc/net/udp", 1, "")
	if err != nil {
		return result, err
	}

	for _, line := range lines {
		res := tcpudpFillLine(line, "udp")
		result = append(result, res)
	}

	return result, nil
}

func tcpudpFillLine(line string, prot string) netstats {
	var res netstats
	us, _ := strconv.Atoi(strings.Fields(line)[7])
	res.Prot = prot
	res.User, _ = LookupLinuxUserById(us)
	res.Pid = findPid(strings.Fields(line)[9])
	res.Name = getNamefromPid(res.Pid)
	res.State = tcpState[strings.Fields(line)[3]]
	res.LocalIp = strings.Split(getIPfromHex(strings.Fields(line)[1]), ":")[0]
	res.LocalPort = strings.Split(getIPfromHex(strings.Fields(line)[1]), ":")[1]
	res.RemoteIp = strings.Split(getIPfromHex(strings.Fields(line)[2]), ":")[0]
	res.RemotePort = strings.Split(getIPfromHex(strings.Fields(line)[2]), ":")[1]
	return res
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

func GetMountSpace(path string) (int64, int64, error) {

	struc := syscall.Statfs_t{}

	err := syscall.Statfs(path, &struc)
	if err != nil {
		return 0, 0, err
	}

	return int64(struc.Bsize) * int64(struc.Blocks), int64(struc.Bsize) * int64(struc.Bfree), nil
}

func parseMounts() ([]mntent, error) {
	result := []mntent{}

	lines, err := ReadLines("/proc/mounts")
	if err != nil {
		return result, err
	}
	for _, line := range lines {
		var lmnt mntent
		lmnt.fsname = strings.Fields(line)[0]
		lmnt.dir = strings.Fields(line)[1]
		lmnt.mnttype = strings.Fields(line)[2]
		lmnt.opts = strings.Fields(line)[3]
		lmnt.freq = stringtoint64(strings.Fields(line)[4])
		lmnt.passno = stringtoint64(strings.Fields(line)[5])
		result = append(result, lmnt)
	}
	return result, nil
}

func GetDF() []strucdf {
	var result []strucdf
	mounts, _ := parseMounts()
	for _, mount := range mounts {
		var st strucdf
		st.Name = mount.dir
		st.Total, st.Free, _ = GetMountSpace(mount.dir)
		//filter out devices with 0 blocks
		if st.Total != 0 {
			result = append(result, st)
		}
	}
	return result
}
