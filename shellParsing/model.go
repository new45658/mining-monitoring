package shellParsing

type WorkerInfo struct {
	HostName string
	IP       string
}

type Workers []WorkerInfo

type P map[string]interface{}

type MinerInfo struct {
	minerId          string   // minerId
	MinerBalance     string   // miner余额
	postBalance      string   // post余额
	workerBalance    string   // worker余额
	pledgeBalance    string   // 抵押
	totalMessages    int      // 消息总数
	rawBytePower     int      // 字节算力
	adjustedPower    int      // 原值算力
	effectivePower   int      // 有效算力
	totalSectors     int      // 总扇区数
	effectiveSectors int      // 有效扇区
	errorSectors     int      // 错误扇区
	recoverySectors  int      // 恢复中扇区
	deletedSectors   int      // 删除扇区
	failSectors      int      // 失败扇区
	workerInfo       []Worker // worker信息
	timestamp        int      // 此次统计时间
}

// ID        Sector  Worker    Hostname       Task  State        Time
//c71e05fc  8598    74d84e37  ya_amd_node36  PC1   running      2h12m29.5s
//
type Task struct {
	Id       string
	Sector   string //扇区Id
	Worker   string
	HostName string
	Task     string //任务类型
	State    string // 任务状态
	Time     string // 耗时
}




type Worker struct {
	hostname     string
	currentQueue []Task  // 当前任务
	pendingQueue []Task  // 队列中任务
	cpuTemper    float32 // cpu问题
	cpuLoad      float32 // cupu负载
	gpuTemper    float32 // gpu温度
	gpuLoad      float32 // gpu负载
	memory       string  // 内存信息
	useDisk      float32 // 磁盘使用率
	diskRW       string  //磁盘IO
	netRW        string  //网络IO

}
