// Implementation of a KeyValueServer. Students should write their code in this file.

package p0

import (
	"net"
	"strconv"

	"bufio"
	"bytes"
	"fmt"
	"io"
)

//定义一下 消息的通道的最大长度
const Max_MESSAGE = 500

type keyValueServer struct {
	// TODO: implement this!

	listener          net.Listener  // 监听,只是开一个的（因为多个客户端都是链接再这一个上）
	currentClients    []*client     //当前连接客户端
	newMessage        chan []byte   //新的消息
	newConnection     chan net.Conn //新的连接
	deadclient        chan *client  //过期的客户端
	quitSignal_main   chan int      //离开信号
	dbQuery           chan *db      //存储数据库命令的管道
	countClientnum    chan int
	clientCount       chan int
	quitSignal_Accept chan int
	debug             bool
}

//客户端的链接
type client struct {
	connetion        net.Conn
	messageQueue     chan []byte
	quitSignal_Read  chan int
	quitSignal_Write chan int
}

//存储一条命令的结构体
type db struct {
	isGet bool
	key   string
	value []byte
}

// New creates and returns (but does not start) a new KeyValueServer.
func New() KeyValueServer {
	// TODO: implement this!

	server := &keyValueServer{
		nil,
		nil,
		make(chan []byte),
		make(chan net.Conn),
		make(chan *client),
		make(chan int),
		make(chan *db),
		make(chan int),
		make(chan int),
		make(chan int),
		true,
	}
	//fmt.Println("starting the server........")
	//go runServer(server)
	return server
}

func (kvs *keyValueServer) Start(port int) error {
	// TODO: implement this!

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port)) //返回数字 i 所表示的字符串类型的十进制数。
	if err != nil {
		return err
	}
	kvs.listener = ln
	//初始化存储
	init_db()

	go runServer(kvs)
	go acceptRoutine(kvs)

	return nil
}

func (kvs *keyValueServer) Close() {
	// TODO: implement this!
	kvs.listener.Close()
	kvs.quitSignal_main <- 0
	kvs.quitSignal_Accept <- 0

}

func (kvs *keyValueServer) Count() int {
	// TODO: implement this!
	kvs.countClientnum <- 0

	return <-kvs.clientCount
}

// TODO: add additional methods/functions below!
//监听端口开启以后, 等待连接
func acceptRoutine(kvs *keyValueServer) {

	for {
		select {
		case <-kvs.quitSignal_Accept:
			return
		default:
			//接受新的连接
			conn, err := kvs.listener.Accept()
			if err == nil {
				//新的连接加入到管道中，这样避免锁
				kvs.newConnection <- conn
			}

		}

	}
}

//开始运行服务
func runServer(kvs *keyValueServer) {

	for {
		select {
		//把消息分发个每一个客户端
		case newmessage := <-kvs.newMessage:
			//fmt.Println("get the newmessage from chan: "+string(newmessage))
			// 遍历客户端， 舍弃数组的下标
			for _, c := range kvs.currentClients {

				//当通道中的内容达到最大的时候，
				if len(c.messageQueue) == Max_MESSAGE {

				} else {
					//新的消息写入
					c.messageQueue <- newmessage
				}

			}

			//新的客户端链接进来
		case newConnection := <-kvs.newConnection:
			c := &client{
				newConnection,
				make(chan []byte, Max_MESSAGE),
				make(chan int),
				make(chan int), //由于换行了就需要加上 逗号
			}
			//新的连接加到数组中去
			kvs.currentClients = append(kvs.currentClients, c)

			// 读取服务器上的消息
			go ReadAndSendToServer(kvs, c)

			//向服务器传递数据
			go WriteMessage(c)

			//读出一条命令
		case dbQuery := <-kvs.dbQuery:

			//  Get 命令
			if dbQuery.isGet {
				value := get(dbQuery.key)
				if kvs.debug {
					//fmt.Printf("get %s value %s  \n",dbQuery.key, string(value))
				}

				//fmt.Println("put message byte in the newMessage")

				newValue := append(append([]byte(dbQuery.key), ","...), value...)
				//fmt.Println("get message :" +string(newValue))
				//自己给自己发送任务就会出现死锁的现象，解决办法是可以把channeal的空间设置大一点，相当于有个缓存的俄存在。
				kvs.newMessage <- newValue
				//上述的办法还有一个就是直接把结果发给client端，不去存在自己的当中
				//死锁的产生是因为,kvs.newMessage必须被读出来才会继续运行下去，但是如果想被运行出来,程序还是必须回到上边,但是程序已经回不去了.
				//fmt.Println("put......")

			} else { //put命令
				if kvs.debug {
					//fmt.Println("put value in .....")
				}
				put(dbQuery.key, dbQuery.value)

				if kvs.debug {
					//fmt.Println("put value completely .....")
				}
			}

			//服务器退出
		case <-kvs.quitSignal_main:
			for _, c := range kvs.currentClients {
				if kvs.debug {
					//fmt.Println(i)
				}
				c.connetion.Close()
				c.quitSignal_Write <- 0
				c.quitSignal_Read <- 0
			}
			return

			//计算当前的客户端数量
		case <-kvs.countClientnum:
			kvs.clientCount <- len(kvs.currentClients)

		case deadclient := <-kvs.deadclient:

			for i, c := range kvs.currentClients {

				if c == deadclient {

					c.quitSignal_Write <- 0
					c.quitSignal_Read <- 0
					kvs.currentClients = append(kvs.currentClients[:i], kvs.currentClients[i+1:]...)
					break
				}
			}

		}

	}

}

//服务器读取客户端的消息
func ReadAndSendToServer(kvs *keyValueServer, c *client) {

	//读入缓存中
	clientReader := bufio.NewReader(c.connetion)

	for {
		select {
		//相当于读取离开信号
		case <-c.quitSignal_Read:
			fmt.Println("client exit read")
			return

		default:
			//message 中读取的是byte类型的数据
			message, err := clientReader.ReadBytes('\n')
			// 终端了，就是结束了
			if err == io.EOF {
				//退出的客户端
				kvs.deadclient <- c

			} else if err != nil {
				return
			} else {
				//有两种命令  put,key,value 或 get,key
				tokens := bytes.Split(message, []byte(","))

				// 强制的类型转换  byte转换成 string
				if string(tokens[0]) == "put" {

					key := string(tokens[1][:])
					value := tokens[2][:] //去掉换行
					kvs.dbQuery <- &db{
						isGet: false,
						key:   key,
						value: value, //byte类型的数据
					}
					if kvs.debug {
						fmt.Printf("put key %s with value %s \n", key, string(value))
					}

				} else if string(tokens[0]) == "get" {

					//tokens[1] 就是key但是注意它有“\n”字符
					key := string(tokens[1][:len(tokens[1])-1])

					kvs.dbQuery <- &db{
						isGet: true,
						key:   key,
					}

				}

			}
		}
	}

}

//服务器向客户端返回的数据
func WriteMessage(c *client) {

	for {
		select {
		case <-c.quitSignal_Write:
			fmt.Println("cleint exit write")
			return
		case message := <-c.messageQueue:
			//fmt.Println("send message from server to client")
			c.connetion.Write(message)
			//fmt.Println("message send ok")
		}

	}
}
