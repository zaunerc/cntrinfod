package system

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/zaunerc/cntinsight/types"
	"os"
	"os/exec"
)

func FetchContainerHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func FetchFirstMac() string {

	var firstMac string

	nics, _ := net.Interfaces()
	for _, nic := range nics {
		if nic.Name != "lo" {
			firstMac = nic.HardwareAddr
			break
		}
	}

	return firstMac
}

func FetchFirstIp() string {

	var firstIp string

	nics, _ := net.Interfaces()
	for _, nic := range nics {
		if nic.Name != "lo" {
			if len(nic.Addrs) > 0 {
				firstIp = nic.Addrs[0].Addr
			}
			break
		}
	}

	return firstIp
}

func FetchTcp46SocketInfo() []types.TcpSocketInfo {

	var socketInfo []types.TcpSocketInfo

	socketInfo = append(socketInfo, fetchTcpSocketInfo("tcp4")...)
	socketInfo = append(socketInfo, fetchTcpSocketInfo("tcp6")...)

	return socketInfo
}

/*
 * See https://github.com/shirou/gopsutil/blob/f20771d/net/net_linux.go#L258
 * for connection kinds.
 */
func fetchTcpSocketInfo(kind string) []types.TcpSocketInfo {

	var socketInfo []types.TcpSocketInfo
	connections, _ := net.Connections(kind)

	for _, con := range connections {

		// We are only interested in listening TCP sockets.
		if con.Status != "LISTEN" {
			continue
		}

		process, _ := process.NewProcess(con.Pid)

		user, _ := process.Username()
		programName, _ := process.Name()

		info := types.TcpSocketInfo{Protocol: kind, LocalIP: con.Laddr.IP,
			LocalPort: convertPortToStr(con.Laddr.Port), RemoteIP: con.Raddr.IP,
			RemotePort: convertPortToStr(con.Raddr.Port), State: con.Status, User: user,
			Pid: con.Pid, ProgramName: programName}

		socketInfo = append(socketInfo, info)
	}

	return socketInfo
}

func FetchUdp46SocketInfo() []types.UdpSocketInfo {

	var socketInfo []types.UdpSocketInfo

	socketInfo = append(socketInfo, fetchUdpSocketInfo("udp4")...)
	socketInfo = append(socketInfo, fetchUdpSocketInfo("udp6")...)

	return socketInfo
}

/*
 * See https://github.com/shirou/gopsutil/blob/f20771d/net/net_linux.go#L258
 * for connection kinds.
 */
func fetchUdpSocketInfo(kind string) []types.UdpSocketInfo {

	var socketInfo []types.UdpSocketInfo
	connections, _ := net.Connections(kind)

	for _, con := range connections {

		process, _ := process.NewProcess(con.Pid)

		user, _ := process.Username()
		programName, _ := process.Name()

		info := types.UdpSocketInfo{Protocol: kind, LocalIP: con.Laddr.IP,
			LocalPort: convertPortToStr(con.Laddr.Port), RemoteIP: con.Raddr.IP,
			RemotePort: convertPortToStr(con.Raddr.Port), User: user,
			Pid: con.Pid, ProgramName: programName}

		socketInfo = append(socketInfo, info)
	}

	return socketInfo
}

func convertPortToStr(port uint32) string {
	var localPort string
	if port == 0 {
		localPort = "*"
	} else {
		localPort = fmt.Sprint(port)
	}

	return localPort
}

func FetchProcessInfo() []types.ProcessInfo {

	var processInfo []types.ProcessInfo
	pids, _ := process.Pids()

	for _, pid := range pids {

		p, _ := process.NewProcess(pid)

		user, _ := p.Username()
		tty, _ := p.Terminal()
		name, _ := p.Name()
		exe, _ := p.Exe()
		cwd, _ := p.Cwd()
		cmd, _ := p.Cmdline()

		info := types.ProcessInfo{Pid: pid, User: user,
			Tty: tty, Name: name,
			Exe: exe, Cwd: cwd, Cmd: cmd}

		processInfo = append(processInfo, info)
	}

	return processInfo
}

func FetchProcessTree() string {

	out, _ := exec.Command("pstree", "-p").Output()
	return string(out)
}
