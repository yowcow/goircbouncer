package recorder

import (
	_ "context"
	"fmt"
	"log"

	_ "github.com/yowcow/goircparser"
)

type subscribers map[chan<- string]bool

type Recorder struct {
	logger      *log.Logger
	subscribers subscribers
}

func New(logger *log.Logger, size int) *Recorder {
	subscribers := make(subscribers)
	return &Recorder{logger, subscribers}
}

func (r *Recorder) AddSubscriber(out chan<- string) error {
	if _, ok := r.subscribers[out]; ok {
		return fmt.Errorf("already subscribing")
	}
	r.subscribers[out] = true
	return nil
}

func (r *Recorder) RemoveSubscriber(out chan<- string) error {
	if _, ok := r.subscribers[out]; ok {
		delete(r.subscribers, out)
		return nil
	}
	return fmt.Errorf("not subscribing")
}

func (r Recorder) BroadcastToSubscribers(command string) {
	for out, _ := range r.subscribers {
		out <- command
	}
}

//func (r Recorder) Start(ctx context.Context) chan<- *parser.Row {
//	in := make(chan *parser.Row)
//	go r.worker(ctx, in)
//	return in
//}
//
//func (r Recorder) worker(ctx context.Context, in <-chan *parser.Row) {
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case row := <-in:
//			r.BroadcastToSubscribers(row.RawLine)
//		}
//	}
//}
