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

### ps

### pstree


