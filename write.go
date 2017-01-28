package zeroformatter

func (s *serializer) writeSize1Int64(value int64, offset uint32) {
	s.create[offset] = byte(value)
}

func (s *serializer) writeSize2Int64(value int64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
}

func (s *serializer) writeSize4Int64(value int64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}

func (s *serializer) writeSize8Int64(value int64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
	s.create[offset+4] = byte(value >> 32)
	s.create[offset+5] = byte(value >> 40)
	s.create[offset+6] = byte(value >> 48)
	s.create[offset+7] = byte(value >> 56)
}

func (s *serializer) writeSize1Uint64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
}

func (s *serializer) writeSize2Uint64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
}

func (s *serializer) writeSize4Uint64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}

func (s *serializer) writeSize8Uint64(value uint64, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
	s.create[offset+4] = byte(value >> 32)
	s.create[offset+5] = byte(value >> 40)
	s.create[offset+6] = byte(value >> 48)
	s.create[offset+7] = byte(value >> 56)
}

func (s *serializer) writeSize4Int(value int, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}

func (s *serializer) writeSize4Uint32(value uint32, offset uint32) {
	s.create[offset] = byte(value)
	s.create[offset+1] = byte(value >> 8)
	s.create[offset+2] = byte(value >> 16)
	s.create[offset+3] = byte(value >> 24)
}
