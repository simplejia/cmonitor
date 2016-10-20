// 进程监控服务.
// author: simplejia
// date: 2014/12/01
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/simplejia/clog"
	_ "github.com/simplejia/cmonitor/clog"
	"github.com/simplejia/cmonitor/comm"
	"github.com/simplejia/cmonitor/conf"
	"github.com/simplejia/cmonitor/svr"
	"github.com/simplejia/utils"
)

func showenv() {
	fmt.Printf("env: %s\n", conf.Env)
	for k, v := range conf.C.Svrs {
		if v != "" {
			fmt.Printf("%s\t%s\n", k, v)
		}
	}
}

func request(command string, service string) {
	url := fmt.Sprintf(
		"http://%s:%d/?command=%s&service=%s",
		utils.GetLocalIp(), conf.C.Port, command, service,
	)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %v [cmonitor maybe down!]\n", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if g, e := resp.StatusCode, http.StatusOK; g != e {
		fmt.Printf("Error: %v, check your conf\n", string(content))
		return
	}

	if content != nil {
		fmt.Println(string(content))
	}

	return
}

func main() {
	var env bool
	var start, stop, restart, status string

	flag.BoolVar(&env, comm.ENV, false, "show env info")
	flag.StringVar(&start, comm.START, "", "start a svr")
	flag.StringVar(&stop, comm.STOP, "", "stop a svr")
	flag.StringVar(&restart, comm.RESTART, "", "restart a svr")
	flag.StringVar(&status, comm.STATUS, "", "status a svr")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Another process monitor\n")
		fmt.Fprintf(os.Stderr, "version: 1.7, Created by simplejia [12/2014]\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	switch {
	case env:
		showenv()
	case start != "":
		request(comm.START, start)
	case stop != "":
		request(comm.STOP, stop)
	case restart != "":
		request(comm.RESTART, restart)
	case status != "":
		request(comm.STATUS, status)
	default:
		clog.Info("main() StartSvr")
		svr.StartSvr()
	}
}
