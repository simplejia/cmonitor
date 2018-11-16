package svr

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/cmonitor/comm"
	"github.com/simplejia/cmonitor/conf"
	"github.com/simplejia/cmonitor/procs"
	"github.com/simplejia/utils"
)

func StartHttpSvr() {
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%d", utils.LocalIp, conf.C.Port),
	}

	http.HandleFunc("/", indexHandler)

	if err := s.ListenAndServe(); err != nil {
		clog.Error("StartHttpSvr() %v", err)
		os.Exit(-1)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	command := r.FormValue("command")
	service := r.FormValue("service")

	services := []string{}
	if service != "all" {
		for _, service := range strings.Split(strings.Replace(service, " ", "", -1), ",") {
			if service == "" {
				continue
			}
			if _, ok := conf.C.Svrs[service]; !ok {
				w.Write([]byte(fmt.Sprintf("Error: serivce[%s] not configure\n", service)))
			}
			services = append(services, service)
		}
	} else {
		for service := range conf.C.Svrs {
			if service == "" {
				continue
			}
			services = append(services, service)
		}
	}

	switch command {
	case comm.STATUS:
		statusoks, statusfails := []string{}, []string{}
		for _, service := range services {
			fullpath := filepath.Join(conf.C.RootPath, conf.C.Svrs[service])
			process, err := procs.GetProc(fullpath)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("\nget service %s error: %v\n", service, err)))
			}
			if process != nil {
				statusoks = append(statusoks, service+" PID:"+strconv.Itoa(process.Pid))
			} else {
				statusfails = append(statusfails, service)
			}
		}

		sort.Strings(statusoks)
		w.Write([]byte(fmt.Sprintf("\n*****STATUS OK SERVICE LIST*****\n%s\n", strings.Join(statusoks, "\n"))))

		sort.Strings(statusfails)
		w.Write([]byte(fmt.Sprintf("\n*****STATUS FAIL SERVICE LIST*****\n%s\n", strings.Join(statusfails, "\n"))))
	default:
		failServices := []string{}
		oldProcesses := map[string]*os.Process{}

		for _, service := range services {
			fullpath := filepath.Join(conf.C.RootPath, conf.C.Svrs[service])
			oldProcess, _ := procs.GetProc(fullpath)
			if oldProcess != nil {
				oldProcesses[service] = oldProcess
			}

			select {
			case ProcChs[service] <- &Msg{Command: command}:
			default:
			}
		}

		for _, service := range services {
			oldProcess := oldProcesses[service]
			step := 6
			for ; step > 0; step-- {
				fullpath := filepath.Join(conf.C.RootPath, conf.C.Svrs[service])
				newProcess, err := procs.GetProc(fullpath)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("\nget service %s error: %v\n", service, err)))
					continue
				}

				if command == comm.STOP {
					if newProcess == nil {
						break
					}
				} else if command == comm.START {
					if newProcess != nil {
						break
					}
				} else if command == comm.RESTART || command == comm.GRESTART {
					if newProcess != nil {
						if oldProcess == nil {
							break
						} else if newProcess.Pid != oldProcess.Pid {
							break
						}
					}
				}
				time.Sleep(time.Millisecond * 300)
			}

			if step == 0 {
				failServices = append(failServices, service)
			}
		}

		if len(failServices) == 0 {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte(fmt.Sprintf("FAIL: %v", strings.Join(failServices, ", "))))
		}
	}
}
