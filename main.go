package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type config struct {
	IP       string
	Port     int
	Interval time.Duration
	Script   string
}

func main() {
	// Load config.
	var conf config
	_, err := toml.DecodeFile("./config.toml", &conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config load ./config.toml: "+err.Error())
		os.Exit(1)
	}

	// Check if TCP is running on the other RPi, set as primary accordingly.
	var primary bool
	conn, err := net.Dial("tcp", conf.IP)
	if err == nil {
		primary = true
	}

	if primary {
		// We actually did this before, but oh well...
		fmt.Println("[primary] Launching TCP dialer...\n")

		for {
			fmt.Println("\n[primary] Running script on primary...")
			fmt.Fprintf(conn, "~[secondary] Executing script on primary...\n")

			cmd := exec.Command(conf.Script)
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "[primary] Tried to execute script but failed!")
				fmt.Fprintf(conn, "[secondary] Tried to execute script but failed!\n")
				fmt.Fprintln(os.Stderr, "[primary] Changing secondary into primary!\n")
				fmt.Fprintf(conn, "[secondary] Changing secondary into primary!~\n")
				// TODO
			}

			fmt.Println("[primary] Script executed on primary without any problems.")
			fmt.Fprintf(conn, "[secondary] Script executed without any problems.\n")

			time.Sleep(conf.Interval * time.Second)
		}
	} else {
		fmt.Println("[secondary] Launching TCP listener...\n")

		ln, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Port))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		conn, err = ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for {
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Print(strings.Replace(msg, "~", "\n", -1))
		}
	}
}
