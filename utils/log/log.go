//go:build linux

// Copyright (c) 2024 Generic API Server All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package log 로그 유틸 패키지
*/
package log

import (
	"fmt"
	"openkms/config"
	"openkms/utils/file"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 로그 정보 구조체
type Logger struct {
	logger *zap.SugaredLogger
}

var logger Logger

// init 패키지 초기화
func init() {
	file.MakeDirectory("log") // log 디렉터리 생성
	initLogger()              // 로거 초기화
}

// initLogger 로거 초기화
func initLogger() {
	// lumberjack 로테이션 설정
	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.LogFilePath, // 로그 파일 경로
		MaxSize:    100,                // 최대 크기(MB)
		MaxBackups: 10,                 // 보관할 백업 파일 수
		MaxAge:     30,                 // 보관할 최대 일수
		Compress:   true,               // 압축 여부
	})

	// 로그 출력 포맷 설정
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",                                               // 시간 키
		LevelKey:         "level",                                              // 레벨 키
		NameKey:          "logger",                                             // 로거 이름 키
		CallerKey:        "caller",                                             // 호출자 키
		MessageKey:       "msg",                                                // 메시지 키
		StacktraceKey:    "stacktrace",                                         // 스택트레이스 키
		LineEnding:       zapcore.DefaultLineEnding,                            // 줄 끝
		EncodeLevel:      customCapitalLevelEncoder,                            // 레벨 대문자 인코딩
		EncodeTime:       zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05]"), // 시간 포맷 지정
		EncodeDuration:   zapcore.StringDurationEncoder,                        // 지속시간 인코딩
		EncodeCaller:     zapcore.ShortCallerEncoder,                           // 호출자 인코딩
		ConsoleSeparator: " ",                                                  // 로그 필드 구분자
	}

	// 코어 생성
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // human-readable 형식의 출력
		logWriter,                                // lumberjack과 연동
		zapcore.DebugLevel,                       // 로그 레벨 설정
	)

	// <호출자 정보 추가>
	// 1단계 스택 깊이 스킵
	// ERROR 레벨 이상에서 스택 트레이스 추가
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel))
	// SugaredLogger로 변환하여 가변 인자 지원
	logger.logger = zapLogger.Sugar()
}

// FinalizeLog 로그 자원 정리
func FinalizeLog() {
	logger.logger.Sync() // 프로그램 종료 시 남은 로그가 모두 기록되도록 함
}

// customCapitalLevelEncoder 로그 레벨 인코더 커스텀 설정
//
// Parameters:
//   - level: 로그 레벨
//   - enc: 인코더
func customCapitalLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// LogInfo 정보 로그를 출력하는 함수 (가변 인자 처리)
//
// Parameters:
//   - format: 로그 포맷
//   - args: 가변 인자
func LogInfo(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.logger.Info(message)
}

// LogWarn 경고 로그를 출력하는 함수 (가변 인자 처리)
//
// Parameters:
//   - format: 로그 포맷
//   - args: 가변 인자
func LogWarn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.logger.Warn(message)
}

// LogError 에러 로그를 출력하는 함수 (가변 인자 처리)
//
// Parameters:
//   - format: 로그 포맷
//   - args: 가변 인자
func LogError(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.logger.Error(message)
}

// LogDebug 디버그 로그를 출력하는 함수 (가변 인자 처리)
//
// Parameters:
//   - format: 로그 포맷
//   - args: 가변 인자
func LogDebug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	logger.logger.Debug(message)
}
