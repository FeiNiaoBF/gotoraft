// Package store 提供存储接口和内存存储实现
package store

// Store 定义了存储接口
type Store interface {
	// 在这里添加你需要的存储方法
	// 例如：Get, Set, Delete 等
}

// MemoryStore 提供基于内存的存储实现
type MemoryStore struct {
	data map[string]interface{}
}

// NewMemoryStore 创建一个新的内存存储实例
func NewMemoryStore() Store {
	return &MemoryStore{
		data: make(map[string]interface{}),
	}
}
