package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/braintree/manners"
)

// This struct contains the config keys, check config.toml for a short
// description of what each key does.
type config struct {
	IP       string
	Port     int
	Interval time.Duration
	Script   string
}

func main() {
	var conf config

	log.Println("loading config file")
	_, err := toml.DecodeFile("./config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}

start:
	// Check if http server is already running up, if this is the case
	// it means the other RPi is already running and the script is (most likely)
	// also working. Thus we make this RPi secondary. The secondary RPi will just
	// stay in this loop until the primary somehow goes down. If after 6 failed
	// GETs the primary is still not up, the secondary will take over the role
	// of primary.
	log.Println("checking http://" + conf.IP + " status")
	var e int
	var primary bool
	for !primary {
		// GET IP as specified in the config and check it for errors.
		_, err = http.Get("http://" + conf.IP)
		if err != nil {
			log.Println(err)

			// We'll take over the role of primary after 6 failes GETs.
			e++
			if e > 6 {
				log.Println("http://" + conf.IP + " is down, we are now primary")

				primary = true
				break
			}
		} else {
			// Reset error count to 0 if the other RPi is working properly.
			e = 0
			log.Println("http://" + conf.IP + " is working properly, we are secondary")
		}

		time.Sleep(4 * time.Second)
	}

	// Launch http server.
	log.Println("starting new http server on port " + strconv.Itoa(conf.Port))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Aasgier is currently running!")
	})
	go manners.ListenAndServe(":"+strconv.Itoa(conf.Port), mux)

	e = 0
	for {
		// Just to be sure we check if the (should be) secondary http
		// server isn't actually running, if this is the case something went
		// wrong. We don't want two concurrent http servers, because that
		// means we are executing the script two times on both the primary
		// and on the "secondary".
		//
		// Doing that means we could easily run into nasty problems such as the
		// barrier getting conflicting commands to close *and* to open at the same
		// time. We could easily fix this by making two separate functions, one for
		// "close barrier" and one for "open barrier", but hey. Doing that would
		// also make this entire script and part of the SNE part of our project
		// useless :^).
		//
		// TODO: Since we are running an http server anyways, I could easily parse
		// the output of the script (if it prints something sane, like something I
		// could parse with the toml decoder) info a nice looking websites with
		// graphs and stuff. That would make this entire script actually useful as
		// well.
		_, err = http.Get("http://" + conf.IP)
		if err == nil {
			log.Println("secondary is executing the scripts as well, ceasing to be primary")
			manners.Close()

			// Become secondary (hopefully).
			goto start
		}

		// Execute script and check if everything went well.
		log.Println("executing " + conf.Script)
		cmd := exec.Command(conf.Script)
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Println(err)

			// We'll accept 6 consecutive failed executions, after that
			// we will give up being primary, and the secondary will (hopefulle)
			// take over the primary role.
			e++
			if e > 6 {
				log.Println("script failed too many times, ceasing to be primary")
				manners.Close()

				// Become secondary (hopefully).
				time.Sleep(4 * time.Second * 2)
				goto start
			}
		} else {
			// Reset error count to 0 if the script executed properly.
			e = 0
		}

		time.Sleep(conf.Interval * time.Second)
	}
}
