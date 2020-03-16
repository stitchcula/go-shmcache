package shmcache

import "time"

type ShareType int

const (
	MMap ShareType = 2
	Shm  ShareType = 1
)

func NewDefaultOptions() *Options {
	return &Options{
		ShareType:           MMap,
		MaxMemory:           256 * 1024 * 1024,
		MinMemory:           0,
		SegmentSize:         8 * 1024 * 1024,
		MaxKeyCount:         200000,
		MaxValueSize:        200 * 1024,
		HashFunc:            "simple_hash",
		RecycleKeyOnce:      0,
		RecycleValidEntries: true,
		AllocatorPolicy: AllocatorPolicy{
			AvgKeyTTL:                    0,
			DiscardMemorySize:            128,
			MaxFailTimes:                 5,
			SleepWhenRecycleValidEntries: 1000 * time.Microsecond,
		},
		LockPolicy: LockPolicy{
			TryLockInterval:        200 * time.Microsecond,
			DetectDeadlockInterval: 1000 * time.Millisecond,
		},
		SyslogLevel: "error",
	}
}

// Options contains configurable options.
type Options struct {
	// shared memory type, shm or mmap
	// shm for SystemV pure shared memory
	// mmap for POSIX shared memory based file
	// default value is mmap
	// Note: when type is shm, the shm limit is too small in FreeBSD and MacOS,
	//  you should increase following kernel parameters:
	//    kern.sysv.shmmax
	//    kern.sysv.shmall
	//    kern.sysv.shmseg
	//  you can refer to http://www.unidata.ucar.edu/software/mcidas/2008/users_guide/workstation.html
	ShareType ShareType

	// the memory limit
	// the oldest memory will be recycled when this max memory reached
	MaxMemory int64

	// the min memory
	// default: 0, means do NOT set min memory
	MinMemory int64

	// the memory segment size for incremental memory allocation
	SegmentSize int64

	// the key number limit
	MaxKeyCount int64

	// the size limit for one value
	MaxValueSize int64

	// the hash function in libfastcommon/src/hash.h
	HashFunc string

	// recycle key number once when reach max keys
	// <= 0 means recycle one memory striping
	RecycleKeyOnce int

	// if recycle valid entries by FIFO
	RecycleValidEntries bool

	AllocatorPolicy AllocatorPolicy
	LockPolicy      LockPolicy

	// standard log level as syslog, case insensitive, value list:
	//  emerg for emergency
	//  alert
	//  crit for critical
	//  error
	//  warn for warning
	//  notice
	//  info
	//  debug
	SyslogLevel string
}

// AllocatorPolicy contains value allocator policy
type AllocatorPolicy struct {
	// avg. key TTL threshold for recycling memory
	// <= 0 for never recycle memory until reach memory limit (max_memory)
	AvgKeyTTL time.Duration

	// when the remain memory <= this parameter, discard it
	DiscardMemorySize int

	// when a value allocator allocate fail times > this parameter,
	// means it is almost full, discard it
	MaxFailTimes int

	// sleep time to avoid other processes read dirty data when
	// recycle more than one valid (in TTL / not expired) KV entries
	// 0 for never sleep
	SleepWhenRecycleValidEntries time.Duration
}

type LockPolicy struct {
	// try lock interval in us, must great than zero
	TryLockInterval time.Duration

	// the interval to detect deadlock caused by the crushed process
	// must great than zero
	DetectDeadlockInterval time.Duration
}
