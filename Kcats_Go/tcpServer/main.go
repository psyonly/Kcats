package main

import (
	"Gin/Webs/tcp.NETv2.0/tools"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// tcp/server/main.go

// TCP server端

// 处理函数
// 建立连接后首先创建该连接对象的管道，并以socket-chan的方式绑定
// 创建一个接收来自用户消息的管道
// 分别执行回传和接收两个协程独立运行
// 轮询检查两个管道缓冲区中是否有消息需要处理
// 缓冲区有消息待发送就将数据从管道中读出到连接中
// 缓冲区有消息收到就打印在服务器中，可以选择不打印减少执行次数
func process(conn net.Conn) {
	// 对应一个连接新建一个协程处理
	var responseBuff = make(chan string, 100)
	chans[conn.RemoteAddr().String()] = &responseBuff
	var usersMSGBuff = make(chan string, 200)
	go response(responseBuff)
	go umsgrecv(conn, usersMSGBuff)
	defer conn.Close() // 关闭连接
	for {
		for len(responseBuff) > 0 {
			str := <-responseBuff   // 屏蔽细节，一个连接的发送器只能发送本缓冲区的消息
			conn.Write([]byte(str)) // 发送数据
		}
		for len(usersMSGBuff) > 0 {
			contentClient := <-usersMSGBuff
			fmt.Println(contentClient)
		}
	}
}

// 回应客户端的输入函数
func response(responseBuff chan string) {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\r') // 读取服务器输入
		inputInfo := strings.Trim(input, "\r\n")
		responseBuff <- inputInfo
		//fmt.Println("LEN ", len(responseBuff))
	}
}

// 接收到用户发来的消息
// 将消息传送给指定接受者的缓冲区
func umsgrecv(conn net.Conn, usersMSGBuff chan string) {
	reader := bufio.NewReader(conn)
	var buf [128]byte
	for {
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		usersMSGBuff <- string(buf[:n])
		msg := tools.ParseMSG(buf[:n])                  // 转换为消息结构体
		usrMap[msg.Sender] = conn.RemoteAddr().String() // 建立用户名-socket对应表用于查询
		fmt.Println(conn.RemoteAddr().String(), msg.Sender)
		if msg.Receiver == "NIL" { // 初始化
			continue
		}
		if msg.Receiver == "SYN" {
			clientList := "NIL" + " " + msg.Sender + "\n"
			clientList += "C"
			for k := range usrMap {
				clientList += " " + k
			}
			clientList += "\r"
			*chans[usrMap[msg.Sender]] <- clientList
			continue
		}
		if v, ok := usrMap[msg.Receiver]; ok { // 检查有没有接收者的chan
			*chans[v] <- string(buf[:n]) // 提取出结构体的接收者，装载消息到指定接受者的chan中缓冲准备分发
		} else { // 若无则返回一个错误模板（详见定义）
			errMSG := tools.ErrMessage
			errMSG.Receiver = msg.Sender
			fmt.Println(string(errMSG.Segment()))
			*chans[usrMap[msg.Sender]] <- string(errMSG.Segment())
		}
	}
}

// ServerInit init the server
func ServerInit() {
	fmt.Println("请输入昵称:")
	userName := ""
	fmt.Scanf("%s\n", &userName)
	if userName != "" {
		nickName = userName
	}
	fmt.Println("是否手动配置监听地址?(Y/n)")
	c := 'n'
	fmt.Scanf("%c\n", &c)
	switch c {
	case 'Y':
		fmt.Scanf("%s\n", &thisSocket)
	default:
		//NULL
	}
}

var thisSocket = "127.0.0.1:9874"
var nickName = "kcats"
var chans = make(map[string]*chan string)
var usrMap = make(map[string]string) // nickName(key) -> socket(value)

func main() {
	ServerInit()
	listen, err := net.Listen("tcp", thisSocket)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() // 建立连接

		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn) // 启动一个goroutine处理连接
	}
}
