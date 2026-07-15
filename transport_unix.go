//go:build !windows

package main

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

func newTransport(socketPath string) Transport {
	return &unixTransport{addr: socketPath}
}

type unixTransport struct {
	addr string
}

func (t *unixTransport) Call(req []byte) ([]byte, error) {
	conn, err := net.DialTimeout("unix", t.addr, apiTimeout)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	deadline := time.Now().Add(apiTimeout)
	conn.SetDeadline(deadline)

	if _, err := conn.Write(append(req, '\n')); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	var buf bytes.Buffer
	chunk := make([]byte, 4096)
	for {
		n, err := conn.Read(chunk)
		if err != nil {
			return nil, fmt.Errorf("read: %w", err)
		}
		buf.Write(chunk[:n])
		if bytes.Contains(buf.Bytes(), []byte{'\n'}) {
			break
		}
	}

	out := buf.Bytes()
	nl := bytes.IndexByte(out, '\n')
	return out[:nl], nil
}

func (t *unixTransport) Close() error {
	return nil
}
