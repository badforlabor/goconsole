/**
 * Auth :   liubo
 * Date :   2018/11/19 14:58
 * Comment: 
 */

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"goconsole/console"
	"net/http"
	"runtime"
)

var port = flag.Int("p", 8888, "Listen port")
var ip = flag.String("ip", "0.0.0.0", "Listen ip")
var name = flag.String("name", "shell", "Driver name, shell docker and ssh")
var host = flag.String("host", "", "Connect to the host address, It has to be docker or SSH Driver")
var cid = flag.String("cid", "", "Docker Container id, It has to be docker Driver")
var cmd = flag.String("cmd",
	func() string {
		if runtime.GOOS != "windows" {
			return "sh"
		}
		return "cmd"
	}(),
	"command to execute")
var disable = flag.Bool("d", false, "Disable url parameters")

type myHttpHandlers struct {
	Router *mux.Router
	Handlers []http.Handler
	FsHandler *console.Static
}
func (self *myHttpHandlers)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var matcher mux.RouteMatch
	if self.Router.Match(r, &matcher) {
		self.Router.ServeHTTP(w, r)
	} else {
		self.FsHandler.ServeHTTP(w, r, func(http.ResponseWriter, *http.Request) {

		})
	}
}


func main() {
	flag.Parse()

	_myHttpHandlers := &myHttpHandlers{}
	_myHttpHandlers.Router = console.ExecRouter(*disable, &console.ReqCreateExec{
		Name: *name,
		Host: *host,
		CId:  *cid,
		Cmd:  *cmd,
	})
	_myHttpHandlers.FsHandler = console.NewStatic(console.NewFileSystem())

	http.ListenAndServe(fmt.Sprintf("%v:%v", *ip, *port), _myHttpHandlers)
}

