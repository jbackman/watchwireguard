package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	defaultWGHost      = "homefw.jitonline.net"
	defaultMonitorHost = "monitor.jitonline.net"
	defaultInterface   = "wg0"
)

var (
	debug   bool
	dest    string
	intdest string
	iface   string
)

func restartWireguard() {

	service := "wg-quick@" + iface
	if debug {
		fmt.Printf("Restarting %s", service)
	}
	_, err := exec.Command("systemctl", "restart", service).Output()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&dest, "wghost", defaultWGHost, "wireguard host")
	flag.StringVar(&intdest, "monhost", defaultMonitorHost, "monitor host")
	flag.BoolVar(&debug, "debug", false, "enable debugging")
	flag.StringVar(&iface, "iface", defaultInterface, "wg interface")

}

func main() {
	flag.Parse()
	ips, err := net.LookupIP(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		if debug {
			fmt.Printf("%s resolves to %s\n", dest, ip.String())
		}
	}
	out, err := exec.Command("wg", "show", iface, "endpoints").Output()
	if err != nil {
		log.Fatal(err)
	}
	m := regexp.MustCompile(`\s+`)
	if len(string(out)) > 0 {
		res, _, _ := strings.Cut(m.Split(string(out), 2)[1], ":")
		if debug {
			fmt.Printf("endpoint found: %s\n", res)
		}
		if res != ips[0].String() {
			restartWireguard()
		}
	} else {
		restartWireguard()
	}
	if debug {
		fmt.Printf("pinging %s\n", intdest)
	}
	out, _ = exec.Command("ping", intdest, "-c 2", "-i 3", "-w 1").Output()
	if strings.Contains(string(out), "100% packet loss") {
		if debug {
			fmt.Printf("cannot ping %s\n", intdest)
		}
	} else {
		if debug {
			fmt.Printf("successfully pinged %s\n", intdest)
		}
	}
}
