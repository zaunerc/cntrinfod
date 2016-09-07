# [cntinsight ](https://github.com/zaunerc/cntinsight)

* [All log files](log)

## Host info

* Hostname: {{.HostHostname}}
* [Further container host info](hostinfo)

## Container info

* Hostname: {{.ContainerHostname}}

### TCP4/6 Listening Sockets

| Proto | Local Address | Foreign Address | State | User | PID | Program name |
|-------|---------------|-----------------|-------|------|-----|--------------|
{{range .TcpSocketInfo}}| {{.Protocol}} | {{.LocalIP}}:{{.LocalPort}} | {{.RemoteIP}}:{{.RemotePort}} | {{.State}} | {{.User}} | {{.Pid}} | {{.ProgramName}} |
{{end}}

### UDP4/6 Sockets

| Proto | Local Address | Foreign Address | User | PID | Program name |
|-------|---------------|-----------------|------|-----|--------------|
{{range .UdpSocketInfo}}| {{.Protocol}} | {{.LocalIP}}:{{.LocalPort}} | {{.RemoteIP}}:{{.RemotePort}} | {{.User}} | {{.Pid}} | {{.ProgramName}} |
{{end}}

### Running processes

| PID | User | Name |  Current workding dir | Command line | Terminal |
|-----|------|------|-----------------------|--------------|----------|
{{range .ProcessInfo}}| {{.Pid}} | {{.User}} | {{.Name}} | {{.Cwd}} | {{.Cmd}} | {{.Tty}} |
{{end}}

### pstree


