package clog

import (
	"github.com/simplejia/clog"
	"github.com/simplejia/cmonitor/conf"
)

func init() {
	clog.Init("cmonitor", "", conf.C.Log.Level, conf.C.Log.Mode)
}
