package lru

import "container/list"

type Cache struct {
	maxBytes  int64                     // 允许使用的最大内存, 单位为 bytes
	nBytes    int64                     // 当前已使用的内存
	deQueue   *list.List                // 双向链表
	cache     map[string]*list.Element  // 键为字符串, 值为双向链表对应节点指针
	OnEvicted func(key string, v Value) // 驱逐回调函数
}

// entry 键值对, 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// Value 数值, 使用 Len() 统计多少个 bytes
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		deQueue:   list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 根据 key 获取元素
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.deQueue.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除最旧的元素
func (c *Cache) RemoveOldest() {
	element := c.deQueue.Back()
	if element != nil {
		c.deQueue.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 调用驱逐回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 添加元素
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		// 命中缓存, 移动到队头
		c.deQueue.MoveToFront(element)
		kv := element.Value.(*entry)
		// 内存增加量: 新 value 长度 - 旧 value 长度
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		// 赋新值
		kv.value = value
	} else {
		// 未命中缓存, 插入到队头
		element := c.deQueue.PushFront(&entry{key, value})
		c.cache[key] = element
		// 内存增加量: 键 + 值总长度
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 超出内存阈值时, 持续移除最旧的元素直到占用内存量小于最大内存量
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len 返回 Cache 元素个数
func (c *Cache) Len() int {
	return c.deQueue.Len()
}
