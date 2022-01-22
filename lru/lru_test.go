package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("hit failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("miss key")
	}
}

func TestCache_Add(t *testing.T) {
	lru := New(int64(8), nil)
	lru.Add("k1", String("v1"))
	lru.Add("k2", String("v2"))
	lru.Add("k2", String("v2"))
	keys := make([]string, lru.Len())
	idx := 0
	for i := lru.deQueue.Front(); i != nil; i = i.Next() {
		keys[idx] = i.Value.(*entry).key
		idx++
	}
	// fmt.Println(len(keys), ":", keys)
	expect := []string{"k2", "k1"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("add failed")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	n := len(k1 + k2 + v1 + v2)
	lru := New(int64(n), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	// for k, v := range lru.cache {
	// 	fmt.Printf("%v: %v\n", k, v.Value.(*entry))
	// }
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("lru algo failed")
	}
}

func TestCache_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(k string, v Value) {
		keys = append(keys, k)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))
	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("call OnEvicted failed")
	}
}
