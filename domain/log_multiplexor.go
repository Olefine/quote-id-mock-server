package domain

import (
	"io"
	"os"
)

// HTTP2Writer to forward Writes into connection
type HTTP2Writer struct {
	Out chan interface{}
}

func (hw HTTP2Writer) Write(p []byte) (int, error) {
	hw.Out <- p
	return len(p), nil
}

// LogMultiplexor ...
type LogMultiplexor struct {
	Outputs []io.Writer
}

func (lm LogMultiplexor) Write(p []byte) (int, error) {
	for _, v := range lm.Outputs {
		go func(wrt io.Writer) {
			wrt.Write(p)
		}(v)
	}

	return len(p), nil
}

// NewLogMultiplexor constucts LogMultiplexor
func NewLogMultiplexor(broadcastTo io.Writer) *LogMultiplexor {
	outputs := []io.Writer{os.Stdout, broadcastTo}
	return &LogMultiplexor{Outputs: outputs}
}
