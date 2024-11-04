package consistenthash

import (
	"bytes"
	"fmt"
	"hash"
	"hash/crc32"
	"sort"
	"sync"
)

// HashLoop is hashLoop
type HashLoop struct {
	lock       *sync.RWMutex
	hashPool   *sync.Pool
	bytesPool  *sync.Pool
	nodes      nodeSet
	virtualNum int
}

// New 创建一个 hash 环
// virtualNum 为虚拟节点数，为保证平衡性，每个真实节点被多个虚拟节点引用
func New(virtualNum int) *HashLoop {
	return &HashLoop{
		virtualNum: virtualNum,
		nodes:      make([]node, 0),
		lock:       &sync.RWMutex{},
		hashPool: &sync.Pool{
			New: func() any {
				return crc32.NewIEEE()
			},
		},
		bytesPool: &sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 10))
			},
		},
	}
}

// AddNodes 添加N个节点到环内
func (h *HashLoop) AddNodes(ips ...string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for _, ip := range ips {
		for i := 0; i < h.virtualNum; i++ {
			h.nodes = append(
				h.nodes,
				node{
					ip:  ip,
					num: h.hash(virtual(ip, i)),
				},
			)
		}
	}
	sort.Sort(h.nodes)
}

// DelNodes 删除N个环内节点
func (h *HashLoop) DelNodes(ips ...string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	var delIndexs []int
	for _, target := range ips {
		for i, node := range h.nodes {
			if node.ip == target {
				delIndexs = append(delIndexs, i)
			}
		}
	}
	var newns = make([]node, 0, len(h.nodes)-len(delIndexs))
	for i, node := range h.nodes {
		if !have(delIndexs, i) {
			newns = append(newns, node)
		}
	}
	h.nodes = newns
}

// Show 显示环内节点信息
func (h *HashLoop) Show() {
	h.lock.RLock()
	h.lock.RUnlock()
	if len(h.nodes) > 0 {
		fmt.Println("")

		lastNum := h.nodes[len(h.nodes)-1].num
		var useds nodeSet
		for _, n := range h.nodes {
			used := node{ip: n.ip, num: n.num - lastNum}
			useds = append(useds, used)
			// fmt.Println("IP", node.ip, "节点", node.num)
			// fmt.Println("*区间大小", node.num-lastNum)
			// fmt.Printf("*被使用概率%.3f%%\n", used*100)
			// fmt.Println("")
			lastNum = n.num
		}

		sort.Sort(useds)
		fmt.Println("总节点数:", len(useds))
		if len(useds) > 2 {
			fmt.Printf("Min 3rd probability: %.3f%%(%s) %.3f%%(%s) %.3f%%(%s)\n",
				getPb(useds[0]), useds[0].ip,
				getPb(useds[1]), useds[1].ip,
				getPb(useds[2]), useds[2].ip,
			)
			fmt.Printf("Max 3rd probability: %.3f%%(%s) %.3f%%(%s) %.3f%%(%s)\n",
				getPb(useds[len(useds)-1]), useds[len(useds)-1].ip,
				getPb(useds[len(useds)-2]), useds[len(useds)-2].ip,
				getPb(useds[len(useds)-3]), useds[len(useds)-3].ip,
			)
		}
		fmt.Printf("Avg probability: %.3f%%\n", float64(100)/float64(len(useds)))
	}
}

// SelectNode 通过一致性hash寻找响应节点
func (h *HashLoop) SelectNode(sourceIP string) string {

	h.lock.RLock()
	defer h.lock.RUnlock()
	// 空环，长度为0的边界条件
	if len(h.nodes) == 0 {
		return ""
	}

	buf := h.bytesPool.Get().(*bytes.Buffer)
	buf.WriteString(sourceIP)
	index := h.hash(buf.Bytes())
	buf.Reset()
	h.bytesPool.Put(buf)

	for _, node := range h.nodes {
		// todo:二分查询
		if node.num > index {
			return node.ip
		}
	}

	// 遍历环结束仍未发现大于目标的节点，则返回第一个节点
	return h.nodes[0].ip
}

func (h *HashLoop) hash(data []byte) uint32 {
	c := h.hashPool.Get().(hash.Hash32)
	c.Write(data)
	res := c.Sum32()
	c.Reset()
	h.hashPool.Put(c)
	return res
}
