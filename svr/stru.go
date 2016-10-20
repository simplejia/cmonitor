package svr

type Msg struct {
	Command string
}

var ProcChs = make(map[string]chan *Msg)
