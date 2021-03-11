package speed

import (
	"context"
	"golang.org/x/time/rate"
	"time"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

type bucket struct {
	writer  buf.Writer
	limiter *rate.Limiter
}

// RateWriter bucket with rate limit
func RateWriter(writer buf.Writer, limiter *rate.Limiter) buf.Writer {
	return &bucket{
		writer:  writer,
		limiter: limiter,
	}
}

// WriteMultiBuffer writes a MultiBuffer into underlying writer.
func (w *bucket) WriteMultiBuffer(mb buf.MultiBuffer) error {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
	for i := 0; i < 10; i++ {
		err := w.limiter.WaitN(ctx, int(mb.Len())/4)
		if err != nil {
			err = newError("waiting to get a new ticket").AtDebug()
			// close when waiting 1s
			if i == 10 {
				w.Close()
				return err
			}
		} else {
			break
		}
	}

	return w.writer.WriteMultiBuffer(mb)
}

// Close WriteBuffer
func (w *bucket) Close() error {
	return common.Close(w.writer)
}
