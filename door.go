package server

import (
	"net/http"
	"sync"
)

// doorHandler 实现了 http.Handler 接口, 对活动的 http.Request 计数, 安全关闭.
type doorHandler struct {
	http.Handler
	notify chan bool
	mux    sync.Mutex
	count  int
	stop   bool
	io     bool
}

/*
  NewDoor 返回可安全关闭的 http.Handler.

  参数:

  hh  是要包裹的 http.Handler, 安全关闭后通知 notify 通道.
  io  指示具体访问通道的方向:
	- true  表示 input,  c.notify <- true
	- false 表示 output, <-c.notify
*/
func NewDoor(hh http.Handler, notify chan bool, io bool) *doorHandler {
	d := &doorHandler{
		Handler: hh,
		notify:  notify,
	}
	return d
}

// http.Handler 接口
func (c *doorHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c.mux.Lock()
	if c.stop {
		c.mux.Unlock()
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	c.count++
	c.mux.Unlock()
	defer c.out()
	c.Handler.ServeHTTP(res, req)
}

func (c *doorHandler) out() {
	c.mux.Lock()
	c.count--
	if c.stop && c.count == 0 {
		if c.io {
			go func() { c.notify <- true }()
		} else {
			go func() { <-c.notify }()
		}
	}
	c.mux.Unlock()
}

// 通知doorHandler开始阻止请求, 此后所有的请求都会被 StatusServiceUnavailable.
// 所有活动请求都完成后通知 notify 通道.
func (c *doorHandler) Stop() {
	c.mux.Lock()
	c.stop = true
	c.mux.Unlock()
}

// Count 返回活动的 http.Request 数量
func (c *doorHandler) Count() (count int) {
	c.mux.Lock()
	count = c.count
	c.mux.Unlock()
	return
}
