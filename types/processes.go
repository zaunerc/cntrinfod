package types

type ProcessInfo struct {
	Pid  int32
	User string
	Tty  string
	Name string
	// Exe returns executable path of the process.
	Exe string
	Cwd string
	Cmd string
}
