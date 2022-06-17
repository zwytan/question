package cache

import (
	"question/util"
	"runtime"
	"sync"
	"time"
)

type Cache interface {
	//size 是⼀个字符串。⽀持以下参数: 1KB，100KB，1MB，2MB，1GB 等
	SetMaxMemory(size string) bool
	// 设置⼀个缓存项，并且在expire时间之后过期
	Set(key string, val interface{}, expire time.Duration)
	// 获取⼀个值
	Get(key string) (interface{}, bool)
	// 删除⼀个值
	Del(key string) bool
	// 检测⼀个值 是否存在
	Exists(key string) bool
	// 情况所有值
	Flush() bool
	// 返回所有的key 多少
	Keys() int64
	// 返回容量大小
	Size() uint64
}

type item struct {
	v interface{} // 值
	e int64       // 过期时间
}

type cache struct {
	mu     sync.RWMutex
	data   map[string]item
	maxMen uint64
	mem    uint64
}

var _ Cache = (*cache)(nil)

//创建实例
func NewCache() *cache {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	cache := &cache{
		data: make(map[string]item),
		mem:  ms.Alloc,
	}
	//清除过期数据
	tick := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-tick.C
			cache.timeDel()
		}
	}()
	return cache
}

func (c *cache) timeDel() {
	t := time.Now().Unix()
	for k, v := range c.data {
		if t > v.e {
			c.Del(k)
		}
	}
}

//设置最大内存
func (c *cache) SetMaxMemory(size string) bool {
	bytes, err := util.ToBytes(size)
	if err != nil {
		panic(err.Error())
	}
	c.maxMen = bytes
	return true
}

//设置键值
func (c *cache) Set(key string, val interface{}, expire time.Duration) {
	if !c.checkMem(1) {
		panic("内存已满")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = item{
		v: val,
		e: time.Now().Add(expire).Unix(),
	}

}

//获取1个key
func (c *cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if ok {
		if time.Now().Unix() <= item.e {
			return item.v, ok
		} else {
			delete(c.data, key)
		}
	}
	return nil, ok
}

//删除1个key
func (c *cache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return true
}

//查看key是否存在
func (c *cache) Exists(key string) bool {
	item, ok := c.data[key]
	if ok && time.Now().Unix() > item.e {
		return false
	}
	return ok
}

//清空
func (c *cache) Flush() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = nil
	runtime.GC()
	return true
}

//获取key数量
func (c *cache) Keys() int64 {
	return int64(len(c.data))
}

// 获取当前内存容量
func (c *cache) Size() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.Alloc - c.mem
}

//检测容量是否已满
func (c *cache) checkMem(dataSize int) bool {
	if c.maxMen >= (c.Size() + uint64(dataSize)) {
		return true
	}
	return false
}
