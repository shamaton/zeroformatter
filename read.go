package zeroformatter

func (d *deserializer) read_s1(index uint32) (byte, uint32) {
	rb := byte1
	return d.data[index], index + rb
}

func (d *deserializer) read_s2(index uint32) ([]byte, uint32) {
	rb := byte2
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) read_s4(index uint32) ([]byte, uint32) {
	rb := byte4
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) read_s8(index uint32) ([]byte, uint32) {
	rb := byte8
	return d.data[index : index+rb], index + rb
}
