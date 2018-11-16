// 进程监控服务.
// author: simplejia
// date: 2014/12/01
package main

import (
	"fmt"
	"time"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/cmonitor/comm"
	"github.com/simplejia/cmonitor/conf"
	"github.com/simplejia/cmonitor/svr"
	"github.com/simplejia/namecli/api"
	"github.com/simplejia/utils"
)

func request(command string, service string) {
	url := fmt.Sprintf("http://%s:%d", utils.LocalIp, conf.C.Port)
	params := map[string]string{
		"command": command,
		"service": service,
	}
	gpp := &utils.GPP{
		Uri:     url,
		Timeout: time.Second * 8,
		Params:  params,
	}
	body, err := utils.Get(gpp)
	if err != nil {
		fmt.Printf("Error: [cmonitor maybe down!] %v, %s\n", err, body)
		return
	}

	fmt.Println(string(body))
	return
}

func init() {
	clog.AddrFunc = func() (string, error) {
		return api.Name(conf.C.Clog.Addr)
	}
	clog.Init(conf.C.Clog.Name, "", conf.C.Clog.Level, conf.C.Clog.Mode)
}

func main() {
	switch {
	case conf.Start != "":
		request(comm.START, conf.Start)
	case conf.Stop != "":
		request(comm.STOP, conf.Stop)
	case conf.Restart != "":
		request(comm.RESTART, conf.Restart)
	case conf.GraceRestart != "":
		request(comm.GRESTART, conf.GraceRestart)
	case conf.Status != "":
		request(comm.STATUS, conf.Status)
	default:
		clog.Info("main() StartSvr")
		svr.StartSvr()
	}
}
