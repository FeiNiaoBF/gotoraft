package raft

import (
	"fmt"
	"math/rand"
	"time"
)

// DPrintf 用于调试日志输出
func DPrintf(format string, a ...interface{}) {
	debug := false
	if debug {
		fmt.Printf(format+"\n", a...)
	}
}

// randomElectionTimeout 生成一个随机的选举超时时间
func randomElectionTimeout(min, max, jitter time.Duration) time.Duration {
	// 生成一个在 [min, max) 范围内的随机时间
	timeout := min + time.Duration(rand.Int63n(int64(max-min)))
	// 添加一个小的随机偏移量，避免选举冲突
	jitterDuration := time.Duration(rand.Int63n(int64(jitter)))
	return timeout + jitterDuration
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
