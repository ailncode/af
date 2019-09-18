// Copyright 2018 ailn(ailnindex@qq.com). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package af provides easy to use graceful restart a http server

package af

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"time"
)

type AF struct {
	Addr             string
	Handler          http.Handler
	ShutdownTimeOut  time.Duration
	listener         net.Listener
	server           *http.Server
	signalSlice      []os.Signal
	signalHandlerMap map[os.Signal]func(*AF)
	graceful         bool
}

//Get default AF
func Default() *AF {
	return &AF{
		Addr:            ":8080",
		ShutdownTimeOut: time.Second * 10,
	}
}

//Get AF with a http handler
func NewWithHandler(handler http.Handler) *AF {
	return &AF{
		Addr:            ":8080",
		Handler:         handler,
		ShutdownTimeOut: time.Second * 10,
	}
}

//Get AF with a http server
func NewWithServer(server *http.Server) *AF {
	return &AF{
		server:          server,
		Addr:            server.Addr,
		Handler:         server.Handler,
		ShutdownTimeOut: time.Second * 10,
	}
}

//init AF
func (af *AF) init() error {
	if af.server == nil {
		if af.Handler == nil {
			af.Handler = http.DefaultServeMux
		}
		af.server = &http.Server{
			Addr:         af.Addr,
			Handler:      af.Handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	}
	var err error
	if af.graceful {
		f := os.NewFile(3, "")
		af.listener, err = net.FileListener(f)
	} else {
		af.listener, err = net.Listen("tcp", af.server.Addr)
	}
	return err
}

//signal handle
func (af *AF) signalHandle() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, af.signalSlice...)
	for {
		select {
		case s := <-signalChan:
			if f, ok := af.signalHandlerMap[s]; ok {
				f(af)
			}
			break
		}
	}
}

//Give a handler for signals
func (af *AF) HandleSignal(handler func(*AF), signals ...os.Signal) {
	if af.signalHandlerMap == nil {
		af.signalHandlerMap = make(map[os.Signal]func(*AF))
	}
	for _, s := range signals {
		if _, ok := af.signalHandlerMap[s]; !ok {
			af.signalSlice = append(af.signalSlice, s)
		}
		af.signalHandlerMap[s] = handler
	}
}

//Run the AF
func (af *AF) Run() error {
	flag.BoolVar(&af.graceful, "graceful", false, "listen on fd open 3 (internal use only)")
	flag.Parse()
	var err error
	if err = af.init(); err != nil {
		return err
	}
	go func(af *AF) {
		if af.signalHandlerMap == nil {
			af.defaultSignalHandle()
		}
		af.signalHandle()
	}(af)

	log.Println(fmt.Sprintf("AF server is run at pid:%d", os.Getpid()))
	err = af.server.Serve(af.listener)
	return err
}

//Stop the AF
func (af *AF) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), af.ShutdownTimeOut)
	return af.server.Shutdown(ctx)
}

//Reload the AF
func (af *AF) Reload() error {
	listener, ok := af.listener.(*net.TCPListener)
	if !ok {
		return errors.New("AF.listener type must be *net.TCPListener,but got type " + reflect.TypeOf(af.listener).Name())
	}
	file, err := listener.File()
	if err != nil {
		return err
	}
	args := os.Args
	if !af.graceful {
		args = append(args, "-graceful")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{file}
	err = cmd.Start()
	if err == nil {
		return af.Stop()
	}
	return err
}
