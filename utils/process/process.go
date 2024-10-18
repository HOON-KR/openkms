//go:build linux

// Copyright (c) 2024 Generic API Server All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package process 프로세스 유틸 패키지
*/
package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// IsProcessRunning 프로세스 동작 여부 확인
//
// Parameters:
//   - pid: PID
//
// Returns:
//   - bool: 동작중(true), 미동작(false)
func IsProcessRunning(pid int) bool {
	// 프로세스가 존재하는지 확인
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// 시그널 0을 보내 실제로 프로세스가 동작중인지 확인
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// Daemonize 프로세스 데몬화
//
// Return:
//   - error: 성공(nil), 실패(error)
func Daemonize() error {
	// 실행 파일 경로 획득
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %s", err)
	}

	// 자식 프로세스 생성
	cmd := exec.Command(exePath, os.Args[1:]...)
	cmd.Env = append(os.Environ(), "DAEMON=true") // 환경 변수 추가
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// 새로운 세션을 생성하고 부모 프로세스와 분리
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// 자식 프로세스 실행
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %s", err)
	}

	os.Exit(0) // 부모 프로세스 종료
	return nil
}
