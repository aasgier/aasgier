package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cenkalti/rpc2"
)

type config struct {
	IP       string
	Port     int
	Interval time.Duration
	Script   string
}

func client(i time.Duration, c *tls.Conn, s string) error {
	type Args struct{ S string }
	type Reply string

	clt := rpc2.NewClient(c)
	clt.Handle("script", func(client *rpc2.Client, args *Args, reply *Reply) error {
		*reply = Reply(args.S)

		return nil
	})

	for {
		go clt.Run()

		var r Reply
		clt.Call("script", Args{s}, &r)
		fmt.Println(r)

		time.Sleep(i * time.Second)
	}

	return nil
}

func server(p int, t *tls.Config, c *tls.Conn, s string) error {
	type Args struct{ S string }
	type Reply string

	srv := rpc2.NewServer()
	srv.Handle("script", func(client *rpc2.Client, args *Args, reply *Reply) error {
		var r Reply
		client.Call("script", Args{s}, &r)
		fmt.Println(r)

		*reply = Reply(args.S)

		return nil
	})

	l, err := tls.Listen("tcp", ":"+strconv.Itoa(p), t)
	if err != nil {
		return err
	}
	srv.Accept(l)

	return nil
}

func main() {
	// Load config.
	var conf config
	_, err := toml.DecodeFile("./config.toml", &conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config load ./config.toml: "+err.Error())
		os.Exit(1)
	}

	// Create a pool of trusted certs (in our case only our own).
	cp := x509.NewCertPool()
	f, err := ioutil.ReadFile("./certs/aasgier.pem")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cp.AppendCertsFromPEM(f)

	kp, err := tls.LoadX509KeyPair("./certs/aasgier.pem", "./certs/aasgier.key")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cert := tls.Config{
		Certificates:       []tls.Certificate{kp},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          cp,
		InsecureSkipVerify: true,
	}

	now := time.Now()
	cert.Time = func() time.Time { return now }
	cert.Rand = rand.Reader

	// Check if TCP is running on the other RPi, set as primary accordingly.
	var primary bool
	conn, err := tls.Dial("tcp", conf.IP, &cert)
	if err == nil {
		primary = true
	}

	if primary {
		// We actually did this before, but oh well...
		fmt.Println("[primary] Launching RPC client...\n")
		if err := client(conf.Interval, conn, conf.Script); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Println("[secondary] Launching RPC server...\n")
		if err := server(conf.Port, &cert, conn, conf.Script); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
