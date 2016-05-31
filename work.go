package godtrace

/*
#include "dtrace.h"
extern int chew(dtrace_probedata_t *, void *);
extern int chewrec(dtrace_probedata_t *, dtrace_recdesc_t *, void *);
extern int dumpChewrec(dtrace_probedata_t *, dtrace_recdesc_t *, void *);
*/
import "C"
import "unsafe"

//export goChewRec
func goChewRec(data *C.dtrace_probedata_t, rec *C.dtrace_recdesc_t, arg unsafe.Pointer) C.int {
	return C.int((*Handle)(arg).rec((*ProbeData)(data), (*RecDesc)(rec)))
}

//export goChew
func goChew(data *C.dtrace_probedata_t, arg unsafe.Pointer) C.int {
	return C.int((*Handle)(arg).probe((*ProbeData)(data)))
}

func (h *Handle) SetHandlerFunc(f func(*ProbeData) int) {
	h.probe = f
}
func (h *Handle) SetRecHandlerFunc(f func(*ProbeData, *RecDesc) int) {
	h.rec = f
}

func (h *Handle) Work() WorkStatus {
	var p *C.dtrace_consume_probe_f
	if h.probe != nil {
		p = (*C.dtrace_consume_probe_f)(C.chew)
	}
	var r *C.dtrace_consume_rec_f
	if h.rec != nil {
		r = (*C.dtrace_consume_rec_f)(C.chewrec)
	} else {
		r = (*C.dtrace_consume_rec_f)(C.dumpChewrec)
	}
	status := C.dtrace_work(h.handle, C.stdout, p, r, unsafe.Pointer(h))
	return WorkStatus(status)
}

type WorkStatus int

const (
	WS_ERROR WorkStatus = C.DTRACE_WORKSTATUS_ERROR
	WS_OKAY  WorkStatus = C.DTRACE_WORKSTATUS_OKAY
	WS_DONE  WorkStatus = C.DTRACE_WORKSTATUS_DONE
)

func (w WorkStatus) String() string {
	switch w {
	case WS_ERROR:
		return "error"
	case WS_OKAY:
		return "okay"
	case WS_DONE:
		return "done"
	default:
		return "unknown"
	}
}

const (
	ConsumeError = C.DTRACE_CONSUME_ERROR
	ConsumeThis  = C.DTRACE_CONSUME_THIS
	ConsumeNext  = C.DTRACE_CONSUME_NEXT
	ConsumeAbort = C.DTRACE_CONSUME_ABORT
)

// -----------------------------------------------------------------------------
// typedef struct dtrace_recdesc {
//     dtrace_actkind_t dtrd_action;           /* kind of action */
//     uint32_t dtrd_size;                     /* size of record */
//     uint32_t dtrd_offset;                   /* offset in ECB's data */
//     uint16_t dtrd_alignment;                /* required alignment */
//     uint16_t dtrd_format;                   /* format, if any */
//     uint64_t dtrd_arg;                      /* action argument */
//     uint64_t dtrd_uarg;                     /* user argument */
// } dtrace_recdesc_t;

type RecDesc C.dtrace_recdesc_t

// -----------------------------------------------------------------------------
// typedef struct dtrace_probedata {
//     dtrace_hdl_t *dtpda_handle;          /* handle to DTrace library */
//     dtrace_eprobedesc_t *dtpda_edesc;    /* enabled probe description */
//     dtrace_probedesc_t *dtpda_pdesc;     /* probe description */
//     processorid_t dtpda_cpu;             /* CPU for data */
//     caddr_t dtpda_data;                  /* pointer to raw data */
//     dtrace_flowkind_t dtpda_flow;        /* flow kind */
//     const char *dtpda_prefix;            /* recommended flow prefix */
//     int dtpda_indent;                    /* recommended flow indent */
// } dtrace_probedata_t;

type ProbeData C.dtrace_probedata_t

func (p *ProbeData) CPU() int {
	return int(p.dtpda_cpu)
}

type FlowKind int

const (
	FLOW_ENTRY  FlowKind = C.DTRACEFLOW_ENTRY
	FLOW_RETURN FlowKind = C.DTRACEFLOW_RETURN
	FLOW_NONE   FlowKind = C.DTRACEFLOW_NONE
)

func (f FlowKind) String() string {
	switch f {
	case FLOW_ENTRY:
		return "entry"
	case FLOW_RETURN:
		return "return"
	case FLOW_NONE:
		return "none"
	default:
		return "unknown"
	}
}

func (p *ProbeData) Flow() FlowKind {
	return FlowKind(p.dtpda_flow)
}

func (p *ProbeData) PDesc() *ProbeDesc {
	return (*ProbeDesc)(p.dtpda_pdesc)
}

func (p *ProbeData) Indent() int {
	return int(p.dtpda_indent)
}

func (p *ProbeData) Prefix() string {
	return C.GoString(p.dtpda_prefix)
}

// -----------------------------------------------------------------------------
// typedef struct dtrace_probedesc {
//     dtrace_id_t dtpd_id;                    /* probe identifier */
//     char dtpd_provider[DTRACE_PROVNAMELEN]; /* probe provider name */
//     char dtpd_mod[DTRACE_MODNAMELEN];       /* probe module name */
//     char dtpd_func[DTRACE_FUNCNAMELEN];     /* probe function name */
//     char dtpd_name[DTRACE_NAMELEN];         /* probe name */
// } dtrace_probedesc_t;

type ProbeDesc C.dtrace_probedesc_t

func (p *ProbeDesc) Id() int {
	return int(p.dtpd_id)
}

func (p *ProbeDesc) Provider() string {
	return C.GoString(&p.dtpd_provider[0])
}

func (p *ProbeDesc) Mod() string {
	return C.GoString(&p.dtpd_mod[0])
}

func (p *ProbeDesc) Func() string {
	return C.GoString(&p.dtpd_func[0])
}
func (p *ProbeDesc) Name() string {
	return C.GoString(&p.dtpd_name[0])
}
