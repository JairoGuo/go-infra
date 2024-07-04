// @Title
// @Description
// @Author Jairo 2024/5/7 13:10
// @Email jairoguo@163.com

package http

import (
	"log"
	"net/http"
	"syscall"
)

type Provision interface {
}

type Server struct {
	server        *EndlessServer
	enableEndless bool
	addr          string
	handler       http.Handler
	hook          func() []string
}

func Start(addr string, handler http.Handler) {

	err := http.ListenAndServe(addr, handler)
	if err != nil {
		panic(err)
	}

}

func New(addr string) *Server {
	return &Server{
		addr:          addr,
		enableEndless: true,
	}

}

func (s *Server) BindAddr(addr string) {
	s.addr = addr
}

func (s *Server) BindHandler(handler http.Handler) {
	s.handler = handler
}

func (s *Server) Hook(f func() []string) {
	s.hook = f
}

func (s *Server) EnableEndless(e bool) {
	s.enableEndless = e
}

func (s *Server) Start() error {

	if !s.enableEndless {
		server := http.Server{
			Addr:    s.addr,
			Handler: s.handler,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server error: %v\n", err)
		}
		return nil
	}

	s.server = NewServer(s.addr, s.handler)

	s.server.BindAddrHook(s.hook)

	if err := s.server.ListenAndServe(); err != nil {
		log.Printf("Server error: %v\n", err)
	}

	return nil
}

func (s *Server) Stop() error {
	err := s.server.shutdown()
	return err
}

func (s *Server) Restart() error {
	err := s.Stop()
	if err != nil {
		log.Printf("HTTP server stop: %v\n", err)
		return err
	}
	return s.Start()

}

func (s *Server) RestartByFork() error {
	pid := syscall.Getpid()
	log.Println(pid, "RestartByFork")

	addrList := s.server.getAddr()
	if addrList != nil {
		s.server.clearTransmit()
		for _, addr := range addrList {
			if !s.server.contains(addr) {
				s.server.setTransmit(addr, false)
			} else {
				s.server.setTransmit(addr, true)
			}
		}

	}

	if err := s.server.fork(); err != nil {
		log.Println("Fork error:", err)
	}

	return nil

}

func (s *Server) Wait() {
	log.Println("Server wait")
	<-s.server.getWait()
	log.Println("Server exit")
}
