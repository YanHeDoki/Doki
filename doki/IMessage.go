package doki

type IMessage interface {
	//获取消息的ID
	GetMsgId() uint32
	//获取消息的长度
	GetDataLen() uint32
	//获取消息的内容
	GetData() []byte

	//设置消息的id
	SetMsgId(uint32)
	//设置消息的长度
	SetDataLen(uint32)
	//设置消息的内容
	SetData([]byte)
}
