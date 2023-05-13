package model

type Event struct {
	// 通用事件参数
	Time     int64  `json:"time"`
	SelfId   int64  `json:"self_id"`
	PostType string `json:"post_type"`

	// 消息事件参数
	MessageType   string `json:"message_type"`
	MessageId     int32  `json:"message_id"`
	UserId        int64  `json:"user_id"`
	RawMessage    string `json:"raw_message"`
	MessageSender Sender `json:"sender"`
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int32  `json:"age"`
}
