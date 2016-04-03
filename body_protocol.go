package main

import (
	"bytes"
	"encoding/binary"
	//"fmt"
)

type IMMessage struct {
	sender    int64  //8
	receiver  int64  //8
	timestamp int32  //4
	msgid     int32  //4
	content   string //...
}

//认证信息
type AuthMessage struct {
	//authIdLen  int32 //4
	//authPwdLen int32 //4
	authId  string
	authPwd string
}

func (this *AuthMessage) ToData() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, int32(len([]byte(this.authId))))
	binary.Write(buffer, binary.BigEndian, int32(len([]byte(this.authPwd))))
	buffer.Write([]byte(this.authId))
	buffer.Write([]byte(this.authPwd))
	return buffer.Bytes()
}

func (this *AuthMessage) FromData(buff []byte) bool {
	if len(buff) < 8 {
		return false
	}
	var authIdLen, authPwdLen int32
	buffer := bytes.NewBuffer(buff)
	binary.Read(buffer, binary.BigEndian, &authIdLen)
	binary.Read(buffer, binary.BigEndian, &authPwdLen)
	//fmt.Println(this.authIdLen, this.authPwdLen)
	this.authId = string(buff[8 : authIdLen+8])
	this.authPwd = string(buff[authIdLen+8 : authIdLen+authPwdLen+8])
	return true
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
