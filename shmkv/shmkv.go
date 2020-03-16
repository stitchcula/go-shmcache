package shmkv

/*
#cgo LDFLAGS: -lshmcache
#include "shmcache/shmcache.h"
*/
import "C"
import (
	"errors"
	"reflect"
	_ "runtime/cgo"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

type ShmKv struct {
	handle *C.struct_shmcache_context
}

func NewShmKvFromFile(ConfigPath string) (*ShmKv, error) {
	CConfigPath := C.CString(ConfigPath)
	//defer C.free(unsafe.Pointer(CConfigPath))

	handle := C.struct_shmcache_context{}
	eno := C.shmcache_init_from_file(&handle, CConfigPath)
	return &ShmKv{handle: &handle}, errors.New(strconv.Itoa(int(eno)))
}

func (c *ShmKv) Set(k string, v []byte, d time.Duration, flag int) error {
	vh := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return syscall.Errno(C.shmcache_set_ex(c.handle, (*C.struct_shmcache_key_info)(unsafe.Pointer(&k)), &C.struct_shmcache_value_info{
		data:    (*C.char)(unsafe.Pointer(vh.Data)),
		length:  C.int(vh.Len),
		expires: C.time_t(time.Now().Add(d).Unix()),
		options: C.int(flag),
	}))
}

func (c *ShmKv) Delete(k string) {
	_ = syscall.Errno(C.shmcache_delete(c.handle, (*C.struct_shmcache_key_info)(unsafe.Pointer(&k))))
	// ENOENT
}

func (c *ShmKv) GetWithExpiration(k string) ([]byte, time.Time, int, bool) {
	vi := &C.struct_shmcache_value_info{}
	eno := syscall.Errno(C.shmcache_get(c.handle, (*C.struct_shmcache_key_info)(unsafe.Pointer(&k)), vi))
	if eno != 0 { // ENOENT or ETIMEDOUT
		return nil, time.Time{}, 0, false
	}
	return C.GoBytes(unsafe.Pointer(vi.data), C.int(vi.length)), time.Unix(int64(vi.expires), 0), int(vi.options), true
}

func (c *ShmKv) Flush() {
	_ = syscall.Errno(C.shmcache_remove_all(c.handle))
}
