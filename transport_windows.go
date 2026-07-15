//go:build windows

package main

import (
	"fmt"
)

func newTransport(socketPath string) Transport {
	return &windowsTransport{pipePath: socketPath}
}

type windowsTransport struct {
	pipePath string
}

func (t *windowsTransport) Call(req []byte) ([]byte, error) {
	return nil, fmt.Errorf("windows named pipe transport not yet implemented; pipe path: %s", t.pipePath)
}

func (t *windowsTransport) Close() error {
	return nil
}
