package main

import (
	"Gin/Webs/tcp.NETv2.0/tools"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// tcp/client/main.go

//
const (
	RUNNING = true
	SHUTDWN = false
)

var inputReader = bufio.NewReader(os.Stdin)
var nickName = ""
var thisSocket = ""
var targetIP = ""

var (
	// STATUS 描述客户端运行状态
	STATUS = RUNNING
)

// 客户端
func main() {
	initClient()
	conn, err := net.Dial("tcp", targetIP)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	thisSocket = tools.GetIP()
	inputBuff := make(chan []byte, 100)
	recevBuff := make(chan []byte, 200)
	go input(inputBuff)
	go receiver(recevBuff, conn)
	time.Sleep(time.Second * 5)
	defer conn.Close() // 关闭连接
	for {
		if STATUS == SHUTDWN {
			break
		}
		for len(inputBuff) > 0 {
			inputStr := <-inputBuff
			_, err := conn.Write([]byte(inputStr)) // 发送数据
			if err != nil {
				return
			}
		}
		for len(recevBuff) > 0 {
			// 读数据
			receiveINFO := <-recevBuff
			fmt.Println(tools.DeCodeMSG(receiveINFO))
		}
	}
}

func initClient() {
	fmt.Println("请输入昵称:")
	fmt.Scanf("%s\n", &nickName)
	fmt.Println("输入命令'/help 查看全部命令'")
	fmt.Println("Go Client>>> 请输入所要连接的Socket")
	fmt.Scanf("%s\n", &targetIP)
	fmt.Println("连接目标IP为：", targetIP, "正在建立链接中...")
	fmt.Println(targetIP + "连接成功")
}

func input(ch chan []byte) {
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		if input[0] == '/' {
			cmd := strings.ToUpper(input)
			switch cmd {
			case "QUIT":
				STATUS = SHUTDWN
			case "HELP":
				fmt.Println("发送消息的语法：xxx>>tttt，其中xxx是你想要发给的人的昵称，tttt是你的消息内容。")
			}
			continue
		}
		msgStr := strings.Split(input, ">>")
		msg := tools.NewMSG(nickName, msgStr[0], msgStr[1])
		inputInfo := msg.Segment()
		ch <- inputInfo
	}
}

func receiver(ch chan []byte, conn net.Conn) {
	buf := [512]byte{}
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("接收不到来自服务端的消息，原因可能是:", err)
			return
		}
		ch <- buf[:n]
	}
}
