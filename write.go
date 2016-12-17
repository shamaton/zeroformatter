package zeroformatter

func (d *serializer) write_s1_i64(value int64, offset uint32) {
	d.create[offset] = byte(value)
}

func (d *serializer) write_s2_i64(value int64, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
}

func (d *serializer) write_s4_i64(value int64, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
	d.create[offset+2] = byte(value >> 16)
	d.create[offset+3] = byte(value >> 24)
}

func (d *serializer) write_s8_i64(value int64, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
	d.create[offset+2] = byte(value >> 16)
	d.create[offset+3] = byte(value >> 24)
	d.create[offset+4] = byte(value >> 32)
	d.create[offset+5] = byte(value >> 40)
	d.create[offset+6] = byte(value >> 48)
	d.create[offset+7] = byte(value >> 56)
}

func (d *serializer) write_s1_u64(value uint64, offset uint32) {
	d.create[offset] = byte(value)
}

func (d *serializer) write_s2_u64(value uint64, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
}

func (d *serializer) write_s4_u64(value uint64, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
	d.create[offset+2] = byte(value >> 16)
	d.create[offset+3] = byte(value >> 24)
}

func (d *serializer) write_s8_u64(value uint64, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
	d.create[offset+2] = byte(value >> 16)
	d.create[offset+3] = byte(value >> 24)
	d.create[offset+4] = byte(value >> 32)
	d.create[offset+5] = byte(value >> 40)
	d.create[offset+6] = byte(value >> 48)
	d.create[offset+7] = byte(value >> 56)
}

func (d *serializer) write_s4_i(value int, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
	d.create[offset+2] = byte(value >> 16)
	d.create[offset+3] = byte(value >> 24)
}

func (d *serializer) write_s4_u32(value uint32, offset uint32) {
	d.create[offset] = byte(value)
	d.create[offset+1] = byte(value >> 8)
	d.create[offset+2] = byte(value >> 16)
	d.create[offset+3] = byte(value >> 24)
}
