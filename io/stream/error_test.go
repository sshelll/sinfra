package stream

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	ep := NewErrorPasser()

	go func() {
		time.Sleep(time.Second)
		ep.Put(errors.New("an err after 1 sec"))
		ep.Close()
	}()

	err := ep.Get()
	assert.NotNil(t, err, "should got an error from ErrorPasser after 1 sec")
	t.Logf("got an error from ErrorPasser: %v", err)

	err = ep.Get()
	assert.Nil(t, err, "should not get any errors from ErrorPasser because it was closed")
	t.Logf("get error from ErrorPasser again, result: err = %v", err)

}

func TestCheck(t *testing.T) {

	ep := NewErrorPasser()

	go func() {
		time.Sleep(time.Second)
		ep.Put(errors.New("an error after 1 sec"))
		ep.Close()
	}()

	i := 0
	for ; ; i++ {
		err, done := ep.Check()
		t.Logf("check result of the %d time: err = %v, done = %v", i, err, done)
		time.Sleep(time.Millisecond * 200)
		if done {
			t.Log("check done")
			break
		}
	}

	err, done := ep.Check()
	assert.Nil(t, err, "should not get any errors after 'check done'")
	assert.True(t, done, "should always return 'done' after 'check done'")
	t.Logf("check again after 'check done', result: err = %v, done = %v", err, done)

}

func TestCheckAndClose(t *testing.T) {

	ep := NewErrorPasser()

	go func() {
		time.Sleep(time.Second)
		ep.Close()
	}()

	i := 0
	for ; ; i++ {
		err, done := ep.Check()
		t.Logf("check result of the %d time: err = %v, done = %v", i, err, done)
		time.Sleep(time.Millisecond * 200)
		if done {
			t.Log("check done")
			break
		}
	}

	err, done := ep.Check()
	assert.Nil(t, err, "should not get any errors after 'check done'")
	assert.True(t, done, "should always return 'done' after 'check done'")
	t.Logf("check again after 'check done', result: err = %v, done = %v", err, done)

}
