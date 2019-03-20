package logic

import (
	"buguang01/chatservice/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

func UserRun(conn *net.Conn) {
	connect := *conn

	addr := connect.RemoteAddr().String()

	fmt.Println(addr, "接入服务")
	play := &models.MemberModel{
		ConnHandle: connect,
		Hush:       models.GetNewHash(),
	}
	UserGetChat(play)

}

//收到用户发过来的信息
func UserGetChat(play *models.MemberModel) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	userlist := *models.ClientList()
	defer func() {
		play.ConnHandle.Close()
		fmt.Println("连接断开")

	}()
	buff := make([]byte, 10240)
	status := 0
	fmt.Println("新客户端")
	msgbuff := make([]byte, 0, 10240)
	for {
		n, err := (play.ConnHandle).Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端断开链接，", play.MemberID)
				return
			} else {
				fmt.Println(err)
			}
		}
		msgbuff = append(msgbuff, buff[:n]...)

		if status == 0 { //还没有确认用户

			msg := string(msgbuff)
			n = strings.Index(msg, "\n")
			if n > 0 {
				msg = string(msgbuff[:n-1])
				msgbuff = append(msgbuff[n+1:])
			} else {
				continue
			}
			msglist := strings.Split(msg, " ")
			if len(msglist) < 2 {
				fmt.Println(msg)
				break
			} else if msglist[0] == "Haypi" {
				(play.MemberID), err = strconv.Atoi(msglist[1])
				if err != nil {
					fmt.Println("格式不对")
					break
				}
				play.Send = make(chan *[]byte, 3)
				ctx1, cal := context.WithCancel(context.Background())
				play.Ctx = ctx1
				play.Cancel = cal

				userlist[play.MemberID] = play
				defer func() {
					delete(userlist, play.MemberID)
					play.Cancel()
				}()
				time.Now()
				status = 1
				//告诉客户端hash
				fmt.Println("告诉客户端hash" + strconv.FormatInt(int64(play.Hush), 10))
				play.ConnHandle.Write([]byte("hash " + strconv.FormatInt(int64(play.Hush), 10) + "\r\n"))
				go UserSend(play) //开始可以给这个用户发东西了
			}
		} else if status == 1 { //开始收用户的消息
			msg := string(msgbuff)
			n = strings.Index(msg, "\n")
			if n > 0 {
				msg = string(msgbuff[:n-1])
				msgbuff = append(msgbuff[n+1:])
			} else {
				continue
			}
			msglist := strings.SplitN(msg, " ", 5)
			if msglist[0] == strconv.FormatInt(int64(play.MemberID), 10) && msglist[1] == strconv.FormatInt(int64(play.Hush), 10) {
				fmt.Println("拿到一个用户的消息：" + msg)

				//msglist[2]这后面的是消息指令
				var msg models.MsgModel = models.MsgModel{
					Commd:  msglist[2],
					Result: models.ResultChatModel{},
					Synwg:  sync.WaitGroup{},
				}
				num, _ := strconv.Atoi(msglist[3])
				msg.ChatMax = int32(num)
				json.Unmarshal([]byte(msglist[4]), &msg.Params)
				//拿到用户MC里的信息
				mcuser := GetUserInfo(play.MemberID)
				//拿到要查的频道
				msg.ChatNameList = []string{
					"ChatMain",
					"Chat_" + strconv.FormatInt(int64(mcuser.ChanID), 10), //这里应该读用户身上的chatid
					"Custom" + msg.Params["Custom"],
					"Member_" + strconv.FormatInt(int64(play.MemberID), 10),
				}
				// msg.ChatNameList = strings.Split(msg.Params["ChanList"], ",")
				switch msglist[2] {
				case "sendall": //给所有人
					msg.SendMember = *mcuser
					msg.Sendchat = "ChatMain"
					break
				case "sendto": //给指定人
					msg.SendMember = *mcuser
					msg.Sendchat = "Member_" + msg.Params["ToMemberID"]
					msg.Synwg.Add(1)
					chatmd := GetChatByChatName(msg.Sendchat)
					chatmd.SendChan <- &msg
					break
				case "sendchanid": //给联盟
					msg.SendMember = *mcuser
					msg.Sendchat = "Chat_" + strconv.FormatInt(int64(mcuser.ChanID), 10)
					break
				case "sendcustom": //给指定频道
					msg.SendMember = *mcuser
					msg.Sendchat = "Custom" + msg.Params["Custom"]
					break
				case "get": //拿消息
					msg.Sendchat = ""
					break
				case "quit":
					return
				}
				//上面在对消息填充
				for _, chatname := range msg.ChatNameList {
					msg.Synwg.Add(1)
					chatmd := GetChatByChatName(chatname)
					chatmd.SendChan <- &msg
				}
				fmt.Println("等待消息处理完成。")
				msg.Synwg.Wait()

				//要收的信息都处理完了，下面是写到send协程上去
				result, _ := json.Marshal(&msg.Result)
				play.Send <- &result
			}
		}
	}

	// json.Marshal(&userlist)
}

func UserSend(play *models.MemberModel) {
	msgchan := play.Send
	for {
		select {
		case <-play.Ctx.Done():
			play.ConnHandle.Write([]byte("quit 1"))
			return
		case msg := <-msgchan:
			_, err := play.ConnHandle.Write(*msg)
			if err != nil {
				if err == io.EOF {
					fmt.Println("客户端断开链接，", play.MemberID)
					return
				} else {
					fmt.Println(err)
				}
			}
			break
		}
	}
}

func GetUserInfo(memberid int) *models.ChatMemberModel {
	var mcuser = models.ChatMemberModel{}
	mc := models.GetPool().Get()
	defer mc.Close()
	v, err := redis.Bytes(mc.Do("GET", "Member_"+strconv.FormatInt(int64(memberid), 10)))
	if err != nil {
		panic(err)
	}
	json.Unmarshal(v, &mcuser)
	return &mcuser
}
