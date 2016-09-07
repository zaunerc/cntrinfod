package types

type TcpSocketInfo struct {
	Protocol    string
	LocalIP     string
	LocalPort   string
	RemoteIP    string
	RemotePort  string
	State       string
	User        string
	Pid         int32
	ProgramName string
}

type UdpSocketInfo struct {
	Protocol    string
	LocalIP     string
	LocalPort   string
	RemoteIP    string
	RemotePort  string
	User        string
	Pid         int32
	ProgramName string
}
