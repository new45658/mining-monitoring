package shellParsing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// 根据 hostName 分组归纳信息
func ParseJobsInfo(jobs, workerHardwareInfo map[string]interface{}) interface{} {

	// job 根据 hostName 分组
	mapByHostName := make(map[string][]Task)
	for _, job := range jobs {
		sectorInfo := job.(map[string]interface{})

		if _, ok := sectorInfo["hostName"]; !ok {
			continue
		}
		hostName := sectorInfo["hostName"].(string)
		if taskList, ok := mapByHostName[hostName]; ok {
			taskList = append(taskList, mapToTask(sectorInfo))
			mapByHostName[hostName] = taskList
		} else {
			mapByHostName[hostName] = []Task{mapToTask(sectorInfo)}
		}
	}

	// 根据任务运行状态分组
	mapByState := make(map[string]interface{})
	for host, taskList := range mapByHostName {
		tHost := host
		taskMap := make(map[string]interface{})
		var currentQueue, pendQueue []Task

		for i := 0; i < len(taskList); i++ {
			task := taskList[i]
			if task.State == "running" {
				currentQueue = append(currentQueue, task)
			} else {
				pendQueue = append(pendQueue, task)
			}
		}
		taskMap["currentQueue"] = currentQueue
		taskMap["pendingQueue"] = pendQueue
		mapByState[tHost] = taskMap
	}

	// 把按状态分组，在按照任务类型分组
	mapByTask := make(map[string]interface{})
	for hostName, taskQueue := range mapByState {
		tHost := hostName
		result := make(map[string]interface{})
		tq := taskQueue.(map[string]interface{})
		for taskType, queue := range tq {
			q1 := queue.([]Task)
			param := tasksByType(q1)
			result[taskType] = param
		}
		mapByTask[tHost] = result
	}

	// 结合硬件信息
	var res []interface{}
	for hostName, hardwareInfo := range workerHardwareInfo {
		thInfo := hardwareInfo.(map[string]interface{})
		if hInfo, ok := mapByTask[hostName]; ok {
			info := hInfo.(map[string]interface{})
			param := mergeMaps(parseHardwareInfo(thInfo), info)
			res = append(res, param)
		}
	}
	return res

}

func parseHardwareInfo(src map[string]interface{}) map[string]interface{} {
	var gpus []interface{}
	if gpuList, ok := src["gpuInfo"]; ok {
		gpuMap := gpuList.(map[string]interface{})
		for _, gpu := range gpuMap {
			gpus = append(gpus, gpu)
		}
	}
	var netIOes []interface{}
	if netioMap, ok := src["netIO"]; ok {
		ioMap := netioMap.(map[string]interface{})
		for _, io := range ioMap {
			netIOes = append(netIOes, io)
		}
	}
	src["gpuInfo"] = gpus
	src["netIO"] = netIOes
	return src
}

// 整理worker任务信息
func mergeWorkerInfo(tasks []Task, hardwareList []HardwareInfo) interface{} {
	// 根据 hostName分组
	mapByHostName := make(map[string][]Task)
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		if taskList, ok := mapByHostName[task.HostName]; ok {
			taskList = append(taskList, task)
			mapByHostName[task.HostName] = taskList
		} else {
			mapByHostName[task.HostName] = []Task{task}
		}
	}
	// 根据任务运行状态分组
	mapByState := make(map[string]interface{})
	for host, taskList := range mapByHostName {
		tHost := host
		taskMap := make(map[string]interface{})
		var currentQueue, pendQueue []Task

		for i := 0; i < len(taskList); i++ {
			task := taskList[i]
			if task.State == "running" {
				currentQueue = append(currentQueue, task)
			} else {
				pendQueue = append(pendQueue, task)
			}
		}
		taskMap["currentQueue"] = currentQueue
		taskMap["pendingQueue"] = pendQueue
		mapByState[tHost] = taskMap
	}

	// 把按状态分组，在按照任务类型分组
	mapByTask := make(map[string]interface{})
	for hostName, taskQueue := range mapByState {
		tHost := hostName
		result := make(map[string]interface{})
		tq := taskQueue.(map[string]interface{})
		for taskType, queue := range tq {
			q1 := queue.([]Task)
			param := tasksByType(q1)
			result[taskType] = param
		}
		mapByTask[tHost] = result

	}

	// 结合硬件信息
	var res []interface{}
	for i := 0; i < len(hardwareList); i++ {
		hardware := hardwareList[i]
		if info, ok := mapByTask[hardware.HostName]; ok {
			tp := info.(map[string]interface{})
			toMap := structToMapByJson(&hardware)
			mapByTask[hardware.HostName] = mergeMaps(tp, toMap)
			res = append(res, mergeMaps(tp, toMap))
		}
	}
	return res
}

// 根据任务类型分组
func tasksByType(res []Task) map[string]interface{} {
	param := make(map[string]interface{})
	for i := 0; i < len(res); i++ {
		task := res[i]
		if taskList, ok := param[task.Task]; ok {
			tt := taskList.([]Task)
			taskList = append(tt, task)
			param[task.Task] = taskList
		} else {
			param[task.Task] = []Task{task}
		}
	}
	return param
}

func structToMapByJson(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	bytes, err := json.Marshal(obj)
	if err != nil {
		return m
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return m
	}
	return m
}

func structToMapByReflect(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return m
	}
	elem := reflect.ValueOf(obj).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[getHumpName(relType.Field(i).Name)] = elem.Field(i).Interface()
	}
	return m
}

func getHumpName(name string) string {
	if len(name) < 1 {
		return name
	}
	return fmt.Sprintf("%v%v", strings.ToLower(name[0:1]), name[1:])
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			tk := k
			tV := v
			result[tk] = tV
		}
	}
	return result
}

// 比较求两个map得差集,
func DiffMap(oldMap, newMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range newMap {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			if _, ok := oldMap[key]; !ok {
				result[key] = value
			} else {
				tempSrcMap := oldMap[key].(map[string]interface{})
				tempDestMap := value.(map[string]interface{})
				diffMap := DiffMap(tempSrcMap, tempDestMap)
				result[key] = diffMap
			}
		} else {
			if tv, ok := oldMap[key]; !ok || value != tv {
				result[key] = value
			}
		}
	}
	return result
}

func DeleteMapNull(src map[string]interface{}) map[string]interface{} {
	for key, value := range src {
		if value == nil {
			delete(src, key)
		} else {
			if reflect.TypeOf(value).Kind() == reflect.Map {
				tValue := value.(map[string]interface{})
				if len(tValue) == 0 {
					delete(src, key)
				} else {
					DeleteMapNull(tValue)
				}
			}
		}
	}
	return src
}

func getRegexValue(src [][]string) string {
	if len(src) == 0 || len(src[0]) == 0 {
		return "0"
	}
	return strings.ReplaceAll(src[0][1], " ", "")
}

func getRegexValueById(src [][]string, id int) string {
	if len(src) == 0 || len(src[0]) < id {
		return "0"
	}
	return strings.ReplaceAll(src[0][id], " ", "")
}

func mapToTask(src map[string]interface{}) Task {
	return Task{
		Id:       src["id"].(string),
		Sector:   src["sector"].(string),
		Worker:   src["worker"].(string),
		HostName: src["hostName"].(string),
		Task:     src["task"].(string),
		State:    src["state"].(string),
		Time:     src["time"].(string),
	}
}
