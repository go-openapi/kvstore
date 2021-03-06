package persist

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Value) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Value":
			z.Value, err = dc.ReadBytes(z.Value)
			if err != nil {
				return
			}
		case "Version":
			z.Version, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "LastUpdated":
			z.LastUpdated, err = dc.ReadInt64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Value) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Value"
	err = en.Append(0x83, 0xa5, 0x56, 0x61, 0x6c, 0x75, 0x65)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Value)
	if err != nil {
		return
	}
	// write "Version"
	err = en.Append(0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Version)
	if err != nil {
		return
	}
	// write "LastUpdated"
	err = en.Append(0xab, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.LastUpdated)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Value) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Value"
	o = append(o, 0x83, 0xa5, 0x56, 0x61, 0x6c, 0x75, 0x65)
	o = msgp.AppendBytes(o, z.Value)
	// string "Version"
	o = append(o, 0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendUint64(o, z.Version)
	// string "LastUpdated"
	o = append(o, 0xab, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.LastUpdated)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Value) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Value":
			z.Value, bts, err = msgp.ReadBytesBytes(bts, z.Value)
			if err != nil {
				return
			}
		case "Version":
			z.Version, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "LastUpdated":
			z.LastUpdated, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Value) Msgsize() (s int) {
	s = 1 + 6 + msgp.BytesPrefixSize + len(z.Value) + 8 + msgp.Uint64Size + 12 + msgp.Int64Size
	return
}
