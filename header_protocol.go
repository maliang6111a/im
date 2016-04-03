package main

//消息
type Message struct {
	version  int         //1,byte,消息协议版本
	msg_type int         //1,byte,消息类型
	body     interface{} //len
}

func (this *Message) ToData() []byte {
	if this.body != nil {
		if m, ok := this.body.(IMessage); ok {
			return m.ToData()
		}
	}
	return make([]byte, 0)
}

func (this *Message) FromData(buff []byte) bool {
	proVersion := this.version
	if creator, ok := message_creators[proVersion]; ok {
		c := creator()
		r := c.FromData(buff)
		this.body = c
		return r
	}
	return len(buff) == 0
}
