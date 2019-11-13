package gstats_test

import (
	"github.com/cevatbarisyilmaz/gstats"
	"log"
	"net"
	"net/http"
)

func Example() {
	gs := gstats.New("gstats")
	listener, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		log.Fatal(err)
	}
	glistener := gs.Listener(listener)
	http.Handle("/gstats/", gs.Collect(http.StripPrefix("/gstats", gs.Show())))
	http.Handle("/", gs.Collect(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("You are visiting " + request.URL.Path))
	})))
	err = http.Serve(glistener, nil)
	gs.PrepareToExit()
	log.Fatal(err)
}
