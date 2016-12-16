package main


import (
	"net"
	"fmt"
	"bufio"
	"os"
	"strings"
	"io/ioutil"
)

func main() {
	//客户端申请连接
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil{
		fmt.Println("Error dialing", err.Error())
		return
	}

	//开启一个读入流
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("1. what is your name?")
	//读取字段直到换行
	clientname, _ := inputReader.ReadString('\n')
	//过滤掉无用的换行符号
	trimeClientName := strings.Trim(clientname, "\n")


	//发送信息, 直到客户端自己退出
	for   {
		fmt.Println("Send Message To Server, Type Q to Quit")
		input,_ := inputReader.ReadString('\n')
		trimeInput := strings.Trim(input,"\n")

		if(trimeInput == "Q"){
			return
		}
		_,err := conn.Write([] byte(trimeClientName + "says: " + trimeInput ))
		if err != nil{
			return
		}
		fmt.Println("等待接收服务端发送的消息......")

		//接收回复的消息
		handleClient(conn)

		fmt.Println("接收成功")
		//if err != nil {
		//	fmt.Println("接收出错：", err)
		//}
		////var reciveText = string(reciverBuffer[0:len])
		//
		//var reciveText = string(bys)
		//fmt.Println(reciveText)

	}
}
/**
处理接收服务端发送的 消息
 */
func handleClient(conn net.Conn) {

		buf := make([]byte, 1024) // 这个1024可以根据你的消息长度来设置
		n, err := conn.Read(buf) // n为一次Read实际得到的消息长度
		if err != nil{
			fmt.Println("接收的消息出现错误")
		}
		// buf[:n] 就是这次实际读到的消息
		reveiveMsg := string(buf[:n])
		fmt.Println("接收到消息：" +reveiveMsg)

}
/**
处理接收的任务
通过ioutil.ReadAll方法
 */
func handleRecMsgByReadall(conn net.Conn)  {
	bys, err := ioutil.ReadAll(conn) //接收消息。func ReadAll(r io.Reader) ([]byte, error)一定要等到有error或EOF的时候才会返回结果，因此只能等到客户端退出时才会返回结果。
	if err != nil {
		fmt.Println("接收出错：", err)
	}
	//var reciveText = string(reciverBuffer[0:len])

	var reciveText = string(bys)
	fmt.Println(reciveText)
}
