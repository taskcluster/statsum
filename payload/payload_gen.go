package payload

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *CountMetric) DecodeMsg(dc *msgp.Reader) (err error) {
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
func (z CountMetric) EncodeMsg(en *msgp.Writer) (err error) {
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
func (z CountMetric) MarshalMsg(b []byte) (o []byte, err error) {
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
func (z *CountMetric) UnmarshalMsg(bts []byte) (o []byte, err error) {
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

func (z CountMetric) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.Key) + 2 + msgp.Float64Size
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
		case "countMetrics":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.CountMetrics) >= int(xsz) {
				z.CountMetrics = z.CountMetrics[:xsz]
			} else {
				z.CountMetrics = make([]CountMetric, xsz)
			}
			for xvk := range z.CountMetrics {
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
						z.CountMetrics[xvk].Key, err = dc.ReadString()
						if err != nil {
							return
						}
					case "v":
						z.CountMetrics[xvk].Value, err = dc.ReadFloat64()
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
		case "valueMetrics":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.ValueMetrics) >= int(xsz) {
				z.ValueMetrics = z.ValueMetrics[:xsz]
			} else {
				z.ValueMetrics = make([]ValueMetric, xsz)
			}
			for bzg := range z.ValueMetrics {
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
						z.ValueMetrics[bzg].Key, err = dc.ReadString()
						if err != nil {
							return
						}
					case "v":
						var xsz uint32
						xsz, err = dc.ReadArrayHeader()
						if err != nil {
							return
						}
						if cap(z.ValueMetrics[bzg].Value) >= int(xsz) {
							z.ValueMetrics[bzg].Value = z.ValueMetrics[bzg].Value[:xsz]
						} else {
							z.ValueMetrics[bzg].Value = make([]float64, xsz)
						}
						for bai := range z.ValueMetrics[bzg].Value {
							z.ValueMetrics[bzg].Value[bai], err = dc.ReadFloat64()
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
	// write "countMetrics"
	err = en.Append(0x82, 0xac, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.CountMetrics)))
	if err != nil {
		return
	}
	for xvk := range z.CountMetrics {
		// map header, size 2
		// write "k"
		err = en.Append(0x82, 0xa1, 0x6b)
		if err != nil {
			return err
		}
		err = en.WriteString(z.CountMetrics[xvk].Key)
		if err != nil {
			return
		}
		// write "v"
		err = en.Append(0xa1, 0x76)
		if err != nil {
			return err
		}
		err = en.WriteFloat64(z.CountMetrics[xvk].Value)
		if err != nil {
			return
		}
	}
	// write "valueMetrics"
	err = en.Append(0xac, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.ValueMetrics)))
	if err != nil {
		return
	}
	for bzg := range z.ValueMetrics {
		// map header, size 2
		// write "k"
		err = en.Append(0x82, 0xa1, 0x6b)
		if err != nil {
			return err
		}
		err = en.WriteString(z.ValueMetrics[bzg].Key)
		if err != nil {
			return
		}
		// write "v"
		err = en.Append(0xa1, 0x76)
		if err != nil {
			return err
		}
		err = en.WriteArrayHeader(uint32(len(z.ValueMetrics[bzg].Value)))
		if err != nil {
			return
		}
		for bai := range z.ValueMetrics[bzg].Value {
			err = en.WriteFloat64(z.ValueMetrics[bzg].Value[bai])
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
	// string "countMetrics"
	o = append(o, 0x82, 0xac, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.CountMetrics)))
	for xvk := range z.CountMetrics {
		// map header, size 2
		// string "k"
		o = append(o, 0x82, 0xa1, 0x6b)
		o = msgp.AppendString(o, z.CountMetrics[xvk].Key)
		// string "v"
		o = append(o, 0xa1, 0x76)
		o = msgp.AppendFloat64(o, z.CountMetrics[xvk].Value)
	}
	// string "valueMetrics"
	o = append(o, 0xac, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.ValueMetrics)))
	for bzg := range z.ValueMetrics {
		// map header, size 2
		// string "k"
		o = append(o, 0x82, 0xa1, 0x6b)
		o = msgp.AppendString(o, z.ValueMetrics[bzg].Key)
		// string "v"
		o = append(o, 0xa1, 0x76)
		o = msgp.AppendArrayHeader(o, uint32(len(z.ValueMetrics[bzg].Value)))
		for bai := range z.ValueMetrics[bzg].Value {
			o = msgp.AppendFloat64(o, z.ValueMetrics[bzg].Value[bai])
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
		case "countMetrics":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.CountMetrics) >= int(xsz) {
				z.CountMetrics = z.CountMetrics[:xsz]
			} else {
				z.CountMetrics = make([]CountMetric, xsz)
			}
			for xvk := range z.CountMetrics {
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
						z.CountMetrics[xvk].Key, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "v":
						z.CountMetrics[xvk].Value, bts, err = msgp.ReadFloat64Bytes(bts)
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
		case "valueMetrics":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.ValueMetrics) >= int(xsz) {
				z.ValueMetrics = z.ValueMetrics[:xsz]
			} else {
				z.ValueMetrics = make([]ValueMetric, xsz)
			}
			for bzg := range z.ValueMetrics {
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
						z.ValueMetrics[bzg].Key, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "v":
						var xsz uint32
						xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
						if err != nil {
							return
						}
						if cap(z.ValueMetrics[bzg].Value) >= int(xsz) {
							z.ValueMetrics[bzg].Value = z.ValueMetrics[bzg].Value[:xsz]
						} else {
							z.ValueMetrics[bzg].Value = make([]float64, xsz)
						}
						for bai := range z.ValueMetrics[bzg].Value {
							z.ValueMetrics[bzg].Value[bai], bts, err = msgp.ReadFloat64Bytes(bts)
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
	s = 1 + 13 + msgp.ArrayHeaderSize
	for xvk := range z.CountMetrics {
		s += 1 + 2 + msgp.StringPrefixSize + len(z.CountMetrics[xvk].Key) + 2 + msgp.Float64Size
	}
	s += 13 + msgp.ArrayHeaderSize
	for bzg := range z.ValueMetrics {
		s += 1 + 2 + msgp.StringPrefixSize + len(z.ValueMetrics[bzg].Key) + 2 + msgp.ArrayHeaderSize + (len(z.ValueMetrics[bzg].Value) * (msgp.Float64Size))
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *ValueMetric) DecodeMsg(dc *msgp.Reader) (err error) {
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
			for cmr := range z.Value {
				z.Value[cmr], err = dc.ReadFloat64()
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
func (z *ValueMetric) EncodeMsg(en *msgp.Writer) (err error) {
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
	for cmr := range z.Value {
		err = en.WriteFloat64(z.Value[cmr])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ValueMetric) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "k"
	o = append(o, 0x82, 0xa1, 0x6b)
	o = msgp.AppendString(o, z.Key)
	// string "v"
	o = append(o, 0xa1, 0x76)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Value)))
	for cmr := range z.Value {
		o = msgp.AppendFloat64(o, z.Value[cmr])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ValueMetric) UnmarshalMsg(bts []byte) (o []byte, err error) {
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
			for cmr := range z.Value {
				z.Value[cmr], bts, err = msgp.ReadFloat64Bytes(bts)
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

func (z *ValueMetric) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.Key) + 2 + msgp.ArrayHeaderSize + (len(z.Value) * (msgp.Float64Size))
	return
}
