package main

import (
	"buguang01/chatservice/logic"
	"buguang01/chatservice/models"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/garyburd/redigo/redis"
)

func main() {
	// var jsonBlob = []byte(`[
	//     {"Name": "Platypus", "Order": "Monotremata"},
	//     {"Name": "Quoll",    "Order": "Dasyuromorphia"}
	// ]`)
	// var animals interface{}
	// json.Unmarshal(jsonBlob, &animals)
	// fmt.Println(animals)

	// fmt.Println(time.Now())
	// fmt.Println(time.Now().UTC())
	// fmt.Println(time.Now().UTC().Unix())
	// a := models.MemberModel{}
	// b := &a
	// fmt.Printf("%p\n", &a)
	// fmt.Printf("%p", b)

	//超时没有收到信息可以自我推出
	// next := time.NewTimer(time.Second * 3)
	// for {
	// 	select {
	// 	case <-next.C:
	// 		fmt.Println("聊天频道闭关")
	// 		return
	// 	default:
	// 		fmt.Println(".......")
	// 		time.Sleep(time.Second * 1)
	// 		next = time.NewTimer(time.Second * 3)

	// 		break
	// 	}
	// }
	flag.Parse()
	tcpid := flag.Arg(0)
	listener, err := net.Listen("tcp", tcpid)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	pool := models.NewPool(flag.Arg(1))
	defer pool.Close()
	// db, _ := connDB()
	// ctx1, cal := context.WithCancel(context.Background())
	// defer cal()
	go logic.ChatManageRun(context.Background())
	for {
		conndhandle, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			go logic.UserRun(&conndhandle)
		}
	}

}
func connDB() (c redis.Conn, err error) {
	db, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	return db, err
}
