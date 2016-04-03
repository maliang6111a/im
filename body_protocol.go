package main

import (
	"bytes"
	"encoding/binary"
)

type IMMessage struct {
	sender    int64  //8
	receiver  int64  //8
	timestamp int32  //4
	msgid     int32  //4
	content   string //...
}

func (this *IMMessage) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, this.sender)
	binary.Write(buffer, binary.BigEndian, this.receiver)
	binary.Write(buffer, binary.BigEndian, this.timestamp)
	binary.Write(buffer, binary.BigEndian, this.msgid)
	buffer.Write([]byte(this.content))
	return buffer.Bytes()
}

func (this *IMMessage) FromData(buff []byte) bool {
	if len(buff) < 24 {
		return false
	}
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &this.sender)
	binary.Read(buffer, binary.BigEndian, &this.receiver)
	binary.Read(buffer, binary.BigEndian, &this.timestamp)
	binary.Read(buffer, binary.BigEndian, &this.msgid)
	this.content = string(buff[24:])
	return true
}
