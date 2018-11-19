// +build !windows

package console

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"unsafe"

	"github.com/kr/pty"
)

type SessionsShell struct {
	sessions map[string]*os.File
}

var _ = (*SessionsShell)(nil)

func NewShellSessions(host string) (Sessions, error) {
	return &SessionsShell{
		sessions: map[string]*os.File{},
	}, nil
}

func (d *SessionsShell) CreateExec(req *ReqCreateExec) (*RespCreateExec, error) {
	sh := exec.Command(req.Cmd)
	id := "0x" + strconv.FormatUint(uint64(uintptr(unsafe.Pointer(sh))), 16)
	// Start the command with a pty.
	ptmx, err := pty.Start(sh)
	if err != nil {
		return nil, err
	}

	d.sessions[id] = ptmx
	return &RespCreateExec{
		EId: id,
	}, nil
}

func (d *SessionsShell) StartExec(eid string, ws io.ReadWriter) error {
	cli, ok := d.sessions[eid]
	if !ok {
		return fmt.Errorf("Can not find eid " + eid)
	}
	defer func() {
		delete(d.sessions, eid)
		cli.Close()
	}()

	go io.Copy(cli, ws)
	io.Copy(ws, cli)
	return nil
}

func (d *SessionsShell) ResizeExecTTY(req *ReqResizeExecTTY) error {
	return nil
}
