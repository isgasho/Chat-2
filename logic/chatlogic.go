package logic

import (
	"buguang01/chatservice/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var chatgetchan chan string
var chatmodelchan chan *models.ChatChannelModel

//频道的逻辑

func ChatRun(self *models.ChatChannelModel) {
	//处理写入频道和频道信息返回
	next := time.NewTimer(time.Hour * 1)
	for {
		select {
		case <-self.Ctx.Done():
			//正常这里不会进来
			fmt.Println("聊天频道闭关", self.ChatName)
			close(self.SendChan)
			delete(*models.ChatList(), self.ChatName)
			return
		case msg := <-self.SendChan:
			self.UpTime = time.Now()
			next = time.NewTimer(time.Hour * 3)
			if msg != nil {
				//处理逻辑与本频道有关
				if msg.Sendchat == self.ChatName {
					//要写入频道
					var chatmsg = models.ChatModel{
						UID:        models.GetNewUID(),
						CreateTime: time.Now().UTC().Unix(),
						Msg:        msg.Params["Msg"],
					}
					chatmsg.SendMember = msg.SendMember
					if self.ChatList.IsFull() {
						self.ChatList.Poll()
					}
					self.ChatList.Offer(&chatmsg)
				}
				//处理逻辑与本频道无关，只是返回频道信息
				switch self.Chatnum {
				case 0:
					chatbodys, _ := self.ChatList.GetArray(msg.ChatMax)
					msg.Result.Chat1 = chatbodys
					break
				case 1:
					chatbodys, _ := self.ChatList.GetArray(msg.ChatMax)
					msg.Result.Chat2 = chatbodys
					break
				case 2:
					chatbodys, _ := self.ChatList.GetArray(msg.ChatMax)
					msg.Result.Chat3 = chatbodys
					break

				}
			}
			msg.Synwg.Done()
		case <-next.C:
			fmt.Println("聊天频道时间到了闭关", self.ChatName)
			close(self.SendChan)
			delete(*models.ChatList(), self.ChatName)
			return
		}

	}
}

//私人频道
func ChatprivateRun(self *models.ChatChannelModel) {
	//处理写入频道和频道信息返回
	defer close(self.SendChan)
	next := time.NewTimer(time.Hour * 1)
	for {
		select {
		case <-self.Ctx.Done():
			fmt.Println("频道：" + self.ChatName + " 关闭")
			return
		case msg := <-self.SendChan:
			self.UpTime = time.Now()
			next = time.NewTimer(time.Hour * 3)
			if msg != nil {
				fmt.Println("频道：" + self.ChatName + " 处理请求")

				//处理逻辑与本频道有关
				if msg.Sendchat == self.ChatName {
					fmt.Println("频道：" + self.ChatName + " 处理写入请求")

					//要写入频道
					var chatmsg = models.ChatModelToUser{
						ChatModel: models.ChatModel{
							UID:        models.GetNewUID(),
							CreateTime: time.Now().UTC().Unix(),
							Msg:        msg.Params["Msg"],
							SendMember: msg.SendMember,
						},
					}

					tomemid, _ := strconv.Atoi(msg.Params["ToMemberID"])
					chatmsg.ToUser = *GetUserInfo(tomemid)
					if self.ChatList.IsFull() {
						self.ChatList.Poll()
					}
					self.ChatList.Offer(&chatmsg)
				}
				//处理逻辑与本频道无关，只是返回频道信息
				if len(msg.ChatNameList) == 4 && msg.ChatNameList[3] == self.ChatName {
					chatbodys, _ := self.ChatList.GetArray(msg.ChatMax)
					msg.Result.Chat4 = chatbodys
				}
			}
			msg.Synwg.Done()
		case <-next.C:
			fmt.Println("聊天频道时间到了闭关", self.ChatName)
			close(self.SendChan)
			delete(*models.ChatList(), self.ChatName)
			return
		}
	}
}

//聊天的管理协程
func ChatManageRun(ctx context.Context) {
	chatlist := models.ChatList()
	chatgetchan = make(chan string)
	chatmodelchan = make(chan *models.ChatChannelModel)
	fmt.Println("聊天的管理协程")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("聊天的管理协程关闭")
			close(chatgetchan)
			close(chatmodelchan)

			return
		case name := <-chatgetchan:
			if md, ok := (*chatlist)[name]; ok {
				chatmodelchan <- md
			} else {
				//没找到的时候，就新建一个
				md := ChatNewByName(ctx, name)
				chatmodelchan <- md

			}
			break
		default:
			time.Sleep(1)
			break
		}
	}
}

func ChatNewByName(ctx context.Context, name string) *models.ChatChannelModel {
	chatlist := models.ChatList()
	var result = &models.ChatChannelModel{
		ChatList: models.NewMyQueueBySize(100),
		ChatName: name,
		SendChan: make(chan *models.MsgModel, 100),
		UpTime:   time.Now().UTC(),
	}
	if strings.Contains(name, "Member_") {
		result.Chatnum = 3
	} else if strings.Contains(name, "ChatMain") {
		result.Chatnum = 0
	} else if strings.Contains(name, "Chat_") {
		result.Chatnum = 1
	} else if strings.Contains(name, "Custom") {
		result.Chatnum = 2
	}
	result.Ctx, result.Cancel = context.WithCancel(ctx)
	if strings.Contains(name, "Member_") {
		go ChatprivateRun(result)
	} else {
		go ChatRun(result)
	}
	(*chatlist)[name] = result
	fmt.Println("新建了一个聊天频道：" + name)

	return result
}

//获取频道(用户协程上跑的)
func GetChatByChatName(cname string) *models.ChatChannelModel {
	fmt.Println("获取一个聊天频道：" + cname)
	chatlist := models.ChatList()
	var result *models.ChatChannelModel
	if md, ok := (*chatlist)[cname]; ok {
		result = md
	} else {
		chatgetchan <- cname
		result = <-chatmodelchan
	}
	return result
}
