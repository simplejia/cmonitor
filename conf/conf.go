package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/simplejia/cmonitor/comm"
	"github.com/simplejia/utils"
)

type Conf struct {
	Port     int
	RootPath string
	Environ  string
	Svrs     map[string]string
	Clog     *struct {
		Name  string
		Addr  string
		Mode  int
		Level int
	}
}

var (
	Envs                                       map[string]*Conf
	Env                                        string
	C                                          *Conf
	Start, Stop, Restart, GraceRestart, Status string
)

func init() {
	flag.StringVar(&Start, comm.START, "", "start a svr")
	flag.StringVar(&Stop, comm.STOP, "", "stop a svr")
	flag.StringVar(&Restart, comm.RESTART, "", "restart a svr")
	flag.StringVar(&GraceRestart, comm.GRESTART, "", "grace restart a svr")
	flag.StringVar(&Status, comm.STATUS, "", "status a svr")
	flag.StringVar(&Env, "env", "prod", "set env")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Another process monitor\n")
		fmt.Fprintf(os.Stderr, "version: 1.7, Created by simplejia [12/2014]\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	fcontent, err := ioutil.ReadFile(filepath.Join(dir, "conf", "conf.json"))
	if err != nil {
		println("conf.json not found")
		os.Exit(-1)
	}

	fcontent = utils.RemoveAnnotation(fcontent)
	if err := json.Unmarshal(fcontent, &Envs); err != nil {
		fmt.Println("conf.json wrong format:", err)
		os.Exit(-1)
	}

	C = Envs[Env]
	if C == nil {
		fmt.Println("env not right:", Env)
		os.Exit(-1)
	}

	fmt.Printf("Env: %s\nC: %s\n", Env, utils.Iprint(C))
	return
}
