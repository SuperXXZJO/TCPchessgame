package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	ChessBoard()
	Start()
}

// 启动服务器
func Start() {
	listener, err := net.Listen("tcp", "59.110.23.117:9090")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 建立连接池，用于广播消息
	conns := make(map[string]net.Conn)

	// 消息通道
	messageChan := make(chan string, 10)

	// 广播消息
	go BroadMessages(&conns, messageChan)

	// 启动
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept failed:%v\n", err)
			continue
		}

		// 把每个客户端连接扔进连接池
		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)

		// 处理消息
		go Handler(conn, &conns, messageChan)
	}
}

// 向所有连接上的乡亲们发广播
func BroadMessages(conns *map[string]net.Conn, messages chan string) {
	for {

		// 不断从通道里读取消息
		msg := <-messages

		lens := strings.Count(msg, "")
		if lens == 4 {
			msg = Board[0][0] + Board[0][1] + Board[0][2] + Board[0][3] + Board[0][4] + "\n" + Board[1][0] + Board[1][1] + Board[1][2] + Board[1][3] + Board[1][4] + "\n" + Board[2][0] + Board[2][1] + Board[2][2] + Board[2][3] + Board[2][4] + "\n" + Board[3][0] + Board[3][1] + Board[3][2] + Board[3][3] + Board[3][4] + "\n" + Board[4][0] + Board[4][1] + Board[4][2] + Board[4][3] + Board[4][4] + "\n"
		} else {
			msg = "winner" + msg
		}

		// 向所有的乡亲们发消息
		for key, conn := range *conns {
			fmt.Println("connection is connected from ", key)
			fmt.Println(msg)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("broad message to %s failed: %v\n", key, err)
				delete(*conns, key)
			}
		}
	}
}

// 处理客户端发到服务端的消息，将其扔到通道中
func Handler(conn net.Conn, conns *map[string]net.Conn, messages chan string) {
	fmt.Println("connect from client ", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("read client message failed:%v\n", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}

		// 把收到的消息写到通道中
		recvStr := string(buf[0:length])
		res := strings.Split(recvStr, ",")
		zb, winner := chess(res)
		if winner != " " && winner != "0" {
			msg := winner
			messages <- msg
		} else {
			msg := zb
			messages <- msg
		}
	}
}

const MaxRow = 5  //行
const MaxLine = 5 //列
var Board [MaxRow][MaxLine]string

//打印棋盘
func ChessBoard() {
	//初始化一个棋盘
	for row := 0; row < 5; row++ {
		for line := 0; line < 5; line++ {
			Board[row][line] = "0"
		}
	}
}

func chess(mod []string) (string, string) {
	x, _ := strconv.Atoi(mod[1])
	y, _ := strconv.Atoi(mod[2])
	//判断是否位置已经落子
	zb := mod[1] + "," + mod[2]
	Board[x][y] = mod[0]
	winner := CheckWinner()
	//if winner != " "  {
	//	return zb, winner
	//}
	return zb, winner

}

//判断胜负
func CheckWinner() string {
	//判断是否连成一行
	for row := 0; row < 5; row++ {
		if Board[row][0] == Board[row][1] && Board[row][0] == Board[row][2] && Board[row][0] == Board[row][3] && Board[row][0] == Board[row][4] {
			return Board[row][0]

		}
	}
	//判断是否连成一列
	for line := 0; line < 5; line++ {
		if Board[0][line] == Board[1][line] && Board[0][line] == Board[2][line] && Board[0][line] == Board[3][line] && Board[0][line] == Board[4][line] {
			return Board[0][line]

		}
	}
	//判断对角线是否连成一线
	if Board[0][0] == Board[1][1] && Board[0][0] == Board[2][2] && Board[0][0] == Board[3][3] && Board[0][0] == Board[4][4] {
		return Board[0][0]
	}
	if Board[0][4] == Board[1][3] && Board[0][4] == Board[2][2] && Board[0][4] == Board[3][1] && Board[0][4] == Board[4][0] {
		return Board[0][4]
	}
	return " "
}
