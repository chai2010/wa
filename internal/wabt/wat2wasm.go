// 版权 @2023 凹语言 作者。保留所有权利。

package wabt

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"wa-lang.org/wa/internal/config"
	"wa-lang.org/wa/internal/logger"
)

var muWabt sync.Mutex
var wat2wasmPath string

const Wat2WasmName = "wa.wat2wasm.exe"

func init() {
	// 1. exe 同级目录存在 wat2wasm ?
	wat2wasmPath = filepath.Join(curExeDir(), Wat2WasmName)
	if exeExists(wat2wasmPath) {
		return
	}

	// 2. 当前目录存在 wat2wasm ?
	cwd, _ := os.Getwd()
	wat2wasmPath = filepath.Join(cwd, Wat2WasmName)
	if exeExists(wat2wasmPath) {
		return
	}

	// 3. 本地系统存在 wat2wasm ?
	if s, _ := exec.LookPath(Wat2WasmName); s != "" {
		wat2wasmPath = s
		return
	}

	// 4. 查找C盘根目录(为了简化一些测试环境)
	if runtime.GOOS == "windows" {
		wat2wasmPath = "c:/" + Wat2WasmName
		if exeExists(wat2wasmPath) {
			return
		}
	} else {
		wat2wasmPath = "/usr/local/bin/" + Wat2WasmName
		if exeExists(wat2wasmPath) {
			return
		}
	}
}

func Wat2Wasm(watBytes []byte) (wasmBytes []byte, err error) {
	muWabt.Lock()
	defer muWabt.Unlock()

	if wat2wasmPath == "" {
		logger.Tracef(&config.EnableTrace_app, "wat2wasm not found")
		return nil, errors.New("wat2wasm not found")
	}

	var bufStdout bytes.Buffer
	var bufStderr bytes.Buffer

	// wat2wasm - --output=-
	var args = []string{"-", "--output=-"}
	if config.DebugMode {
		args = append(args, "--debug-names")
	}

	cmd := exec.Command(wat2wasmPath, args...)
	cmd.Stdin = bytes.NewReader(watBytes)
	cmd.Stdout = &bufStdout
	cmd.Stderr = &bufStderr

	err = cmd.Run()
	wasmBytes = bufStdout.Bytes()

	if err != nil && bufStderr.Len() > 0 {
		err = errors.New(bufStderr.String())
	}

	return
}

// exe 文件存在
func exeExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	if !fi.Mode().IsRegular() {
		return false
	}
	return true
}

// 当前执行程序所在目录
func curExeDir() string {
	s, err := os.Executable()
	if err != nil {
		logger.Panicf("os.Executable() failed: %+v", err)
	}
	return filepath.Dir(s)
}
