package main

import (
	"fmt"
	"log"
	"net"
)

//type user struct {
//	username string
//	password string
//	gorm.Model
//}
//
//var db *gorm.DB
//
//func init() {
//	DB, err := gorm.Open("mysql", "root:root@(localhost)/user?charset=utf8mb4&parseTime=True&loc=Local")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	DB.AutoMigrate()
//	db = DB
//}
//
////注册
//func Sign(mod *user) {
//	db.Create(&mod)
//}
//
////登录
//func Log(mod *user) int {
//	res := &user{}
//	db.Where("username= ?", mod.username).First(&res)
//	if mod.password != res.password {
//		return 0
//	}
//	return 1
//}
//
func sendmsg(conn net.Conn) {
	username := conn.LocalAddr().String()
	for {
		var input string
		// 接收输入消息，放到input变量中
		fmt.Scanln(&input)
		if len(input) > 0 {

			msg := username + " say:" + input

			_, err := conn.Write([]byte(msg))

			if err != nil {

				conn.Close()

				break

			}

		}

	}
}

func main() {
	//var (
	//	op1 int
	//	mod user
	//)
	//fmt.Println("1.登录")
	//fmt.Println("2.注册\n")
	//fmt.Scanln(&op1)
	//if op1 == 1 {
	//	fmt.Println("请输入用户名：\n")
	//	fmt.Scanln(&mod.username)
	//	fmt.Println("请输入密码\n")
	//	fmt.Scanln(&mod.password)
	//	res := Log(&mod)
	//	if res == 1 {
	//		run()
	//	}
	//}
	//if op1 == 2 {
	//	fmt.Println("请输入用户名：\n")
	//	fmt.Scanln(&mod.username)
	//	fmt.Println("请输入密码\n")
	//	fmt.Scanln(&mod.password)
	//	Sign(&mod)
	//	run()
	//}
	run()
}

func run() {

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	go sendmsg(conn)

	// 接收来自服务器端的广播消息

	buf := make([]byte, 1024)

	for {

		length, err := conn.Read(buf)

		if err != nil {

			log.Printf("recv server msg failed: %v\n", err)

			conn.Close()

			break

		}
		fmt.Println(string(buf[0:length]))

	}
}
