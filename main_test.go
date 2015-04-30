package main

import (
	"testing"

	"github.com/brimstone/go-saverequest"
)

func TestHandleData(t *testing.T) {
	saverequest.TestRequestFiles(t, ".", handleData)
}
