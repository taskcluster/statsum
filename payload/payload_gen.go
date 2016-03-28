package payload

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Counter) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "k":
			z.Key, err = dc.ReadString()
			if err != nil {
				return
			}
		case "v":
			z.Value, err = dc.ReadFloat64()
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
func (z Counter) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "k"
	err = en.Append(0x82, 0xa1, 0x6b)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Key)
	if err != nil {
		return
	}
	// write "v"
	err = en.Append(0xa1, 0x76)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Value)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Counter) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "k"
	o = append(o, 0x82, 0xa1, 0x6b)
	o = msgp.AppendString(o, z.Key)
	// string "v"
	o = append(o, 0xa1, 0x76)
	o = msgp.AppendFloat64(o, z.Value)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Counter) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "k":
			z.Key, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "v":
			z.Value, bts, err = msgp.ReadFloat64Bytes(bts)
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

func (z Counter) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.Key) + 2 + msgp.Float64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Measure) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "k":
			z.Key, err = dc.ReadString()
			if err != nil {
				return
			}
		case "v":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Value) >= int(xsz) {
				z.Value = z.Value[:xsz]
			} else {
				z.Value = make([]float64, xsz)
			}
			for xvk := range z.Value {
				z.Value[xvk], err = dc.ReadFloat64()
				if err != nil {
					return
				}
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
func (z *Measure) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "k"
	err = en.Append(0x82, 0xa1, 0x6b)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Key)
	if err != nil {
		return
	}
	// write "v"
	err = en.Append(0xa1, 0x76)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Value)))
	if err != nil {
		return
	}
	for xvk := range z.Value {
		err = en.WriteFloat64(z.Value[xvk])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Measure) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "k"
	o = append(o, 0x82, 0xa1, 0x6b)
	o = msgp.AppendString(o, z.Key)
	// string "v"
	o = append(o, 0xa1, 0x76)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Value)))
	for xvk := range z.Value {
		o = msgp.AppendFloat64(o, z.Value[xvk])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Measure) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "k":
			z.Key, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "v":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Value) >= int(xsz) {
				z.Value = z.Value[:xsz]
			} else {
				z.Value = make([]float64, xsz)
			}
			for xvk := range z.Value {
				z.Value[xvk], bts, err = msgp.ReadFloat64Bytes(bts)
				if err != nil {
					return
				}
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

func (z *Measure) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.Key) + 2 + msgp.ArrayHeaderSize + (len(z.Value) * (msgp.Float64Size))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Payload) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "counters":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Counters) >= int(xsz) {
				z.Counters = z.Counters[:xsz]
			} else {
				z.Counters = make([]Counter, xsz)
			}
			for bzg := range z.Counters {
				var isz uint32
				isz, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for isz > 0 {
					isz--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "k":
						z.Counters[bzg].Key, err = dc.ReadString()
						if err != nil {
							return
						}
					case "v":
						z.Counters[bzg].Value, err = dc.ReadFloat64()
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
			}
		case "measures":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Measures) >= int(xsz) {
				z.Measures = z.Measures[:xsz]
			} else {
				z.Measures = make([]Measure, xsz)
			}
			for bai := range z.Measures {
				var isz uint32
				isz, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for isz > 0 {
					isz--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "k":
						z.Measures[bai].Key, err = dc.ReadString()
						if err != nil {
							return
						}
					case "v":
						var xsz uint32
						xsz, err = dc.ReadArrayHeader()
						if err != nil {
							return
						}
						if cap(z.Measures[bai].Value) >= int(xsz) {
							z.Measures[bai].Value = z.Measures[bai].Value[:xsz]
						} else {
							z.Measures[bai].Value = make([]float64, xsz)
						}
						for cmr := range z.Measures[bai].Value {
							z.Measures[bai].Value[cmr], err = dc.ReadFloat64()
							if err != nil {
								return
							}
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
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
func (z *Payload) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "counters"
	err = en.Append(0x82, 0xa8, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Counters)))
	if err != nil {
		return
	}
	for bzg := range z.Counters {
		// map header, size 2
		// write "k"
		err = en.Append(0x82, 0xa1, 0x6b)
		if err != nil {
			return err
		}
		err = en.WriteString(z.Counters[bzg].Key)
		if err != nil {
			return
		}
		// write "v"
		err = en.Append(0xa1, 0x76)
		if err != nil {
			return err
		}
		err = en.WriteFloat64(z.Counters[bzg].Value)
		if err != nil {
			return
		}
	}
	// write "measures"
	err = en.Append(0xa8, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Measures)))
	if err != nil {
		return
	}
	for bai := range z.Measures {
		// map header, size 2
		// write "k"
		err = en.Append(0x82, 0xa1, 0x6b)
		if err != nil {
			return err
		}
		err = en.WriteString(z.Measures[bai].Key)
		if err != nil {
			return
		}
		// write "v"
		err = en.Append(0xa1, 0x76)
		if err != nil {
			return err
		}
		err = en.WriteArrayHeader(uint32(len(z.Measures[bai].Value)))
		if err != nil {
			return
		}
		for cmr := range z.Measures[bai].Value {
			err = en.WriteFloat64(z.Measures[bai].Value[cmr])
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Payload) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "counters"
	o = append(o, 0x82, 0xa8, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Counters)))
	for bzg := range z.Counters {
		// map header, size 2
		// string "k"
		o = append(o, 0x82, 0xa1, 0x6b)
		o = msgp.AppendString(o, z.Counters[bzg].Key)
		// string "v"
		o = append(o, 0xa1, 0x76)
		o = msgp.AppendFloat64(o, z.Counters[bzg].Value)
	}
	// string "measures"
	o = append(o, 0xa8, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Measures)))
	for bai := range z.Measures {
		// map header, size 2
		// string "k"
		o = append(o, 0x82, 0xa1, 0x6b)
		o = msgp.AppendString(o, z.Measures[bai].Key)
		// string "v"
		o = append(o, 0xa1, 0x76)
		o = msgp.AppendArrayHeader(o, uint32(len(z.Measures[bai].Value)))
		for cmr := range z.Measures[bai].Value {
			o = msgp.AppendFloat64(o, z.Measures[bai].Value[cmr])
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Payload) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "counters":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Counters) >= int(xsz) {
				z.Counters = z.Counters[:xsz]
			} else {
				z.Counters = make([]Counter, xsz)
			}
			for bzg := range z.Counters {
				var isz uint32
				isz, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for isz > 0 {
					isz--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "k":
						z.Counters[bzg].Key, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "v":
						z.Counters[bzg].Value, bts, err = msgp.ReadFloat64Bytes(bts)
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
			}
		case "measures":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Measures) >= int(xsz) {
				z.Measures = z.Measures[:xsz]
			} else {
				z.Measures = make([]Measure, xsz)
			}
			for bai := range z.Measures {
				var isz uint32
				isz, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for isz > 0 {
					isz--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "k":
						z.Measures[bai].Key, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "v":
						var xsz uint32
						xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
						if err != nil {
							return
						}
						if cap(z.Measures[bai].Value) >= int(xsz) {
							z.Measures[bai].Value = z.Measures[bai].Value[:xsz]
						} else {
							z.Measures[bai].Value = make([]float64, xsz)
						}
						for cmr := range z.Measures[bai].Value {
							z.Measures[bai].Value[cmr], bts, err = msgp.ReadFloat64Bytes(bts)
							if err != nil {
								return
							}
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
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

func (z *Payload) Msgsize() (s int) {
	s = 1 + 9 + msgp.ArrayHeaderSize
	for bzg := range z.Counters {
		s += 1 + 2 + msgp.StringPrefixSize + len(z.Counters[bzg].Key) + 2 + msgp.Float64Size
	}
	s += 9 + msgp.ArrayHeaderSize
	for bai := range z.Measures {
		s += 1 + 2 + msgp.StringPrefixSize + len(z.Measures[bai].Key) + 2 + msgp.ArrayHeaderSize + (len(z.Measures[bai].Value) * (msgp.Float64Size))
	}
	return
}
