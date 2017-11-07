package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/ironzhang/coap"
)

type Server struct {
	svr coap.Server
}

func (s *Server) ListenAndServe(address string) error {
	s.svr.Handler = s
	return s.svr.ListenAndServe(address)
}

func (s *Server) ServeCOAP(w coap.ResponseWriter, r *coap.Request) {
	//log.Print(r.URL.String())
	switch r.URL.Path {
	case "/ping":
		go s.ping(r.RemoteAddr.String(), 1)
	case "/observe":
		go s.observe(r.RemoteAddr.String(), 1)
	}
}

func (s *Server) ping(addr string, n int) {
	req, err := coap.NewRequest(true, coap.POST, fmt.Sprintf("coap://%s/ping", addr), nil)
	if err != nil {
		log.Printf("new coap request: %v", err)
		return
	}

	for i := 0; i < n; i++ {
		_, err := s.svr.SendRequest(req)
		if err != nil {
			log.Printf("send coap request: %v", err)
			break
		}
	}
}

func (s *Server) observe(addr string, n int) {
	urlstr := fmt.Sprintf("coap://%s/observe", addr)
	for i := 0; i < n; i++ {
		if err := s.svr.Observe(genToken(), urlstr, 0); err != nil {
			log.Printf("coap observe: %v", err)
			break
		}
	}
}

func genToken() coap.Token {
	b := make([]byte, 8)
	rand.Read(b)
	return coap.Token(b)
}

func main() {
	coap.Verbose = 0
	coap.EnableCache = false

	var s Server
	if err := s.ListenAndServe(":5683"); err != nil {
		log.Fatal(err)
	}
}
