package recorder

import (
	"bytes"
	_ "context"
	"io"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	_ = New(logger, 1)
}

func Test_AddSubscriber(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	rec := New(logger, 1)
	ch := make(chan string)

	assert.Nil(t, rec.AddSubscriber(ch))
	assert.NotNil(t, rec.AddSubscriber(ch))
}

func Test_RemoveSubscriber(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	rec := New(logger, 1)
	ch := make(chan string)
	rec.AddSubscriber(ch)

	assert.Nil(t, rec.RemoveSubscriber(ch))
	assert.NotNil(t, rec.RemoveSubscriber(ch))
}

func Test_BroadcastToSubscribers(t *testing.T) {
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	rec := New(logger, 1)

	ch1 := make(chan string)
	ch2 := make(chan string)
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)

	writerFunc := func(wg *sync.WaitGroup, ch chan string, buf io.Writer) {
		defer wg.Done()
		for in := range ch {
			buf.Write([]byte(in))
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go writerFunc(wg, ch1, buf1)
	go writerFunc(wg, ch2, buf2)

	assert.Nil(t, rec.AddSubscriber(ch1))
	assert.Nil(t, rec.AddSubscriber(ch2))

	rec.BroadcastToSubscribers("foo")
	rec.BroadcastToSubscribers("bar")
	rec.BroadcastToSubscribers("buz")

	assert.Nil(t, rec.RemoveSubscriber(ch1))
	assert.Nil(t, rec.RemoveSubscriber(ch2))

	close(ch1)
	close(ch2)
	wg.Wait()

	assert.Equal(t, "foobarbuz", buf1.String())
	assert.Equal(t, "foobarbuz", buf2.String())
}

//func Test_Start(t *testing.T) {
//	logbuf := new(bytes.Buffer)
//	logger := log.New(logbuf, "", 0)
//	rec := New(logger)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	in := rec.Start(ctx)
//
//	cancel()
//}
