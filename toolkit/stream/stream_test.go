package stream

import (
	"context"
	"fmt"
	"testing"
)

//
//func TestStream(t *testing.T) {
//	// 创建一个 Stream 实例
//	data := []*int{new(int), new(int), new(int), new(int), new(int), new(int)}
//	*data[0] = 1
//	*data[1] = 2
//	*data[2] = 2
//	*data[3] = 3
//	*data[4] = 3
//	*data[5] = 3
//
//	NewStream(data).Filter(func(i *int) bool {
//		return *i >= 2
//	}).ForEach(context.Background(), func(i *int) {
//		fmt.Println(*i)
//	}).
//		Map(func(i *int) *int {
//			val := *i * 2
//			return &val
//		}).ForEach(context.Background(), func(i *int) {
//		fmt.Println("map", *i)
//	}).
//		Distinct(func(i *int, i2 *int) bool {
//			return *i == *i2
//		}).ForEach(context.Background(), func(i *int) {
//		fmt.Println("distinct", *i)
//	})
//}

func TestStream(t *testing.T) {
	data := []*int{new(int), new(int), new(int), new(int), new(int)}
	*data[0] = 1
	*data[1] = 2
	*data[2] = 2
	*data[3] = 3
	*data[4] = 3

	stream := NewStream(data)

	stream.Filter(func(i *int) bool {
		return *i >= 2
	}).Sort(func(a, b *int) bool {
		return *a < *b
	}).ForEach(context.Background(), func(i *int) {
		fmt.Println(*i)
	}).Map(func(i *int) *int {
		val := *i * 2
		newVal := val
		return &newVal // Return the new pointer
	}).ForEach(context.Background(), func(i *int) {
		fmt.Println("map: ", *i)
	}).Distinct(func(i *int, i2 *int) bool {
		return *i == *i2
	}).ForEach(context.Background(), func(i *int) {
		fmt.Println("distinct", *i)
	})
}
