package stream

import (
	"context"
	"sort"
	"sync"
)

type Stream[T any] struct {
	Data []*T
	mu   sync.Mutex // 添加互斥锁
}

// NewStream 方法，创建并返回 Stream 实例
func NewStream[T any](data []*T) *Stream[T] {
	return &Stream[T]{Data: data}
}

// Filter 方法，用于链式过滤
func (fs *Stream[T]) Filter(predicate func(*T) bool) *Stream[T] {
	fs.mu.Lock() // 加锁以保护数据
	defer fs.mu.Unlock()
	// 创建一个通道用于传输过滤后的数据
	out := make(chan *T)
	var wg sync.WaitGroup
	go func() {
		defer close(out)
		for _, v := range fs.Data {
			if predicate(v) {
				out <- v
			}
		}
	}()
	var result []*T
	wg.Add(1)
	go func() {
		defer wg.Done()
		for perm := range out {
			result = append(result, perm)
		}
	}()
	wg.Wait()
	fs.Data = result
	return fs
}

// ForEach 方法，遍历每个元素并执行指定操作
// ctx用于退出ctx, cancel := context.WithCancel(context.Background())
// defer cancel()
// if () { cancel()  }
// 满足条件时取消 context，从而中止遍历
func (fs *Stream[T]) ForEach(ctx context.Context, action func(*T)) *Stream[T] {
	fs.mu.Lock() // 加锁以保护数据
	defer fs.mu.Unlock()
	// 创建一个通道用于传输数据
	in := make(chan *T)
	go func() {
		for _, v := range fs.Data {
			in <- v
		}
		close(in) // 关闭输入通道
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range in {
			select {
			case <-ctx.Done(): // 如果 context 被取消，则提前退出
				return
			default:
				action(v)
			}
		}
	}()
	wg.Wait() // 等待所有 goroutine 完成
	return fs
}

// Distinct 方法，用于去除重复元素
func (fs *Stream[T]) Distinct(equal func(*T, *T) bool) *Stream[T] {
	fs.mu.Lock() // 加锁以保护数据
	defer fs.mu.Unlock()
	unique := make([]*T, 0)
	seen := make(map[*T]struct{})

	for _, v := range fs.Data {
		found := false
		for _, u := range unique {
			if equal(v, u) {
				found = true
				break
			}
		}
		if !found {
			unique = append(unique, v)
			seen[v] = struct{}{}
		}
	}

	fs.Data = unique
	return fs
}

// Map 方法，用于对每个元素应用指定的转换函数
func (fs *Stream[T]) Map(transform func(*T) *T) *Stream[T] {
	fs.mu.Lock() // 加锁以保护数据
	defer fs.mu.Unlock()

	// 创建一个通道用于传输映射后的数据
	out := make(chan struct {
		index int
		value *T
	}, len(fs.Data))
	var wg sync.WaitGroup

	// 启动 goroutine 进行并发映射
	wg.Add(len(fs.Data))
	for i, v := range fs.Data {
		go func(index int, val *T) {
			defer wg.Done()
			out <- struct {
				index int
				value *T
			}{index, transform(val)} // 应用转换函数并发送到通道
		}(i, v)
	}

	go func() {
		wg.Wait()  // 等待所有 goroutine 完成
		close(out) // 关闭输出通道
	}()

	// 创建一个用于保存结果的切片
	transformed := make([]*T, len(fs.Data))
	for result := range out {
		transformed[result.index] = result.value
	}

	return NewStream(transformed) // 返回新的 Stream 实例
}

// Sort 方法，用于对流中的元素进行排序
func (fs *Stream[T]) Sort(compare func(a, b *T) bool) *Stream[T] {
	fs.mu.Lock() // 加锁以保护数据
	defer fs.mu.Unlock()

	// 使用 sort.Slice 进行排序
	sort.Slice(fs.Data, func(i, j int) bool {
		return compare(fs.Data[i], fs.Data[j])
	})

	return fs // 返回当前 Stream 实例
}
