package shellParsing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mining-monitoring/log"
	"time"
)

type Manager struct {
	currentInfo map[string]interface{}
	shellParse  *ShellParse
	Workers     []WorkerInfo
}

func (m *Manager) GetCurrentMinerInfo() interface{} {
	return m.currentInfo
}

func (m *Manager) DoShell() (interface{}, error) {
	if e := recover(); e != nil {
		log.Error("doShell error: %v ", e)
	}
	taskInfo, err := m.shellParse.getTaskInfo()
	if err != nil {
		return nil, err
	}
	// todo
	return taskInfo, nil

}

func (m *Manager) Run(obj chan interface{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Debug("start timer get minerInfo ")
			result, err := m.DoShell()
			if err != nil {
				fmt.Printf("doShell error %v \n", err)
				continue
			}
			obj <- result
		default:

		}
	}
}

func NewManager(path string) (*Manager, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workerhost error %v \n", err)
	}
	var workers []WorkerInfo
	err = json.Unmarshal(data, &workers)
	if err != nil {
		return nil, fmt.Errorf("parse json error: %v \n", err)
	}

	return &Manager{
		currentInfo: map[string]interface{}{},
		shellParse:  NewShellParse(workers),
		Workers:     workers,
	}, nil
}
