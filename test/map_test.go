package test

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMapMem(t *testing.T) {
	v := struct{}{}

	a := make(map[int]struct{})

	for i := 0; i < 10000; i++ {
		a[i] = v
	}
	runtime.GC()
	memData("添加1万个键值对后")
	fmt.Println("删除前Map长度：", len(a))

	for i := 0; i < 10000-1; i++ {
		delete(a, i)
	}
	fmt.Println("删除后Map长度：", len(a))

	// 再次进行手动GC回收
	runtime.GC()
	memData("删除1万个键值对后")

	for i := 0; i < 10000-1; i++ {
		a[i] = v
	}

	// 再次进行手动GC回收
	runtime.GC()
	memData("再一次添加1万个键值对后")

	// 设置为nil进行回收
	a = nil
	runtime.GC()
	memData("设置为nil后")
}

func memData(mag string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%v：分配的内存 = %vKB, GC的次数 = %v\n", mag, m.Alloc/1024, m.NumGC)
}
