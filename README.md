[![Go Report Card](https://goreportcard.com/badge/github.com/aasgier/aasgier)](https://goreportcard.com/report/github.com/aasgier/aasgier)

## Installation

This program depends on `python3` and `go`.

```shell
go get -u -v -fix github.com/aasgier/aasgier
cd $GOPATH/src/github.com/aasgier/aasgier
go build -v
```

if you don't have a `$GOPATH` set, `go` will use `$HOME/go`.


## Usage

Execute the `$GOPATH/src/github.com/aasgier/aasgier/aasgier` binary on one or more systems, make sure to
edit `config.toml` to make it point to the right IP address and such.


## Authors

* Camille Scholtz ([onodera-punpun](https://github.com/onodera-punpun))
* Eelke Feitz ([EelkeFeitz](https://github.com/EelkeFeitz))
* Fabiënne Oskam
* Joost van Bussum
* Ismael Hassan


## Notes

* Probably runs on every OS, but only tested on GNU/Linux.
* Website probably works on every browser, but only tested on Firefox.


---

# Technical design

This program includes multiple components. I will describe the different component below and how they work.


## Program that interfaces with the Raspberry Pi GPIO pins
#### [GPIO/GPIO.py](https://github.com/aasgier/aasgier/blob/master/GPIO/GPIO.py)

In our case this program is written in `python3` but of course, you can use any language you want. This program is used to convert the data send by the different sensors such as the water sensor to something usable. For example the water sensor initially returns a value between 0 and 1023, this program converts this integer to an interger that ranges from 0 to 100. Where 0 is no water touching the water sensor, and 100 is the value where the gate will close. In our case 100 equals 750. These values get returned by simple functions the other program can call. One such function could be `getWaterLevel()` (`int`).

This program also has a few functions that only return an error an do actual tasks, such as closing the gate with for example `closeGate` (`err`).


## Program that calls the hardware interfacing functions and decide if the gate should close or not
#### [script/aasgier.py](https://github.com/aasgier/aasgier/blob/master/script/aasgier.py)

Again, since we try to work modular you could theoretically write this in any language, but we use `python3`. What this script does is call the various functions provided by the earlier mentioned program, these functions could include things such as `getWaterLevel()` (`int`), `getWindLevel()` (`int`) or `isGateClosed()` (`bool`).

The idea of this program is that it evaluates the received values and decides if the gate should close (or open) depending on the combination of circumstances. It then calls a function from the GPIO program to close the gate if needed. It also sends some data such as the waterlevel and if the gate is closed to the HTTP-server.


## Program that creates the HTTP-server, receives data from the above mentioned program, and sends data to the website via websockets
#### [main.go](https://github.com/aasgier/aasgier/blob/master/main.go)

Since I'm learning `go` I've decided to write this in `go`. This not only serves as an opportunity to learn about `go` but also to test the "modularity" of our program(s). You could also do this in `python3` or any other language, but libraries for websockets and hosting a HTTP-server need to be available for the language. Moreover I make use of a bit of concurrency in order to both host the website, send data over websockets, and receive data from the above mentioned programs.

In short this program does a few things:

* Read a config with various options specified by the administrator, these options include the port to run the HTTP-server on, the IP and port of the other measuring-station, the interval at which to update the webpage and the interval at which to check the sensor.

* Check if there already is a HTTP-server running, it this is the case that means that the program/measuring-station will become "secondary". This means that it will idle and check at an interval if the HTTP-server is still running, if this is *not* the case it means something is wrong. This could be various things such as the "primary" not having an internet connection, the program could've crashed, et cetera.

* If there is no HTTP-server running the program will launch a new one using the specified port from the config file. the program/mesuring-station that runs the HTTP-server is by definition the "primary".

* Execute a script at an interval in a loop, the program gets defined in the config, in our case `script/aasgier.py` gets executed.

* It evaluates the information the above mentioned program prints to stdout (though there is probably a better way of doing this (I'll look into it `:)`)), and sends the information over websockets to the website, where a graph gets generated from it.


## The website itself
#### [HTTP/index.html](https://github.com/aasgier/aasgier/blob/master/HTTP/index.html)
#### [HTTP/assets/style.css](https://github.com/aasgier/aasgier/blob/master/HTTP/assets/style.css)
#### [HTTP/assets/script.js](https://github.com/aasgier/aasgier/blob/master/HTTP/assets/script.js)

I've written this in a mix of `html`, `css` and `javascript`.

* The `html` portion of the website just creates the general layout, only the static text gets defined in the html.

* The `css` is of course used to style the website and make it all pretty.

* The `javascript` part of the website it used to receive data that is send via websockets. This data gets used to generate the graph (using `graph.js`) and to fill in the various dynamic parts of the website such as the div that displays the current water level. I also use `jquery` to make this a bit simpler (in terms of writing), but of course this isn't strictly necessary.
