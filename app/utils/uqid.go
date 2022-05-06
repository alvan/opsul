package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"os"
	"sync/atomic"
	"time"
)

var (
	Uqid = new(uqid)
)

type uqid struct {
	counterId uint32
	machineId []byte
}

func (this *uqid) Init() {
	this.machineId = make([]byte, 3)

	name, err := os.Hostname()
	if err == nil && len(name) != 0 {
		hash := md5.New()
		hash.Write([]byte(name))
		copy(this.machineId, hash.Sum(nil))
	} else {
		rand.Reader.Read(this.machineId)
	}
}

func (this *uqid) Next() string {
	var b [12]byte

	// Timestamps, 4 bytes
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine ID, 3 bytes
	b[4] = this.machineId[0]
	b[5] = this.machineId[1]
	b[6] = this.machineId[2]
	// Process ID, 2 bytes
	p := os.Getpid()
	b[7] = byte(p >> 8)
	b[8] = byte(p)
	// Counter ID, 3 bytes
	i := atomic.AddUint32(&this.counterId, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	return hex.EncodeToString(b[:])
}

func (this *uqid) Less() string {
	var b [6]byte

	// Timestamps, 4 bytes
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Counter ID, 2 bytes
	i := atomic.AddUint32(&this.counterId, 1)
	b[4] = byte(i >> 8)
	b[5] = byte(i)

	return hex.EncodeToString(b[:])
}

func (this *uqid) Time(id string) (ret time.Time, err error) {
	b, err := hex.DecodeString(id[0:8])
	if err == nil {
		ret = time.Unix(int64(binary.BigEndian.Uint32(b)), 0)
	}

	return
}

func init() {
	Uqid.Init()
}
