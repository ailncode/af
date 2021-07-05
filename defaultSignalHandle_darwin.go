// Copyright 2018 ailn(ailnindex@qq.com). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package af provides easy to use graceful restart a http server

// +build darwin

package af

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

//Give some default signal handler if AF' signalHandlerMap is nil
func (af *AF) defaultSignalHandle() {
	af.HandleSignal(func(af *AF) {
		log.Println(fmt.Sprintf("AF server is shutdown pid:%d", os.Getpid()))
		af.Stop()
	}, syscall.SIGINT, syscall.SIGTERM)
	af.HandleSignal(func(af *AF) {
		log.Println(fmt.Sprintf("AF server is reload pid:%d err:%v", os.Getpid(), af.Reload()))
	}, syscall.SIGUSR2)
}
