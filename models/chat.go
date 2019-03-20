package models

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

//从cache里读出来的信息
type ChatMemberModel struct {
	MemberID int
	Name     string
	ChanID   int
	Info     interface{}
}

//私聊的聊天信息
type ChatModelToUser struct {
	ChatModel
	ToUser ChatMemberModel
}

//要返回的聊天信息
type ChatModel struct {
	SendMember ChatMemberModel
	Msg        string
	CreateTime int64
	UID        int32
}

//返回给用户的频道
type ResultChatModel struct {
	Chat1 *[]interface{}
	Chat2 *[]interface{}
	Chat3 *[]interface{}
	Chat4 *[]interface{}
}

// type ResultChatModel struct {
// 	Chat1 []ChatModel
// 	Chat2 []ChatModel
// 	Chat3 []ChatModel
// 	Chat4 []ChatModel
// }

//切换线程的消息
type MsgModel struct {
	Commd        string            //指令
	Params       map[string]string //参数
	Sendchat     string            //要写入的频道
	SendMember   ChatMemberModel   //要写入的用户
	ChatNameList []string          //要读的频道列表
	ChatMax      int32             //从哪一行开始读
	Result       ResultChatModel
	Synwg        sync.WaitGroup
}

//聊天的频道
type ChatChannelModel struct {
	ChatList *MyQueue
	ChatName string
	SendChan chan *MsgModel //频道收到的消息
	UpTime   time.Time
	Chatnum  int
	Ctx      context.Context
	Cancel   context.CancelFunc
}

var uidmax int32 = 1

//拿每个频道的新的ID
func GetNewUID() int32 {
	result := atomic.AddInt32(&uidmax, 1)
	return result
}
