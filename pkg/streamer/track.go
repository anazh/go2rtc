package streamer

import (
	"fmt"
	"github.com/pion/rtp"
	"sync"
)

type WriterFunc func(packet *rtp.Packet) error
type WrapperFunc func(push WriterFunc) WriterFunc

type Track struct {
	Codec     *Codec
	Direction string
	sink      map[*Track]WriterFunc
	sinkMu    sync.RWMutex
}

func (t *Track) String() string {
	s := t.Codec.String()
	s += fmt.Sprintf(", sinks=%d", len(t.sink))
	return s
}

func (t *Track) WriteRTP(p *rtp.Packet) error {
	t.sinkMu.RLock()
	for _, f := range t.sink {
		_ = f(p)
	}
	t.sinkMu.RUnlock()
	return nil
}

func (t *Track) Bind(w WriterFunc) *Track {
	t.sinkMu.Lock()

	if t.sink == nil {
		t.sink = map[*Track]WriterFunc{}
	}

	clone := &Track{
		Codec: t.Codec, Direction: t.Direction, sink: t.sink,
	}
	t.sink[clone] = w

	t.sinkMu.Unlock()

	return clone
}

func (t *Track) Unbind() {
	t.sinkMu.Lock()
	delete(t.sink, t)
	t.sinkMu.Unlock()
}

func (t *Track) GetSink(from *Track) {
	t.sink = from.sink
}

func (t *Track) HasSink() bool {
	t.sinkMu.RLock()
	defer t.sinkMu.RUnlock()
	return len(t.sink) > 0
}
