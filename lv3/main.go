package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const filePath = "./users.data"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userHash map[string]string

type Checker struct {
	uh userHash	// 用户信息
	registerUsers []User // 注册了但未保存的用户
}

func (c *Checker) SignIn() {
	fmt.Println("请输入用户名和密码")
	var username, password string
	fmt.Scan(&username, &password)
	if _, ok := c.uh[username]; !ok {
		fmt.Println("查无此人")
		return
	}
	if c.uh[username] != password {
		fmt.Println("用户名密码错误")
		return
	}
	
	fmt.Println("登录成功")
}

func (c *Checker) SignUp() {
	fmt.Println("请输入用户名")
	var username, password string
	fmt.Scan(&username)
	if _, ok := c.uh[username]; ok {
		fmt.Println("用户名已被占用")
		return
	}
	fmt.Println("请输入密码")
	for {
		fmt.Scan(&password)
		if len(password) >= 6 {
			break
		}
		fmt.Println("密码长度应大于六位，请重新输入")
	}

	// 先写入缓存，再异步写入文件
	c.registerUsers = append(c.registerUsers, User{
		Username: username,
		Password: password,
	})
	if len(c.registerUsers) > 10 {
		go c.Save()
	}
	c.uh[username] = password
}

func (c *Checker) Save() {
	fail := saveUsers(c.registerUsers)
	c.registerUsers = fail
}

func initUsers() (userHash, error){
	f, err := os.OpenFile(filePath, os.O_CREATE | os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer f.Close()

	uh := make(userHash)
	reader := bufio.NewReader(f)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return nil, err
		}
		var user User
		err = json.Unmarshal(buf, &user)
		if err != nil {
			fmt.Println(err)
			continue
		}
		uh[user.Username] = user.Password
	}
	return uh, nil
}

func saveUsers(users []User) (fail []User){
	// 以追加的方式写入文件
	f, err := os.OpenFile(filePath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, user := range users{
		buf, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
			fail = append(fail, user)
			continue
		}
		n, err := writer.Write(append(buf, byte('\n')))
		if err != nil {
			fmt.Println(n, err)
			fail = append(fail, user)
			continue
		}
	}
	writer.Flush()
	return
}

func showList() {
	fmt.Println("请选择操作：")
	fmt.Println("1、登录")
	fmt.Println("2、注册")
	fmt.Println("3、退出")
}

func main() {
	checker := Checker{}
	var err error
	checker.uh, err = initUsers()
	if err != nil {
		return
	}
	

	var opt int
	for {
		showList()
		_, err := fmt.Scanln(&opt)
		if err != nil || opt < 1 || opt > 3 {
			fmt.Println("请输入正确的操作序号")
			continue
		}

		switch opt {
		case 1:
			checker.SignIn()
		case 2:
			checker.SignUp()
		case 3:
			checker.Save()
			return
		}
	}
}