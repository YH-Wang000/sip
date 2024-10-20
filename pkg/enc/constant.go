package enc

import "time"

const (
	DefaultMessageLength  = 8192
	InitDecoderBufferSize = 1024

	DefaultReadHeaderTimeout = 20 * time.Second
	DefaultReadBodyTimeout   = 30 * time.Second
)
