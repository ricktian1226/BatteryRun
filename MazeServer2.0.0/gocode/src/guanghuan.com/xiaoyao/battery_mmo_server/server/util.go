// util
package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func genID() string {
	u := make([]byte, 16)
	io.ReadFull(rand.Reader, u)
	return hex.EncodeToString(u)
}

func protoErr(state int, i int, buf []byte) error {
	snip := protoSnippet(i, buf)
	err := fmt.Errorf("Parser ERROR, state=%d, i=%d: proto='%s...'",
		state, i, snip)
	return err
}

func protoSnippet(start int, buf []byte) string {
	stop := start + PROTO_SNIPPET_SIZE
	if stop > len(buf) {
		stop = len(buf) - 1
	}
	return fmt.Sprintf("%q", buf[start:stop])
}
