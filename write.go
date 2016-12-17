package zeroformatter

func (s *serializer) write_s1_i64(value int64, offset uint32) {
	s.create[offset] = byte(value)
}

func (s *serializer) write_s2_i64(value int64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
}

func (s *serializer) write_s4_i64(value int64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}

func (s *serializer) write_s8_i64(value int64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
	s.create[offset+4] = byte(value >> 32)
	s.create[offset+5] = byte(value >> 40)
	s.create[offset+6] = byte(value >> 48)
	s.create[offset+7] = byte(value >> 56)
}

func (s *serializer) write_s1_u64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
}

func (s *serializer) write_s2_u64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
}

func (s *serializer) write_s4_u64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}

func (s *serializer) write_s8_u64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
	s.create[offset+4] = byte(value >> 32)
	s.create[offset+5] = byte(value >> 40)
	s.create[offset+6] = byte(value >> 48)
	s.create[offset+7] = byte(value >> 56)
}

func (s *serializer) write_s4_i(value int, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}

func (s *serializer) write_s4_u32(value uint32, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}
