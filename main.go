package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"golang.org/x/net/websocket"

	"github.com/BurntSushi/toml"
	"github.com/braintree/manners"
)

// This struct contains the config keys, check config.toml for a short
// description of what each key does.
type config struct {
	IP       string
	Port     int
	Interval time.Duration
	History  int
	Script   string
}

// This struct contains the keys the script (should) print.
type script struct {
	Vibrate    bool
	WaterLevel int
}

var conf config
var scr script
var waterLevelList []int
var initTime = time.Now().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")

type message struct {
	Hostname       string `json:"hostname"`
	Uptime         string `json:"uptime"`
	Vibrate        bool   `json:"vibrate"`
	WaterLevelList []int  `json:"waterLevelList"`
}

func root(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("./http/index.html")
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "%s", f)
}

func socket(ws *websocket.Conn) {
	hn, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}

	for {
		// Send the waterLevelList to websocket
		if err := websocket.JSON.Send(ws, message{hn, initTime, scr.Vibrate, waterLevelList}); err != nil {
			log.Println(err)
			break
		}

		// Receive messages.
		var m message
		if err := websocket.JSON.Receive(ws, &m); err != nil {
			log.Println(err)
			break
		}

		// TODO: Use some kind of even here to continue the loop.
		// If I fix this I can remove a lot of "useless" code in script.js as well.
		time.Sleep(conf.Interval * time.Second / 2)
	}
}

func init() {
	log.Println("decoding config file")
	_, err := toml.DecodeFile("./config.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// If this is true (which it can only be on first run) we will
	// skip waiting for 6 failed GETs..
	init := true

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

		if _, err := http.Get("http://" + conf.IP); err != nil {
			log.Println(err)

			// We'll take over the role of primary after 6 failes GETs.
			e++
			if e > 6 || init {
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
	init = false

	log.Println("starting new http server on port " + strconv.Itoa(conf.Port))
	mux := http.NewServeMux()

	// Set location of our assets and websocket stuff.
	mux.Handle("/socket", websocket.Handler(socket))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./http/assets"))))

	// Read stuff into template, or something.
	mux.HandleFunc("/", root)

	// Start the server.
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
		if _, err := http.Get("http://" + conf.IP); err == nil {
			log.Println("secondary is executing the scripts as well, ceasing to be primary")
			manners.Close()

			// Become secondary (hopefully).
			goto start
		}

		// Execute script and check if everything went well.
		log.Println("executing " + conf.Script)
		cmd := exec.Command(conf.Script)
		var b bytes.Buffer
		cmd.Stdout = &b
		if err := cmd.Run(); err != nil {
			log.Println(err)

			// We'll accept 6 consecutive failed executions, after that
			// we will give up being primary, and the secondary will (hopefully)
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

		// Parse script output.
		if _, err := toml.Decode(b.String(), &scr); err != nil {
			log.Println(err)
		}
		waterLevelList = append(waterLevelList, scr.WaterLevel)
		if len(waterLevelList) >= conf.History+1 {
			waterLevelList = waterLevelList[1 : conf.History+1]
		}

		time.Sleep(conf.Interval * time.Second)
	}
}
