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

	"github.com/simplejia/clog"
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

	err := s.ListenAndServe()
	clog.Error("StartHttpSvr() %v", err)
	os.Exit(-1)
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
				w.Write([]byte(fmt.Sprintf("Error: serivce[%s] not configure", service)))
			}
			services = append(services, service)
		}
	} else {
		for service, _ := range conf.C.Svrs {
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
		services_fail := []string{}
		processes_old := []*os.Process{}

		for _, service := range services {
			fullpath := filepath.Join(conf.C.RootPath, conf.C.Svrs[service])
			process_old, _ := procs.GetProc(fullpath)
			processes_old = append(processes_old, process_old)

			select {
			case ProcChs[service] <- &Msg{Command: command}:
			default:
			}
		}

		for pos, service := range services {
			process_old := processes_old[pos]
			step := 5
			for ; step > 0; step-- {
				fullpath := filepath.Join(conf.C.RootPath, conf.C.Svrs[service])
				process_new, err := procs.GetProc(fullpath)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("\nget service %s error: %v\n", service, err)))
					continue
				}

				if command == comm.STOP {
					if process_new == nil {
						break
					}
				} else if command == comm.START {
					if process_new != nil {
						break
					}
				} else if command == comm.RESTART {
					if process_new != nil {
						if process_old == nil {
							break
						} else if process_new.Pid != process_old.Pid {
							break
						}
					}
				}
				time.Sleep(time.Millisecond * 300)
			}

			if step == 0 {
				services_fail = append(services_fail, service)
			}
		}

		if len(services_fail) == 0 {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte(fmt.Sprintf("FAIL: %v", strings.Join(services_fail, ", "))))
		}
	}
}
