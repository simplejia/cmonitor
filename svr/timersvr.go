package svr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/simplejia/clog"
	"github.com/simplejia/cmonitor/comm"
	"github.com/simplejia/cmonitor/conf"
	"github.com/simplejia/cmonitor/procs"
)

func proc(service string, cmd string) {
	defer func() {
		if err := recover(); err != nil {
			clog.Error("proc() recover %v, %s", err, debug.Stack())
			os.Exit(-1)
		}
	}()

	fullpath := filepath.Join(conf.C.RootPath, cmd)

	process, err := procs.GetProc(fullpath)
	if err != nil {
		clog.Error("proc() GetProc %s error: %v", service, err)
		return
	}

	tick1 := time.Tick(time.Millisecond * 300)
	tick2 := time.Tick(time.Minute)
	tick3 := time.Tick(time.Hour * 24)
	failNum := 0
	status := 2 // 1: stop 2: start 3: restart 4: grace restart
	msgCh := ProcChs[service]

	for {
		select {
		case <-tick1:
			switch status {
			case 1: // stop
				if ok := procs.CheckProc(process); ok {
					if failNum++; failNum > 5 {
						clog.Error("proc() stop %s always fail, must check it", service)
						failNum = 0
						time.Sleep(time.Second * 3)
					}
					if err := procs.StopProc(process); err != nil {
						clog.Error("proc() StopProc %s error: %v", service, err)
					} else {
						process = nil
					}
				}
			case 2: // start
				if ok := procs.CheckProc(process); !ok {
					if failNum++; failNum > 5 {
						clog.Error("proc() start %s always fail, must check it", service)
						failNum = 0
						time.Sleep(time.Second * 3)
					}
					if process_i, err := procs.StartProc(fullpath, conf.C.Environ); err != nil || process_i == nil {
						clog.Error("proc() StartProc %s error: %v", service, err)
					} else {
						process = process_i
					}
					time.Sleep(time.Second)
				}
			case 3: // restart
				if ok := procs.CheckProc(process); ok {
					if failNum++; failNum > 5 {
						clog.Error("proc() stop %s always fail, must check it", service)
						failNum = 0
						time.Sleep(time.Second * 3)
					}
					if err := procs.StopProc(process); err != nil {
						clog.Error("proc() StopProc %s error: %v", service, err)
					} else {
						process = nil
					}
				} else {
					status = 2
				}
			case 4: // grace restart
				if ok := procs.CheckProc(process); ok {
					if failNum++; failNum > 5 {
						clog.Error("proc() stop %s always fail, must check it", service)
						failNum = 0
						time.Sleep(time.Second * 3)
					}
					if err := procs.GStopProc(process); err != nil {
						clog.Error("proc() GStopProc %s error: %v", service, err)
					} else {
						process = nil
					}
				} else {
					status = 2
				}
			}
		case <-tick2:
			if status == 2 {
				if process, err := procs.GetProc(fullpath); err != nil || process == nil {
					clog.Error("proc() GetProc %s error: %v", service, err)
				}
			}
		case <-tick3:
			dirname := ""
			pos := strings.Index(fullpath, " ")
			if pos != -1 {
				dirname = filepath.Dir(fullpath[:pos])
			} else {
				dirname = filepath.Dir(fullpath)
			}
			cmdStr := fmt.Sprintf(
				"cd %s; head -c`wc -c cmonitor.log|awk '{print $1}'` cmonitor.log >cmonitor.%d.log && cat /dev/null >cmonitor.log",
				dirname, time.Now().Day(),
			)
			err := exec.Command("sh", "-c", cmdStr).Run()
			if err != nil {
				clog.Error("proc() exec.Command error: %v", err)
			}
		case msg := <-msgCh:
			switch msg.Command {
			case comm.STOP:
				failNum, status = 0, 1
			case comm.START:
				failNum, status = 0, 2
			case comm.RESTART:
				failNum, status = 0, 3
			case comm.GRESTART:
				failNum, status = 0, 4
			default:
				clog.Error("proc() unexpected command: %s", msg.Command)
			}
		}
	}
}

func StartCronSvr() {
	for k, v := range conf.C.Svrs {
		if k != "" && v != "" {
			go proc(k, v)
		}
	}
}
