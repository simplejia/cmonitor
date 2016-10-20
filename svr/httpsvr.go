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
		Addr:        fmt.Sprintf("%s:%d", utils.GetLocalIp(), conf.C.Port),
		ReadTimeout: time.Second * 3,
	}

	http.HandleFunc("/", indexHandler)

	err := s.ListenAndServe()
	clog.Error("StartHttpSvr() %v", err)
	os.Exit(-1)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		}
	}()
	command := r.FormValue("command")
	service := r.FormValue("service")
	if command == "" || service == "" {
		panic("command or serivce is empty")
	}

	services := []string{}
	if service != "all" {
		if _, ok := conf.C.Svrs[service]; !ok {
			panic("serivce not configure")
		}
		services = append(services, service)
	} else {
		for _service, _ := range conf.C.Svrs {
			if _service != "" {
				services = append(services, _service)
			}
		}
	}

	switch command {
	case comm.STATUS:
		statusoks, statusfails := []string{}, []string{}
		for _, _service := range services {
			fullpath := filepath.Join(conf.C.RootPath, conf.C.Svrs[_service])
			process, err := procs.GetProc(fullpath)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("\nget service %s error: %v\n", _service, err)))
			}
			if process != nil {
				statusoks = append(statusoks, _service+" PID:"+strconv.Itoa(process.Pid))
			} else {
				statusfails = append(statusfails, _service)
			}
		}

		sort.Strings(statusoks)
		w.Write([]byte(fmt.Sprintf("\n*****STATUS OK SERVICE LIST*****\n%s\n", strings.Join(statusoks, "\n"))))

		sort.Strings(statusfails)
		w.Write([]byte(fmt.Sprintf("\n*****STATUS FAIL SERVICE LIST*****\n%s\n", strings.Join(statusfails, "\n"))))
	default:
		services_fail := []string{}
		var process_old *os.Process
		if len(services) == 1 {
			process_old, _ = procs.GetProc(filepath.Join(conf.C.RootPath, conf.C.Svrs[services[0]]))
		}

		for _, _service := range services {
			select {
			case ProcChs[_service] <- &Msg{Command: command}:
			}
		}

		if len(services) == 1 {
			service := services[0]
			i := 5
			for ; i > 0; i-- {
				process_new, err := procs.GetProc(filepath.Join(conf.C.RootPath, conf.C.Svrs[service]))
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

			if i == 0 {
				services_fail = append(services_fail, service)
			}
		}

		if len(services_fail) == 0 {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte(fmt.Sprintf("FAIL: %v", services_fail)))
		}
	}
}
