package main

import (
	"fmt"
	"time"
)

// 复读机
func repeater(content string, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			fmt.Println(content)
		}
	}
}

// 定时器。精确到分
func timing(content string, h, m int) {
	// 先算出第一个输出的时间点
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, time.Local)
	if t.Before(now) {
		t = t.Add(24 * time.Hour)
	}

	timer := time.NewTimer(t.Sub(now))
	for {
		<-timer.C
		fmt.Println(content)
		timer.Reset(24 * time.Hour) // 直接重置就可以了
	}
}

func main() {
	go repeater("芜湖！起飞！", time.Hour)
	go timing("没有困难的工作，只有勇敢的打工人！", 2, 0)
	go timing("早安，打工人！", 8, 0)

	// 保持运行
	control := make(chan bool)
	<-control
}