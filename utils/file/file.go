//go:build linux

// Copyright (c) 2024 Generic API Server All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package file 파일 유틸 패키지
*/
package file

import (
	"fmt"
	"os"
)

// WriteTextFile 제네릭 함수를 사용하여 데이터를 텍스트 파일에 기록하는 함수
//
// Parameters:
//   - filePath: 파일 경로
//   - data: 파일에 기록할 데이터
//
// Returns:
//   - error: 성공(nil), 실패(error)
func WriteTextFile[T any](filePath string, data T) error {
	// 파일 열기 (없으면 생성, 있으면 덮어쓰기)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	// 데이터를 문자열로 변환하여 파일에 기록
	_, err = fmt.Fprint(file, data)
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}

// MakeDirectory 디렉터리 생성 함수
//
// Parameters:
//   - dirPath: 디렉터리 경로
//
// Returns:
//   - error: 성공(nil), 실패(error)
func MakeDirectory(dirPath string) error {
	// 디렉터리 상태 정보 획득
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// 디렉터리가 존재하지 않으면 생성
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %s", err)
		}
	} else if err != nil {
		// 에러 발생
		return fmt.Errorf("error checking directory: %s", err)
	}

	return nil
}
