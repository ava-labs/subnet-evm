package utils

import "sync"

var (
	mu sync.RWMutex

	outputFile string
	pluginDir  string

	execPath string

	vmGenesisPath string

	mode string
)

func SetOutputFile(filepath string) {
	mu.Lock()
	outputFile = filepath
	mu.Unlock()
}

func GetOutputPath() string {
	mu.RLock()
	e := outputFile
	mu.RUnlock()
	return e
}

func SetExecPath(p string) {
	mu.Lock()
	execPath = p
	mu.Unlock()
}

func GetExecPath() string {
	mu.RLock()
	e := execPath
	mu.RUnlock()
	return e
}

func SetPluginDir(dir string) {
	mu.Lock()
	pluginDir = dir
	mu.Unlock()
}

func GetPluginDir() string {
	mu.RLock()
	p := pluginDir
	mu.RUnlock()
	return p
}

func SetVmGenesisPath(p string) {
	mu.Lock()
	vmGenesisPath = p
	mu.Unlock()
}

func GetVmGenesisPath() string {
	mu.RLock()
	p := vmGenesisPath
	mu.RUnlock()
	return p
}

func SetMode(m string) {
	mu.Lock()
	mode = m
	mu.Unlock()
}

func GetMode() string {
	mu.RLock()
	m := mode
	mu.RUnlock()
	return m
}
