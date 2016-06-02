// const
package server

import (
	"time"
)

const (
	// VERSION is the current version for the server.
	VERSION = "1.0.0"

	// ACCEPT_MIN_SLEEP is the minimum acceptable sleep times on temporary errors.
	ACCEPT_MIN_SLEEP = 10 * time.Millisecond
	// ACCEPT_MAX_SLEEP is the maximum acceptable sleep times on temporary errors
	ACCEPT_MAX_SLEEP = 1 * time.Second

	// PROTO_SNIPPET_SIZE is the default size of proto to print on parse errors.
	PROTO_SNIPPET_SIZE = 32

	//报文正文内容，最大设置为10k
	PROTO_MAX_CONTENT_SIZE = 10240
)
