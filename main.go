package main

import (
	"buguang01/chatservice/logic"
	"buguang01/chatservice/models"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

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

	// db, _ := connDB()
	// ctx1, cal := context.WithCancel(context.Background())
	// defer cal()
	ctx1, cal1 := context.WithCancel(context.Background())
	ctx2, cal2 := context.WithCancel(context.Background())
	wg.Add(1)
	go logic.ChatManageRun(ctx1, &wg)
	wg.Add(1)
	go SocketListner(ctx2)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
	cal1()
	cal2()
	fmt.Println("等待关闭")
	wg.Wait()
	fmt.Println("关闭完成")

}

var wg sync.WaitGroup

func SocketListner(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("socket close!", err)
		}
	}()
	defer wg.Done()
	flag.Parse()
	tcpid := flag.Arg(0)
	listener, err := net.Listen("tcp", tcpid)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	pool := models.NewPool(flag.Arg(1))
	defer pool.Close()
	go func(ctx context.Context) {
		<-ctx.Done()
		listener.Close()
	}(ctx)
	for {
		conndhandle, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
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
