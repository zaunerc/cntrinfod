package system

import (
	"fmt"
	"github.com/drael/GOnetstat"
	"os"
)

func FetchContainerHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

/*
 * FetchNetstatTcp gets TCP information and show like netstat.
 * Information like 'user' and 'name' of some processes will
 * not show if you don't have root permissions
 */
func FetchNetstatTcp() string {

	d := GOnetstat.Tcp()

	// format header
	fmt.Printf("Proto %16s %20s %14s %24s\n", "Local Adress", "Foregin Adress",
		"State", "Pid/Program")

	for _, p := range d {

		// Check STATE to show only Listening connections
		if p.State == "LISTEN" {
			// format data like netstat output
			ip_port := fmt.Sprintf("%v:%v", p.Ip, p.Port)
			fip_port := fmt.Sprintf("%v:%v", p.ForeignIp, p.ForeignPort)
			pid_program := fmt.Sprintf("%v/%v", p.Pid, p.Name)

			fmt.Printf("tcp %16v %20v %16v %20v\n", ip_port, fip_port,
				p.State, pid_program)
		}
	}

	return ""
}

/*
 * FetchNetstatTcp6 gets TCP information and show like netstat.
 * Information like 'user' and 'name' of some processes will
 * not show if you don't have root permissions
 */
func FetchNetstatTcp6() string {

	d := GOnetstat.Tcp6()

	// format header
	fmt.Printf("Proto %16s %20s %14s %24s\n", "Local Adress", "Foregin Adress",
		"State", "Pid/Program")

	for _, p := range d {

		// Check STATE to show only Listening connections
		if p.State == "LISTEN" {
			// format data like netstat output
			ip_port := fmt.Sprintf("%v:%v", p.Ip, p.Port)
			fip_port := fmt.Sprintf("%v:%v", p.ForeignIp, p.ForeignPort)
			pid_program := fmt.Sprintf("%v/%v", p.Pid, p.Name)

			fmt.Printf("tcp %16v %20v %16v %20v\n", ip_port, fip_port,
				p.State, pid_program)
		}
	}

	return ""
}
