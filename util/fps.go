package util

import (
	"fmt"
	"time"
)

// FPSManager 控制帧速率的结构体
type FPSManager struct {
	frameDuration time.Duration
	lastTime      time.Time
}

// NewFPSManager 创建一个新的帧速率控制器
func NewFPSManager(fps int) *FPSManager {
	return &FPSManager{
		frameDuration: time.Second / time.Duration(fps),
		lastTime:      time.Now(),
	}
}

// Wait 等待直到下一帧应该被绘制时返回
func (frc *FPSManager) Wait() {
	now := time.Now()
	elapsed := now.Sub(frc.lastTime)
	if elapsed < frc.frameDuration {
		fmt.Println(frc.frameDuration - elapsed)
		time.Sleep(frc.frameDuration - elapsed)
	}
	frc.lastTime = now
}
