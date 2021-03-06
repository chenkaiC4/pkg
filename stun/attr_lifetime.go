package stun

import (
	"github.com/pkg/errors"
)

type Lifetime struct {
	Duration uint32
}

func (x *Lifetime) Pack(message *Message) error {
	v := make([]byte, 4)
	enc.PutUint32(v, x.Duration)
	message.AddAttribute(AttrLifetime, v)
	return nil
}

func (x *Lifetime) Unpack(message *Message, rawAttribute *RawAttribute) error {
	v := rawAttribute.Value

	if len(v) != 4 {
		return errors.Errorf("invalid lifetime length %d != %d (expected)", len(v), 4)
	}

	x.Duration = enc.Uint32(v)

	return nil
}
