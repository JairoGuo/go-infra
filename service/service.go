// @Title
// @Description
// @Author Jairo 2023/12/11 17:17
// @Email jairoguo@163.com

package service

import (
	"fmt"
	"github.com/jairoguo/go-infra/util/path"
	"github.com/kardianos/service"
	"os"
)

type Program struct {
	instance      service.Service
	programHandle func()
}

func (p *Program) Start(s service.Service) error {
	go p.programHandle()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	return nil
}

func (p *Program) InitService(config service.Config) service.Service {

	exePath := path.GetExecutePath()
	os.Chdir(exePath)

	prg := &Program{}
	s, err := service.New(prg, &config)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	p.instance = s
	return s
}

func (p *Program) BindProgram(f func()) {
	p.programHandle = f
}

func (p *Program) RunService() {
	s := p.instance

	if len(os.Args) > 1 {
		err := service.Control(s, os.Args[1])
		if err != nil {
			fmt.Printf("service Control failed, err: %v\n", err)
			os.Exit(1)
		}
		return
	}

	err := s.Run()
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func New() *Program {
	return &Program{}
}
