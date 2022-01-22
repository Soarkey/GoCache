# GoCache
> 基于Go的分布式缓存数据库, 参考 https://github.com/geektutu/7days-golang

## 1. LRU 缓存淘汰策略
- 双端队列 + map
- 自定义传入缓存淘汰函数