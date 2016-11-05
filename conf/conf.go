package conf

import (
	"chatroombase/cmonitor/comm"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/simplejia/utils"
)

type Conf struct {
	Port     int
	RootPath string
	Environ  string
	Svrs     map[string]string
	Clog     *struct {
		Mode  int
		Level int
	}
}

var (
	Envs                         map[string]*Conf
	Env                          string
	C                            *Conf
	Start, Stop, Restart, Status string
)

func init() {
	flag.StringVar(&Start, comm.START, "", "start a svr")
	flag.StringVar(&Stop, comm.STOP, "", "stop a svr")
	flag.StringVar(&Restart, comm.RESTART, "", "restart a svr")
	flag.StringVar(&Status, comm.STATUS, "", "status a svr")
	flag.StringVar(&Env, "env", "prod", "set env")

	var conf string
	flag.StringVar(&conf, "conf", "", "set custom conf")

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

	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("conf not right:", err)
				os.Exit(-1)
			}
		}()
		ccs := strings.Split(conf, "::")
		for _, cs := range ccs {
			pos := strings.Index(cs, "=")
			if pos == -1 {
				continue
			}
			name, value := strings.TrimSpace(cs[:pos]), strings.TrimSpace(cs[pos+1:])

			rv := reflect.Indirect(reflect.ValueOf(C))
			for _, field := range strings.Split(name, ".") {
				rv = reflect.Indirect(rv.FieldByName(strings.Title(field)))
			}
			switch rv.Kind() {
			case reflect.String:
				rv.SetString(value)
			case reflect.Bool:
				b, err := strconv.ParseBool(value)
				if err != nil {
					panic(err)
				}
				rv.SetBool(b)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					panic(err)
				}
				rv.SetInt(i)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				u, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					panic(err)
				}
				rv.SetUint(u)
			case reflect.Float32, reflect.Float64:
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				rv.SetFloat(f)
			}
		}
	}()

	fmt.Printf("Env: %s\nC: %s\n", Env, utils.Iprint(C))

	return
}
