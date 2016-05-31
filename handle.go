package godtrace

import (
	"errors"
	"unsafe"
)

/*
#cgo LDFLAGS: -ldtrace
#include <stdlib.h>
#include "dtrace.h"
extern void freeArray(int n, char** chs);
extern void freeString(char* chs);
*/
import "C"

type Handle struct {
	handle *C.struct_dtrace_hdl

	consumer consumer

	probe func(*ProbeData) int
	rec   func(*ProbeData, *RecDesc) int
}

func Open(flags int) (*Handle, error) {
	var errno C.int
	cflags := C.int(flags)
	hdl := C.dtrace_open(C.DTRACE_VERSION, cflags, &errno)
	if hdl == nil {
		return nil, errors.New(ErrMsg(nil, int(errno)))
	}
	return newHandle(hdl), nil
}

func newHandle(handle *C.struct_dtrace_hdl) *Handle {
	hdl := &Handle{
		handle: handle,
	}
	return hdl
}

func (h *Handle) Close() {
	C.dtrace_close(h.handle)
}

type Prog struct {
	prog *C.dtrace_prog_t
}

type ProbeSpec C.dtrace_probespec_t

const (
	ProbeSpecNone     ProbeSpec = C.DTRACE_PROBESPEC_NONE
	ProbeSpecProvider ProbeSpec = C.DTRACE_PROBESPEC_PROVIDER
	ProbeSpecMod      ProbeSpec = C.DTRACE_PROBESPEC_MOD
	ProbeSpecFunc     ProbeSpec = C.DTRACE_PROBESPEC_FUNC
	ProbeSpecName     ProbeSpec = C.DTRACE_PROBESPEC_NAME
)

// DTRACE_C_PSPEC

const (
	C_PSPEC int = C.DTRACE_C_PSPEC
)

// dtrace_prog_t *dtrace_program_strcompile (dtrace_hdl_t *handle,
//    char *source, dtrace_probespec_t c_spec, int cflags,
//    int argc, char **argv);
func (h *Handle) Compile(source string, spec ProbeSpec, cflags int, args []string) (*Prog, error) {
	argc, argv := CArray(args)   // free
	csource := C.CString(source) // free
	cspec := C.dtrace_probespec_t(spec)
	ccflags := C.uint_t(cflags)
	prog := C.dtrace_program_strcompile(h.handle, csource, cspec, ccflags, argc, argv)
	if prog == nil {
		return nil, h.GetError()
	}
	C.freeArray(argc, argv)
	C.freeString(csource)
	return &Prog{prog}, nil
}

func (h *Handle) Go() error {
	success := C.dtrace_go(h.handle)
	if success == -1 {
		return h.GetError()
	}
	return nil
}

type ProgInfo C.dtrace_proginfo_t

func (p *ProgInfo) Matches() int {
	return int(p.dpi_matches)
}

func (h *Handle) Exec(prog *Prog) (*ProgInfo, error) {
	var info C.dtrace_proginfo_t
	retno := C.dtrace_program_exec(h.handle, prog.prog, &info)
	if retno == -1 {
		return nil, h.GetError()
	}
	return (*ProgInfo)(&info), nil
}

type ProcHandle C.struct_ps_prochandle

func (h *Handle) ProcGrab(pid int) (*ProcHandle, error) {
	hdl := C.dtrace_proc_grab(h.handle, C.pid_t(pid), 0)
	if hdl == nil {
		return nil, h.GetError()
	}
	return (*ProcHandle)(hdl), nil
}

func (h *Handle) Errno() int {
	return int(C.dtrace_errno(h.handle))
}

func (h *Handle) SetBufSize(value string) {
	h.SetOpt("bufsize", value)
}

func (h *Handle) SetOpt(key, value string) {
	ckey := C.CString(key)     // free
	cvalue := C.CString(value) // free
	C.dtrace_setopt(h.handle, ckey, cvalue)
	C.freeString(ckey)
	C.freeString(cvalue)
}

func (h *Handle) Stop() {
	C.dtrace_stop(h.handle)
}

func (h *Handle) Sleep() {
	C.dtrace_sleep(h.handle)
}

func (h *Handle) GetError() error {
	return errors.New(ErrMsg(h, h.Errno()))
}

func ErrMsg(hdl *Handle, errno int) string {
	var chdl *C.struct_dtrace_hdl
	if hdl != nil {
		chdl = hdl.handle
	}
	return C.GoString(C.dtrace_errmsg(chdl, C.int(errno)))
}

// try freeArray
func CArray(s []string) (C.int, **C.char) {
	if len(s) == 0 {
		return C.int(0), nil
	}
	char := make([]*C.char, len(s))
	for idx, item := range s {
		char[idx] = C.CString(item) // freeArray
	}
	return C.int(len(s)), (**C.char)(unsafe.Pointer(&char[0]))
}
