/*
  TypePress server module.
*/
package server

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/braintree/manners"
	"github.com/codegangsta/martini"
	"github.com/typepress/core"
	"github.com/typepress/core/types"
)

var (
	AppVersion string
	srv        *manners.GracefulServer
	onceInit   sync.Once
	door       *doorHandler
	stopString string
)

func Run(m *martini.Martini, r martini.Router) error {
	if stopString == "" {
		stopString = time.Now().UTC().Format("15040501022006.000000000")
	}
	stopSig := types.NewStringSignal(stopString, nil)

	onceInit.Do(initOnce)
	m.Map(http.Dir(staticPath))

	srv = manners.NewServer()
	core.ListenSignal(sigReceive(srv.Shutdown), os.Interrupt, os.Kill, stopSig)

	stop := make(chan bool)

	go func() {
		stop <- true
		core.FireSignal(types.NewStringSignal(core.ServerShutDown, nil))
		srv.Shutdown <- true
	}()

	door = NewDoor(m, stop, false)

	return srv.ListenAndServe(laddr, door)
}

func Simple() error {
	return Run(core.Martini())
}

// 返回一个停止信号, 可用 core.FireSignal 触发信号通知服务器安全关闭.
func StopSignal() os.Signal {
	if stopString == "" {
		stopString = time.Now().UTC().Format("15040501022006.000000000")
	}
	return types.NewStringSignal(stopString, nil)
}

func initOnce() {
	LoadConfig()
}
