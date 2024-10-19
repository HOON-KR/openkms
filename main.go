//go:build linux

// Copyright (c) 2024 Generic API Server All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package main 메인 패키지
*/
package main

import (
	"flag"
	"fmt"
	"openkms/config"
	"openkms/utils/file"
	"openkms/utils/log"
	"openkms/utils/process"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

// options 명령행 옵션 정보 구조체
type options struct {
	version bool
	help    bool
}

// setOptions 옵션 값 설정
func (o *options) setOptions() {
	flag.BoolVar(&o.version, "v", false, "Print version")
	flag.BoolVar(&o.help, "h", false, "Print help")
}

// getVersion 버전 정보 출력
//
// Returns:
//   - string: 버전 정보
func (o *options) getVersion() string {
	return fmt.Sprintf("%s version %s", config.ModuleName, Version)
}

// usage 사용법 출력
func (o *options) usage() {
	fmt.Println(o.getVersion())
	fmt.Println("Build Date:", BuildDate)
	fmt.Println("Command: start | stop")
	flag.Usage()
}

// processOption 명령행 옵션 처리
func (o *options) processOption() {
	if o.version {
		fmt.Println(o.getVersion())
		os.Exit(0)
	}

	if o.help {
		o.usage()
		os.Exit(0)
	}
}

// main 메인 함수
func main() {
	var option options

	option.setOptions() // 명령행 옵션 설정

	if len(os.Args) <= 1 {
		option.usage()
		os.Exit(0)
	}

	flag.Parse() // 명령행 옵션 파싱

	// 작업 경로를 현재 실행 파일의 경로로 변경
	if err := changeWorkDir(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	option.processOption() // 명령행 옵션 처리

	// 동작 명령어 체크
	switch os.Args[1] {
	case "start":
		// 이미 동작중인 프로세스가 존재하는지 확인
		if isProcessRunWithPidFile(config.PidFilePath) {
			fmt.Println("There is already a working process")
			os.Exit(0)
		}
	case "stop":
		// 프로세스 종료 시그널(SIGTERM) 전송
		pid, err := stopProcess(config.PidFilePath)
		if err != nil {
			fmt.Println(err)
		} else if pid != 0 {
			fmt.Printf("Stop %s process (pid: %d)\n", config.ModuleName, pid)
		}
		os.Exit(0)
	default:
		option.usage()
		os.Exit(0)
	}

	sigChan := make(chan os.Signal, 1)
	stopChan := make(chan bool)

	setupSignal(sigChan) // 시그널 설정

	// 환경 변수를 체크하여 데몬 프로세스인지 확인
	if os.Getenv("DAEMON") != "true" {
		// 프로세스 데몬화
		// 데몬화 성공 시 함수 내부에서 프로세스 종료
		err := process.Daemonize()
		// 프로세스 데몬화 실패
		fmt.Println(err)
		os.Exit(1)
	}

	// 종료 시그널 처리
	go func() {
		sig := <-sigChan
		log.LogInfo("Receive SIGNAL: %d", sig)
		stopChan <- true
	}()

	initialization() // 초기화
	defer func() {
		finalization() // 종료 전 작업 정리
	}()

	// 데몬 프로세스인 경우 PID를 파일에 기록
	err := file.WriteTextFile[int](config.PidFilePath, os.Getpid())
	if err != nil {
		log.LogWarn("%s", err)
		return
	}

	<-stopChan // 종료 대기
}

// changeWorkDir 작업 경로를 현재 실행 파일의 경로로 변경
//
// Returns:
//   - error: 성공(nil), 실패(error)
func changeWorkDir() error {
	// 실행 파일 경로 획득
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %s", err)
	}

	exeDir := filepath.Dir(exePath) // 실행 파일의 디렉터리 경로 추출

	// 작업 디렉토리 변경
	err = os.Chdir(exeDir)
	if err != nil {
		return fmt.Errorf("error changing working directory: %s", err)
	}

	return nil
}

// isProcessRunWithPidFile 파일에서 PID를 읽고, 해당 PID를 가진 프로세스가 동작 중인지 확인
//
// Parameters:
//   - pidFilePath: PID 파일 경로
//
// Returns:
//   - bool: 동작중(true), 미동작(false)
func isProcessRunWithPidFile(pidFilePath string) bool {
	// PID 파일 읽기
	pidBytes, err := os.ReadFile(pidFilePath)
	if err != nil {
		return false
	}

	// 파일에서 읽은 PID를 정수로 변환
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return false
	}

	// 프로세스가 동작중인지 확인
	return process.IsProcessRunning(pid)
}

func stopProcess(pidFilePath string) (int, error) {
	// PID 파일 읽기
	pidBytes, err := os.ReadFile(pidFilePath)
	if err != nil {
		return 0, nil
	}

	// 파일에서 읽은 PID를 정수로 변환
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return 0, nil
	}

	// 프로세스가 존재하는지 확인
	process, err := os.FindProcess(pid)
	if err != nil {
		return 0, nil
	}

	// 시그널 0을 보내 실제로 프로세스가 동작중인지 확인
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return 0, nil
	}

	// SIGTERM 시그널을 보내 프로세스 종료
	err = process.Signal(syscall.Signal(syscall.SIGTERM))
	if err != nil {
		return 0, fmt.Errorf("failed to send signal (pid: %d): %s", pid, err)
	}

	os.Remove(pidFilePath) // PID 파일 삭제

	return pid, nil
}

// initialization 초기화 함수
func initialization() {
	file.MakeDirectory("var") // var 디렉터리 생성
	file.MakeDirectory("log") // log 디렉터리 생성
	log.InitLogger()          // 로거 초기화
}

// finalization 모듈 종료 시 작업 정리 함수
func finalization() {
	log.FinalizeLog() // 로그 자원 정리
}

// setupSignal 시그널 설정
//
// Parameters:
//   - sigChan: 시그널을 수신할 채널
func setupSignal(sigChan chan os.Signal) {
	signal.Ignore(syscall.SIGABRT, syscall.SIGALRM, syscall.SIGHUP, syscall.SIGTSTP,
		syscall.SIGILL, syscall.SIGPROF, syscall.SIGQUIT, syscall.SIGVTALRM)

	signal.Notify(sigChan, syscall.SIGTERM)
}
