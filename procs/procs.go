package procs

import (
	"bytes"
	"errors"
	"fmt"
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

	pid := ""

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		_pid, _ppid, _cmd := fields[0], fields[1], strings.Join(fields[2:], " ")

		if strings.Join(strings.Fields(cmd), " ") != _cmd {
			continue
		}

		var tpid string
		if _ppid == "1" {
			tpid = _pid
		} else {
			tpid = _ppid
		}

		if pid == "" {
			pid = tpid
			continue
		}

		if tpid != pid {
			err = fmt.Errorf("GetProc() %s multi process exist, pids: %v,%v", cmd, tpid, pid)
			return
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

func GStopProc(process *os.Process) (err error) {
	if process == nil {
		return
	}
	// SIGHUP: 1
	if err = process.Signal(syscall.Signal(1)); err != nil {
		return
	}
	process.Release()
	return
}

func CheckProc(process *os.Process) (ok bool) {
	if process == nil {
		return
	}
	err := process.Signal(syscall.Signal(0))
	if err == nil {
		ok = true
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
		"cd %s && %s && $(nohup %s >>cmonitor.log 2>&1 &)",
		dirname, env, cmd,
	)

	command := exec.Command("sh", "-c", cmdStr)
	errBuf := new(bytes.Buffer)
	command.Stderr = errBuf
	if err = command.Run(); err != nil {
		if errBuf.Len() > 0 {
			err = errors.New(errBuf.String())
		}
		return
	}

	process, err = GetProc(cmd)
	if err != nil {
		return
	}
	if process != nil {
		return
	}

	logfile := filepath.Join(dirname, "cmonitor.log")
	fi, err := os.Stat(logfile)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New("check cmonitor log")
			return
		}
		return
	}

	logpos := int64(0)
	logbuf := make([]byte, 256)
	if n, m := fi.Size(), int64(len(logbuf)); n > m {
		logpos = n - m
	}

	f, err := os.Open(logfile)
	if err != nil {
		return
	}
	defer f.Close()

	f.ReadAt(logbuf, logpos)

	err = errors.New(string(logbuf))
	return
}
