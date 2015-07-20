package test

import (
	"io"
)

type mockIo struct {
}

func (*mockIo) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (*mockIo) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (*mockIo) Close() error {
	return nil
}

func NewMockIo() io.ReadWriteCloser {
	return new(mockIo)
}
