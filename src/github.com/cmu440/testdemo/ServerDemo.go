package main

import (
	_ "fmt"
	_ "net"
	"fmt"
	"net"
	"strings"
)

var wait = true
func main()  {

	fmt.Println("创建一个服务器.....")

	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil{
		fmt.Println("监听出现错误", err.Error())
		return
	}


	//监听等待连接
	for  wait{
		conn, err := listener.Accept()
		if err != nil{
			fmt.Println("Error  accepting", err.Error())
			return
		}

		go DoServerStuff(conn)

	}
}
//得到一个客户端的连接以后开出处理
func DoServerStuff(conn net.Conn)  {

	for  {
		// 缓存定义的小的话会循环读取
		buf := make([]byte,  1024)
		_,err := conn.Read(buf)
		if err != nil{
			fmt.Println("Error Reading", err.Error())
			return
		}
		receiveStr := string(buf)
		if strings.Compare(receiveStr, "qinjiawei:SH ")==0 {
			fmt.Println("Server ShutDown")
			wait = false
			return
		}

		//收到的任务命令
		fmt.Printf("Receive data: %v \n", receiveStr)
		//判断是不是有效的命令
		if strings.HasPrefix(receiveStr, "get") || strings.HasPrefix(receiveStr, "put") {
			order :=[]byte("ok")
			fmt.Printf("发送.... %v \n", string(order))
			conn.Write(order)
		}else{
			illegalOrder :=[]byte("illegal Order")
			fmt.Printf("发送.... %v \n", string(illegalOrder))
			conn.Write(illegalOrder)
		}

		//conn.Write([]byte("hello, Nice to meet you, my name is SongXingzhu"))
		fmt.Println("回复成功")

	}
}
