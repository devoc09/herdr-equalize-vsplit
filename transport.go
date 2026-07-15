package main

import "time"

type Transport interface {
	Call(req []byte) ([]byte, error)
	Close() error
}

const apiTimeout = 5 * time.Second
