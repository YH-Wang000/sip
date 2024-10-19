package enc

import "sip/pkg/sipmsg"

type Encoder interface {
	Encode(data *sipmsg.GenericMessage) ([]byte, error)
}

type EncoderOption func(e Encoder)

func NewEncoder(opts ...EncoderOption) Encoder {
	e := &encoder{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

type encoder struct {
}

func (e *encoder) Encode(data *sipmsg.GenericMessage) ([]byte, error) {
	return nil, nil
}
