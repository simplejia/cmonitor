package clog

import (
	"github.com/simplejia/clog"
	"github.com/simplejia/cmonitor/conf"
	"github.com/simplejia/namecli/api"
)

func init() {
	clog.AddrFunc = func() (string, error) {
		return api.Name(conf.C.Clog.Addr)
	}
	clog.Init(conf.C.Clog.Name, "", conf.C.Clog.Level, conf.C.Clog.Mode)
}
