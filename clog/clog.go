package clog

import (
	"github.com/simplejia/clog"
	"github.com/simplejia/cmonitor/conf"
)

func init() {
	clog.Init("cmonitor", "", conf.C.Clog.Level, conf.C.Clog.Mode)
}
