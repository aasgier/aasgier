[![Go Report Card](https://goreportcard.com/badge/github.com/aasgier/aasgier)](https://goreportcard.com/report/github.com/aasgier/aasgier)

## Installation

```shell
go get github.com/aasgier/aasgier
```

This program depends on `python3` and `go`.

Make sure to edit `config.toml` to suit your needs.


## Usage

Execute the `aasgier` binary on one or more systems, again, make sure to
edit `config.toml` to make it point to the right IP address and such.


## Authors

* Camille Scholtz ([onodera-punpun](https://github.com/onodera-punpun))
* Eelke Feitz
* Fabiënne Oskam
* Joost van Bussum
* Ismael Hassan


## Notes

* Probably runs on every OS, but only tested on GNU/Linux.
* Website probably works on every browser, but only tested on Firefox.

---

# Technical design

This program includes multiple components. I will describe the different component below and how they work.


## Script that interfaces with the Raspberry Pi GPIO
#### [GPIO/GPIO.py](https://github.com/aasgier/aasgier/blob/master/GPIO/GPIO.py)

In our case this script is written in `python3` but of course, you can use any language you want. This script is used
to convert the data send by the different sensors such as the water sensor to something usable. For example
the water sensor initially returns a value from 0 to 1023, but this script converts this integer to a percentage value
that ranger from 0 to 100. Where 0 is no water touching the water sensor, and 100 is the value where the floodgate will
close, which in our case is a 750. These values get returned by simple functions. One such example could be
`getWaterLevel()` (`int`).


## Script that calls the hardware interfacing functions and decide if the gate should close or not
#### [script/aasgier.py](https://github.com/aasgier/aasgier/blob/master/script/aasgier.py)

Again, since we try to work modular you could theoretically write this in any language, but we use `python3`.
What this script does is call the various functions provided by the earlier mentioned script, these functions could
include things such as `getWaterLevel()` (`int`), `getWindLevel()` (`int`) or `isGateClosed()` (`bool`).

The idea of this script is that it evaluates the received values and decides if the floodgate should close (or open).
It then calls a function from the GPIO script to close the gate if needed. It also sends some data to the http-server.


## Program that creates the http-server, receives data from the above mentioned script, and sends data to website via websockets
#### [main.go](https://github.com/aasgier/aasgier/blob/master/main.go)

Since I'm learning `go` I've decided to write this in `go`. This not only serves as an opportunity to learn about `go`
but also to test the "modularity" of our program(s). You could also do this in `python3` or any other language, but
libraries for websockets and hosting a http server need to be available for the language. Moreover I make use of a
bit of concurrency in order to both host the website, send data over websockets, and receive data from the
above mentioned scripts.

In short this program does a few things:

* Read a config with various options specified by the administrator, these options include the port to run the
  http-server on, the IP and port of the other measuring station, the interval at which to update the webpage and the
  interval at which to check the sensor 

* Check if there already is a http-server is already running, it this is the case that means that the program will become
  "secondary", this means that it will idle, check at an interval if the http server is still running, if this is *not* the
  case it means something is wrong. This could be various things such as the "primary" not having an internet, the script
  could've crashed and more.

* If there is no http server running the program will launch a new one using the specified port found in the config file. the
  machine/program that runs the http server is by definition the "primary".

* Execute a script at an interval in a loop, the script gets defined in `config.toml`, in our case `script/aasgier.py` gets executed.

* It evaluates the information the script prints to stdout, though there is probably a better way of doing this (I'll look into it :)), and sends the information over websockets to the website, where a graph gets generated from it.


## The website itself
#### [http/index.html](https://github.com/aasgier/aasgier/blob/master/http/index.html)
#### [http/assets/style.css](https://github.com/aasgier/aasgier/blob/master/http/assets/style.css)
#### [http/assets/script.js](https://github.com/aasgier/aasgier/blob/master/http/assets/script.js)

I've written this in a mix of `html`, `css` and `javascript`. Below I briefly describe what the various language
components do.

* The `html` portion of the website just creates the general layout, only the static text gets defined in the html.

* The `css` is of course used to style the website and make it pretty.

* The `javascript` part of the website it used to receive data that is send via websockets. This data gets used to
  genrate the graph (using `graph.js`) and to fill in the various dynamic part of the website, such as the div that
  displays the current water level. I also use `jquery` to make this a bit simpler (in terms of writing), but of course
  his isn't strictly necessary.
