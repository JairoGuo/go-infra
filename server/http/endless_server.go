// @Title
// @Description
// @Author Jairo 2024/5/22 14:15
// @Email jairoguo@163.com

package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	PRE_SIGNAL = iota
	POST_SIGNAL

	STATE_INIT
	STATE_RUNNING
	STATE_SHUTTING_DOWN
	STATE_TERMINATE
)

var (
	runningServerReg     sync.RWMutex
	runningServers       map[string]*EndlessServer
	runningServersOrder  []string
	socketPtrOffsetMap   map[string]uint
	runningServersForked bool

	DefaultReadTimeOut    time.Duration
	DefaultWriteTimeOut   time.Duration
	DefaultMaxHeaderBytes int
	DefaultHammerTime     time.Duration

	isChild     bool
	socketOrder string

	hookableSignals []os.Signal
)

func init() {
	runningServers = make(map[string]*EndlessServer)
	runningServersOrder = []string{}
	socketPtrOffsetMap = make(map[string]uint)

	DefaultMaxHeaderBytes = 0 // use http.DefaultMaxHeaderBytes - which currently is 1 << 20 (1MB)

	DefaultHammerTime = 60 * time.Second

	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	}
}

type EndlessServer struct {
	http.Server
	listener    net.Listener
	SignalHooks map[int]map[os.Signal][]func()
	sigChan     chan os.Signal
	waitChan    chan struct{}
	isChild     bool
	transmit    map[string]bool
	state       uint8
	lock        sync.RWMutex
	BeforeBegin func(addr string)
	addrHook    func() []string
}

func NewServer(addr string, handler http.Handler) (srv *EndlessServer) {
	runningServerReg.Lock()
	defer runningServerReg.Unlock()

	socketOrder = os.Getenv("ENDLESS_SOCKET_ORDER")
	value := os.Getenv("ENDLESS_CONTINUE")
	isChild = value != ""
	transmit := make(map[string]bool)

	if len(socketOrder) > 0 {
		for i, addr := range strings.Split(socketOrder, ",") {
			socketPtrOffsetMap[addr] = uint(i)
			transmit[addr] = true
		}
	} else {
		socketPtrOffsetMap[addr] = uint(len(runningServersOrder))
	}

	srv = &EndlessServer{
		sigChan:  make(chan os.Signal),
		waitChan: make(chan struct{}),
		isChild:  isChild,
		transmit: transmit,
		SignalHooks: map[int]map[os.Signal][]func(){
			PRE_SIGNAL:  make(map[os.Signal][]func()),
			POST_SIGNAL: make(map[os.Signal][]func()),
		},
		state: STATE_INIT,
	}

	srv.Server.Addr = addr
	srv.Server.ReadTimeout = DefaultReadTimeOut
	srv.Server.WriteTimeout = DefaultWriteTimeOut
	srv.Server.MaxHeaderBytes = DefaultMaxHeaderBytes
	srv.Server.Handler = handler

	srv.BeforeBegin = func(addr string) {
		log.Println(syscall.Getpid(), addr)
	}

	runningServersOrder = append(runningServersOrder, addr)
	runningServers[addr] = srv

	return
}

func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}

func (s *EndlessServer) getState() uint8 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.state
}

func (s *EndlessServer) setState(st uint8) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.state = st
}

func (s *EndlessServer) setTransmit(addr string, state bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.transmit[addr] = state
}

func (s *EndlessServer) clearTransmit() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.transmit = make(map[string]bool)
}

func (s *EndlessServer) Serve() error {
	defer log.Println(syscall.Getpid(), "Serve() returning...")
	s.setState(STATE_RUNNING)
	err := s.Server.Serve(s.listener)
	log.Println(syscall.Getpid(), "Waiting for connections to finish...")
	s.setState(STATE_TERMINATE)
	return err
}

func (s *EndlessServer) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":http"
	}

	go s.handleSignals()

	l, err := s.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	s.listener = l

	if s.isChild {
		syscall.Kill(syscall.Getppid(), syscall.SIGTERM)
	}

	s.BeforeBegin(s.Addr)

	return s.Serve()
}

func (s *EndlessServer) ListenAndServeTLS(certFile, keyFile string) error {
	addr := s.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if s.TLSConfig != nil {
		*config = *s.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	config.Certificates = make([]tls.Certificate, 1)
	var err error
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	go s.handleSignals()

	l, err := s.getListener(addr)
	if err != nil {
		log.Println(err)
		return err
	}

	s.listener = tls.NewListener(l, config)

	if s.isChild {
		syscall.Kill(syscall.Getppid(), syscall.SIGTERM)
	}

	log.Println(syscall.Getpid(), s.Addr)
	return s.Serve()
}

func (s *EndlessServer) getListener(laddr string) (net.Listener, error) {
	// TODO 如果修改端口，不使用父进程的文件描述符
	if s.isChild && s.transmit[laddr] {
		ptrOffset := socketPtrOffsetMap[laddr]

		f := os.NewFile(uintptr(3+ptrOffset), "")
		l, err := net.FileListener(f)
		if err != nil {
			return nil, fmt.Errorf("net.FileListener error: %v", err)
		}

		addr, err := net.ResolveTCPAddr("", laddr)

		if l.Addr().(*net.TCPAddr).Port == addr.Port {
			return l, nil
		}
	}

	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, fmt.Errorf("net.Listen error: %v", err)
	}
	return l, nil
}

func (s *EndlessServer) handleSignals() {
	signal.Notify(s.sigChan, hookableSignals...)
	pid := syscall.Getpid()
	for {
		sig := <-s.sigChan
		s.signalHooks(PRE_SIGNAL, sig)
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "Received SIGHUP. Forking.")
			s.clearTransmit()
			for _, addr := range runningServersOrder {
				s.setTransmit(addr, true)
			}
			if err := s.fork(); err != nil {
				log.Println("Fork error:", err)
			}
		case syscall.SIGUSR1:
			log.Println(pid, "Received SIGUSR1.")
			addrList := s.getAddr()
			if addrList != nil {
				s.clearTransmit()
				for _, addr := range addrList {
					if !s.contains(addr) {
						s.setTransmit(addr, false)
					} else {
						s.setTransmit(addr, true)
					}
				}
			}

			if err := s.fork(); err != nil {
				log.Println("Fork error:", err)
			}
		case syscall.SIGUSR2:
			log.Println(pid, "Received SIGUSR2.")
			s.clearTransmit()
			for _, addr := range runningServersOrder {
				s.setTransmit(addr, false)
			}
			if err := s.fork(); err != nil {
				log.Println("Fork error:", err)
			}
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println(pid, "Received SIGINT/SIGTERM.")
			err := s.shutdown()
			if err != nil {
				log.Println(pid, "Error during shutdown:", err)
			}
		case syscall.SIGTSTP:
			log.Println(pid, "Received SIGTSTP.")
		default:
			log.Printf("Received %v: nothing I care about...\n", sig)
		}
		s.signalHooks(POST_SIGNAL, sig)
	}
}

func (s *EndlessServer) signalHooks(ppFlag int, sig os.Signal) {
	if hooks, exists := s.SignalHooks[ppFlag][sig]; exists {
		for _, f := range hooks {
			f()
		}
	}
}

func (s *EndlessServer) shutdown() error {
	if s.getState() != STATE_RUNNING {
		return nil
	}

	s.setState(STATE_SHUTTING_DOWN)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultHammerTime)
	defer cancel()

	if err := s.listener.Close(); err != nil {
		log.Println(syscall.Getpid(), "Listener.Close() error:", err)
		return err
	}
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Println(syscall.Getpid(), "Server.Shutdown() error:", err)
	}

	log.Println(syscall.Getpid(), s.listener.Addr(), "Listener closed.")

	s.SetKeepAlivesEnabled(false)
	s.setState(STATE_TERMINATE)

	for i, s2 := range runningServersOrder {
		if s2 == s.Addr {
			runningServersOrder = remove(runningServersOrder, i)
			break
		}
	}
	delete(runningServers, s.Addr)
	close(s.waitChan)
	return nil
}

func remove(slice []string, i int) []string {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func (s *EndlessServer) fork() error {
	runningServerReg.Lock()
	defer runningServerReg.Unlock()

	// TODO
	//if runningServersForked {
	//	return errors.New("Another process already forked. Ignoring this one.")
	//}

	runningServersForked = true
	var fs []*os.File

	os.Unsetenv("ENDLESS_CONTINUE")
	os.Unsetenv("ENDLESS_SOCKET_ORDER")
	os.Clearenv()

	transmitCount := 0
	transmitAddr := []string{}
	transmitAddrOrderMap := make(map[string]uint)

	fmt.Println("fork", s.transmit)
	for addr, state := range s.transmit {
		if state && runningServers[addr] != nil {
			transmitAddr = append(transmitAddr, addr)
			transmitAddrOrderMap[addr] = uint(transmitCount)
			transmitCount++

		}
	}

	if transmitCount > 0 {
		runningServersOrderStr := strings.Join(transmitAddr, ",")

		files := make([]*os.File, len(transmitAddr))
		for _, srv := range runningServers {
			if s.transmit[srv.Addr] {
				file, err := srv.getListenerFile()
				if err != nil {
					return err
				}
				files[transmitAddrOrderMap[srv.Addr]] = file

			}
		}
		fs = append([]*os.File{os.Stdin, os.Stdout, os.Stderr}, files...)
		os.Setenv("ENDLESS_CONTINUE", "1")
		os.Setenv("ENDLESS_SOCKET_ORDER", runningServersOrderStr)
	} else {

		fs = append([]*os.File{os.Stdin, os.Stdout, os.Stderr})
		os.Setenv("ENDLESS_CONTINUE", "2")

	}

	path := os.Args[0]
	args := append([]string{path}, os.Args[1:]...)

	ppid := syscall.Getpid()
	env := os.Environ()

	process, err := os.StartProcess(path, args, &os.ProcAttr{
		Dir:   "",
		Env:   env,
		Files: fs,
		Sys:   &syscall.SysProcAttr{},
	})
	if err != nil {
		return err
	}

	log.Println(ppid, "Forked child", process.Pid)
	return nil
}

func (s *EndlessServer) getListenerFile() (*os.File, error) {
	tcpListener := s.listener.(*net.TCPListener)
	file, err := tcpListener.File()
	if err != nil {
		return nil, fmt.Errorf("getListenerFile: %v", err)
	}
	return file, nil
}

func (s *EndlessServer) BindAddrHook(f func() []string) {
	s.addrHook = f
}

func (s *EndlessServer) consistent(addr string) bool {
	tcpAddr, _ := net.ResolveTCPAddr("", addr)

	if s.listener.Addr().(*net.TCPAddr).Port != tcpAddr.Port {
		return false
	}

	return true
}

func (s *EndlessServer) contains(addr string) bool {
	tcpAddr, _ := net.ResolveTCPAddr("", addr)

	if value, ok := runningServers[addr]; ok {
		addr := value.listener.Addr().(*net.TCPAddr)
		if addr.Port != tcpAddr.Port {
			return false
		}
	} else {
		return false
	}
	return true
}

func (s *EndlessServer) getAddr() []string {
	if s.addrHook == nil {
		return nil
	}
	addr := s.addrHook()
	return addr
}

func (s *EndlessServer) getWait() chan struct{} {
	return s.waitChan
}
