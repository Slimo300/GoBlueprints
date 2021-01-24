package trace

import (
	"fmt"
	"io"
)

// Tracer interface
type Tracer interface {
	Trace(...interface{})
}

// New func returns a new Tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off func creates nilTracer
func Off() Tracer {
	return &nilTracer{}
}
