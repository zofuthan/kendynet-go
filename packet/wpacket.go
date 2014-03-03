package packet

type Wpacket struct{
	writeidx uint32
	buffer *ByteBuffer
}

func NewWpacket(buffer *ByteBuffer)(*Wpacket){
	if buffer == nil {
		return nil
	}
	buffer.PutUint32(0,0)
	return &Wpacket{writeidx:4,buffer:buffer}
}


func (this *Wpacket)Buffer()(*ByteBuffer){
	return this.buffer
}


func (this *Wpacket)PutUint16(value uint16)(error){
	if this.buffer == nil {
		return ErrInvaildData
	}
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutUint16(this.writeidx,value)
	if err != nil{
		return err
	}
	size += 2
	this.writeidx += 2
	this.buffer.PutUint32(0,size)
	return nil
}

func (this *Wpacket)PutUint32(value uint32)(error){
	if this.buffer == nil {
		return ErrInvaildData
	}
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutUint32(this.writeidx,value)
	if err != nil{
		return err
	}
	size += 4
	this.writeidx += 4
	this.buffer.PutUint32(0,size)
	return nil
}

func (this *Wpacket)PutString(value string)(error){
	if this.buffer == nil {
		return ErrInvaildData
	}
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutString(this.writeidx,value)
	if err != nil{
		return err
	}
	size += (4+(uint32)(len(value))+1)
	this.writeidx += (4+(uint32)(len(value))+1)
	this.buffer.PutUint32(0,size)
	return nil
}

func (this *Wpacket)PutBinary(value []byte)(error){
	if this.buffer == nil {
		return ErrInvaildData
	}
	size,err := this.buffer.Uint32(0)
	if err != nil {
		return err
	}
	err = this.buffer.PutBinary(this.writeidx,value)
	if err != nil{
		return err
	}
	size += (4+(uint32)(len(value)))
	this.writeidx += (4+(uint32)(len(value)))
	this.buffer.PutUint32(0,size)
	return nil
}
