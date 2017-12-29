package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/ironzhang/coap"
	"github.com/ironzhang/x/collect"
)

type Server struct {
	svr coap.Server
}

func (s *Server) ListenAndServe(address string) error {
	s.svr.ReadBytes = 512 * 1024
	s.svr.Handler = s
	return s.svr.ListenAndServe(address)
}

func (s *Server) ServeCOAP(w coap.ResponseWriter, r *coap.Request) {
	collectf := collect.Default.Collect("main.Server.ServeCOAP")
	defer collectf()

	switch r.URL.Path {
	case "/echo":
		go s.echo(r.RemoteAddr.String(), 1)
	}
}

func (s *Server) echo(addr string, n int) {
	req, err := coap.NewRequest(true, coap.POST, fmt.Sprintf("coap://%s/echo", addr), nil)
	if err != nil {
		log.Printf("new request: %v", err)
		return
	}
	for i := 0; i < n; i++ {
		if _, err = s.svr.SendRequest(req); err != nil {
			log.Printf("send request: %v", err)
			break
		}
	}

	req, err = coap.NewRequest(true, coap.POST, fmt.Sprintf("coap://%s/echoFinish", addr), nil)
	if err != nil {
		log.Printf("new request: %v", err)
		return
	}
	if _, err = s.svr.SendRequest(req); err != nil {
		log.Printf("send request: %v", err)
		return
	}
}

func main() {
	coap.Verbose = 0
	coap.EnableCache = false
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go collect.PrintDefaultCollect(os.Stdout, 5*time.Second, true)

	go func() {
		var _ = pprof.Index
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatal(err)
		}
	}()

	var s Server
	if err := s.ListenAndServe(":5683"); err != nil {
		log.Fatal(err)
	}
}
