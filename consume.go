package godtrace

// #include "dtrace.h"
// extern int bufhandler(dtrace_bufdata_t *bufdata, void *arg);
import "C"

import (
	"io"
	"unsafe"
)

type ConsumeMode int

const (
	ConsumeModeNone ConsumeMode = iota
	ConsumeModeFunc
	ConsumeModeChan
	ConsumeModePipe
	ConsumeModeFile
)

type consumer struct {
	mode ConsumeMode

	// ConsumeModeHandler
	handleFunc func(*BufData) int

	// ConsumeModeChan
	handleChan chan *BufData

	// ConsumeModeBuffer
	handlePipe *io.PipeWriter
}

func (c *consumer) Inited() bool {
	return c.mode != ConsumeModeNone
}

func (c *consumer) ConsumeChan() <-chan *BufData {
	c.mode = ConsumeModeChan
	c.handleChan = make(chan *BufData, 32)
	return c.handleChan
}

func (c *consumer) ConsumeFunc(f func(*BufData) int) {
	c.mode = ConsumeModeFunc
	c.handleFunc = f
}

func (c *consumer) ConsumePipe() *io.PipeReader {
	c.mode = ConsumeModePipe
	pr, pw := io.Pipe()
	c.handlePipe = pw
	return pr
}

func (c *consumer) Consume(buf *BufData) int {
	switch c.mode {
	case ConsumeModeFunc:
		return c.handleFunc(buf)
	case ConsumeModeChan:
		c.handleChan <- buf
	case ConsumeModePipe:
		c.handlePipe.Write([]byte(buf.Buffered()))
	case ConsumeModeFile:
		panic("not supported yet")
	}
	return ConsumeThis
}

// -----------------------------------------------------------------------------

//export goBufHandler
func goBufHandler(bufdata *C.dtrace_bufdata_t, arg unsafe.Pointer) C.int {
	return C.int((*Handle)(arg).consumer.Consume((*BufData)(bufdata)))
}

func (h *Handle) ConsumeChan() (<-chan *BufData, error) {
	if !h.consumer.Inited() {
		if err := h.initHandleBuffered(); err != nil {
			return nil, err
		}
	}

	return h.consumer.ConsumeChan(), nil
}

func (h *Handle) ConsumeFunc(f func(*BufData) int) error {
	if !h.consumer.Inited() {
		if err := h.initHandleBuffered(); err != nil {
			return err
		}
	}
	h.consumer.ConsumeFunc(f)
	return nil
}

func (h *Handle) ConsumePipe() (*io.PipeReader, error) {
	if !h.consumer.Inited() {
		if err := h.initHandleBuffered(); err != nil {
			return nil, err
		}
	}
	pr := h.consumer.ConsumePipe()
	return pr, nil
}

func (h *Handle) initHandleBuffered() error {
	handler := (*C.dtrace_handle_buffered_f)(C.bufhandler)
	ret := C.dtrace_handle_buffered(h.handle, handler, unsafe.Pointer(h))
	if ret == -1 {
		return h.GetError()
	}
	return nil
}
