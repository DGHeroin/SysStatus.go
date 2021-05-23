package SysStatus

import (
    "bytes"
    "fmt"
    "math"
    "runtime"
    "time"
)

var (
    Status   St
    initTime = time.Now()
)

func init() {
    UpdateSystemStatus()
}

type St struct {
    Uptime       string // 运行时间
    NumGoroutine int    // 当前 Goroutines 数量

    // General statistics.
    MemAllocated string // 当前内存使用量
    MemTotal     string // 所有被分配的内存
    MemSys       string // 内存占用量
    Lookups      uint64 // 指针查找次数
    MemMallocs   uint64 // 内存分配次数
    MemFrees     uint64 // 内存释放次数

    // Main allocation heap statistics.
    HeapAlloc    string // 当前 Heap 内存使用
    HeapSys      string // Heap 内存占用量
    HeapIdle     string // Heap 内存空闲量
    HeapInuse    string // 正在使用的 Heap 内存
    HeapReleased string // 被释放的 Heap 内存
    HeapObjects  uint64 // Heap 对象数量

    // Low-level fixed-size structure allocator statistics.
    //	Inuse is bytes used now.
    //	Sys is bytes obtained from system.
    StackInuse  string // 启动 Stack 使用量
    StackSys    string // 被分配的 Stack 内存
    MSpanInuse  string // MSpan 结构内存使用量
    MSpanSys    string // 被分配的 MSpan 结构内存
    MCacheInuse string // MCache 结构内存使用量
    MCacheSys   string // 被分配的 MCache 结构内存
    BuckHashSys string // 被分配的剖析哈希表内存
    GCSys       string // 被分配的 GC 元数据内存
    OtherSys    string // 其它被分配的系统内存

    // Garbage collector statistics.
    NextGC       string // 下次 GC 内存回收量
    LastGC       string // 距离上次 GC 时间
    PauseTotalNs string // GC 暂停时间总量
    PauseNs      string // 上次 GC 暂停时间
    NumGC        uint32 // GC 执行次数
}

func UpdateSystemStatus() {
    m := new(runtime.MemStats)
    runtime.ReadMemStats(m)
    Status.NumGoroutine = runtime.NumGoroutine()

    Status.MemAllocated = getByteSize(int64(m.Alloc))
    Status.MemTotal = getByteSize(int64(m.TotalAlloc))
    Status.MemSys = getByteSize(int64(m.Sys))
    Status.Lookups = m.Lookups
    Status.MemMallocs = m.Mallocs
    Status.MemFrees = m.Frees

    Status.HeapAlloc = getByteSize(int64(m.HeapAlloc))
    Status.HeapSys = getByteSize(int64(m.HeapSys))
    Status.HeapIdle = getByteSize(int64(m.HeapIdle))
    Status.HeapInuse = getByteSize(int64(m.HeapInuse))
    Status.HeapReleased = getByteSize(int64(m.HeapReleased))
    Status.HeapObjects = m.HeapObjects

    Status.StackInuse = getByteSize(int64(m.StackInuse))
    Status.StackSys = getByteSize(int64(m.StackSys))
    Status.MSpanInuse = getByteSize(int64(m.MSpanInuse))
    Status.MSpanSys = getByteSize(int64(m.MSpanSys))
    Status.MCacheInuse = getByteSize(int64(m.MCacheInuse))
    Status.MCacheSys = getByteSize(int64(m.MCacheSys))
    Status.BuckHashSys = getByteSize(int64(m.BuckHashSys))
    Status.GCSys = getByteSize(int64(m.GCSys))
    Status.OtherSys = getByteSize(int64(m.OtherSys))

    Status.NextGC = getByteSize(int64(m.NextGC))
    Status.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
    Status.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
    Status.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
    Status.NumGC = m.NumGC
}

func getByteSize(s int64) string {
    sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
    return humanateBytes(uint64(s), 1024, sizes)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
    if s < 10 {
        return fmt.Sprintf("%d B", s)
    }
    e := math.Floor(logn(float64(s), base))
    suffix := sizes[int(e)]
    val := float64(s) / math.Pow(base, math.Floor(e))
    f := "%.0f"
    if val < 10 {
        f = "%.1f"
    }

    return fmt.Sprintf(f+" %s", val, suffix)
}

func logn(n, b float64) float64 {
    return math.Log(n) / math.Log(b)
}

func (s St) String() string {
    buf := bytes.NewBufferString("SysStatus:\n")
    buf.WriteString(fmt.Sprintf("当前 Goroutines 数量:     %v\n", Status.NumGoroutine))
    buf.WriteString(fmt.Sprintf("当前内存使用量:           %v\n", Status.MemAllocated))
    buf.WriteString(fmt.Sprintf("所有被分配的内存:         %v\n", Status.MemTotal))
    buf.WriteString(fmt.Sprintf("内存占用量:               %v\n", Status.MemSys))
    buf.WriteString(fmt.Sprintf("指针查找次数:             %v\n", Status.Lookups))
    buf.WriteString(fmt.Sprintf("内存分配次数:             %v\n", Status.MemMallocs))
    buf.WriteString(fmt.Sprintf("内存释放次数:             %v\n", Status.MemFrees))
    buf.WriteString(fmt.Sprintf("当前 Heap 内存使用量:     %v\n", Status.HeapAlloc))
    buf.WriteString(fmt.Sprintf("Heap 内存占用:            %v\n", Status.HeapSys))
    buf.WriteString(fmt.Sprintf("Heap 内存空闲量:          %v\n", Status.HeapIdle))
    buf.WriteString(fmt.Sprintf("正在使用的 Heap 内存:     %v\n", Status.HeapInuse))
    buf.WriteString(fmt.Sprintf("被释放的 Heap 内存:       %v\n", Status.HeapReleased))
    buf.WriteString(fmt.Sprintf("Heap 对象数量:            %v\n", Status.HeapObjects))
    buf.WriteString(fmt.Sprintf("启动 Stack 使用量:        %v\n", Status.StackInuse))
    buf.WriteString(fmt.Sprintf("被分配的 Stack 内存:      %v\n", Status.StackSys))
    buf.WriteString(fmt.Sprintf("MSpan 结构内存使用量:     %v\n", Status.MSpanInuse))
    buf.WriteString(fmt.Sprintf("被分配的 MSpan 结构内存:  %v\n", Status.MSpanSys))
    buf.WriteString(fmt.Sprintf("MCache 结构内存使用量:    %v\n", Status.MCacheInuse))
    buf.WriteString(fmt.Sprintf("被分配的 MCache 结构内存: %v\n", Status.MCacheSys))
    buf.WriteString(fmt.Sprintf("被分配的剖析哈希表内存:   %v\n", Status.BuckHashSys))
    buf.WriteString(fmt.Sprintf("被分配的 GC 元数据内存:   %v\n", Status.GCSys))
    buf.WriteString(fmt.Sprintf("其它被分配的系统内存:     %v\n", Status.OtherSys))
    buf.WriteString(fmt.Sprintf("下次 GC 内存回收量:       %v\n", Status.NextGC))
    buf.WriteString(fmt.Sprintf("距离上次 GC 时间:         %v\n", Status.LastGC))
    buf.WriteString(fmt.Sprintf("GC 暂停时间总量:          %v\n", Status.PauseTotalNs))
    buf.WriteString(fmt.Sprintf("上次 GC 暂停时间:         %v\n", Status.PauseNs))
    buf.WriteString(fmt.Sprintf("GC 执行次数:              %v\n", Status.NumGC))

    return buf.String()
}
