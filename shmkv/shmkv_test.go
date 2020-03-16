package shmkv

import (
	"bytes"
	"testing"
	"time"
)

func TestShmKv(t *testing.T) {
	kv, err := NewShmKvFromFile("/etc/libshmcache.conf")
	if err != nil {
		t.Fatal(err)
	}
	v, tm, f, ok := kv.GetWithExpiration("hello")
	if ok {
		t.Fatal("found hello ok", v, tm, f)
	}
	kv.Set("hello", []byte("world"), 2*time.Second, 1)
	v, tm, f, ok = kv.GetWithExpiration("hello")
	if !ok {
		t.Fatal("found hello !ok")
	} else if bytes.Equal(v, []byte("world")) {
		t.Fatal(string(v))
	} else if f != 1 {
		t.Fatal(f)
	}
	time.Sleep(3 * time.Second)
	v, tm, f, ok = kv.GetWithExpiration("hello")
	if ok {
		t.Fatal("found hello ok", v, tm, f)
	}
	kv.Set("hello", []byte("world"), 2*time.Second, 1)
	kv.Delete("hello")
	v, tm, f, ok = kv.GetWithExpiration("hello")
	if ok {
		t.Fatal("found hello ok", v, tm, f)
	}
}
