package procs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func GetProc(cmd string) (process *os.Process, err error) {
	output, err := exec.Command("ps", "-e", "-opid", "-oppid", "-ocommand").CombinedOutput()
	if err != nil {
		err = fmt.Errorf("err: %v, output: %s", err, output)
		return
	}

	lines := strings.Split(string(output), "\n")
	pid := ""
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		_pid, _ppid, _cmd := fields[0], fields[1], strings.Join(fields[2:], " ")

		if strings.Join(strings.Fields(cmd), " ") != _cmd {
			continue
		}
		if pid == "" {
			if _ppid == "1" {
				pid = _pid
			} else {
				pid = _ppid
			}
		} else {
			if _ppid != pid {
				err = fmt.Errorf("GetProc() %s multi process exist, ppid:%v pid:%v", cmd, _ppid, pid)
				return
			}
		}
	}

	if pid == "" {
		return
	}

	i, err := strconv.Atoi(pid)
	if err != nil {
		return
	}
	process, err = os.FindProcess(i)
	return
}

func StopProc(process *os.Process) (err error) {
	if process == nil {
		return
	}
	if err = process.Kill(); err != nil {
		return
	}

	process.Release()
	return
}

func CheckProc(process *os.Process) (ok bool, err error) {
	if process == nil {
		return
	}
	_err := process.Signal(syscall.Signal(0))
	if _err == nil {
		ok = true
	} else if _err == syscall.ESRCH || _err.Error() == "os: process already finished" {
		ok = false
	} else {
		err = _err
	}
	return
}

func StartProc(cmd string, env string) (process *os.Process, err error) {
	if process, err = GetProc(cmd); err != nil || process != nil {
		return
	}

	dirname := ""
	pos := strings.Index(cmd, " ")
	if pos != -1 {
		dirname = filepath.Dir(cmd[:pos])
	} else {
		dirname = filepath.Dir(cmd)
	}

	env = strings.Trim(env, ";")
	if env == "" {
		env = "true"
	}
	cmdStr := fmt.Sprintf(
		"cd %s; %s; nohup %s >cmonitor.log 2>&1 &",
		dirname, env, cmd,
	)

	err = exec.Command("sh", "-c", cmdStr).Run()
	if err != nil {
		return
	}
	process, err = GetProc(cmd)
	if err != nil {
		return
	}
	if process != nil {
		return
	}
	content, err := ioutil.ReadFile(filepath.Join(dirname, "cmonitor.log"))
	if err != nil {
		return
	}
	err = errors.New(string(content))
	return
}
