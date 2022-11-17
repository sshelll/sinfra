package stream

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {

	stream := NewIOStream()
	ctx := context.Background()

	go func() {
		time.Sleep(time.Second)
		stream.Write(NewSimpleDatapack(context.WithValue(ctx, "key", "extra"), nil))
	}()

	// block read
	data, closed := stream.Read()
	assert.Equal(t, data.Context().Value("key").(string), "extra")
	assert.False(t, closed)

	// close and read
	stream.Close()
	data, closed = stream.Read()
	assert.Nil(t, data)
	assert.True(t, closed)

	// closed again
	stream.Close()
	data, closed = stream.Read()
	assert.Nil(t, data)
	assert.True(t, closed)

}

func TestTryRead(t *testing.T) {

	stream := NewIOStream()

	go func() {
		time.Sleep(time.Second)
		stream.Write(NewSimpleDatapack(context.Background(), nil))
		time.Sleep(time.Second)
		stream.Write(NewSimpleDatapack(context.Background(), nil))
		stream.Write(NewSimpleDatapack(context.Background(), nil))
		stream.Close()
	}()

	for i := 0; ; i++ {
		data, closed := stream.TryRead()
		t.Logf("try read for the %d time, result: data = %v, closed = %v", i, data, closed)
		if closed {
			t.Log("stream closed, break")
			break
		}
		time.Sleep(time.Millisecond * 200)
	}

	data, closed := stream.TryRead()
	assert.Nil(t, data)
	assert.True(t, closed)
	t.Logf("try read for the last time, result: data is nil = %v, closed = %v", data == nil, closed)

}
