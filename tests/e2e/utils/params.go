package utils

import (
	"os"
	"sync"

	"github.com/ava-labs/avalanche-network-runner/client"

	"gopkg.in/yaml.v2"
)

// ClusterInfo represents the local cluster information.
type ClusterInfo struct {
	URIs                  []string `json:"uris"`
	Endpoint              string   `json:"endpoint"`
	PID                   int      `json:"pid"`
	LogsDir               string   `json:"logsDir"`
	SubnetEVMRPCEndpoints []string `json:"subnetEVMRPCEndpoints"`
}

const fsModeWrite = 0o600

func (ci ClusterInfo) Save(p string) error {
	ob, err := yaml.Marshal(ci)
	if err != nil {
		return err
	}
	return os.WriteFile(p, ob, fsModeWrite)
}

var (
	mu sync.RWMutex

	cli client.Client

	outputFile string
	pluginDir  string

	// executable path for "avalanchego"
	execPath      string
	vmGenesisPath string

	skipNetworkRunnerShutdown bool

	clusterInfo ClusterInfo

	contractsFoundryDir string
)

func SetClient(c client.Client) {
	mu.Lock()
	cli = c
	mu.Unlock()
}

func GetClient() client.Client {
	mu.RLock()
	c := cli
	mu.RUnlock()
	return c
}

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

// Sets the executable path for "avalanchego".
func SetExecPath(p string) {
	mu.Lock()
	execPath = p
	mu.Unlock()
}

// Loads the executable path for "avalanchego".
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

func SetSkipNetworkRunnerShutdown(b bool) {
	mu.Lock()
	skipNetworkRunnerShutdown = b
	mu.Unlock()
}

func GetSkipNetworkRunnerShutdown() bool {
	mu.RLock()
	b := skipNetworkRunnerShutdown
	mu.RUnlock()
	return b
}

func SetClusterInfo(c ClusterInfo) {
	mu.Lock()
	clusterInfo = c
	mu.Unlock()
}

func GetClusterInfo() ClusterInfo {
	mu.RLock()
	c := clusterInfo
	mu.RUnlock()
	return c
}

func SetContractsFoundryDir(dir string) {
	mu.Lock()
	contractsFoundryDir = dir
	mu.Unlock()
}

func GetContractsFoundryDir() string {
	mu.RLock()
	dir := contractsFoundryDir
	mu.RUnlock()
	return dir
}
