package cache

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Record is record. Nothing more, nothing less.
type Record[U UidType, D DataType] struct {
	uid            U
	data           D
	volume         int
	lastAccessTime uint
	cache          *Cache[U, D]
	upperRecord    *Record[U, D]
	lowerRecord    *Record[U, D]
}

// NewRecord creates a new cache record.
func NewRecord[U UidType, D DataType](cache *Cache[U, D], uid U, data D, useVolume bool) (rec *Record[U, D], err error) {
	//err = checkUid(uid)
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = checkData(data)
	//if err != nil {
	//	return nil, err
	//}

	rec = &Record[U, D]{
		uid:            uid,
		data:           data,
		lastAccessTime: 0, // See below.
		cache:          cache,
		upperRecord:    nil,
		lowerRecord:    nil,
	}

	if useVolume {
		rec.volume = getVariableSize(data)
	}

	rec.touch()

	return rec, nil
}

//func checkUid[U UidType](uid U) (err error) {
//	// TODO: Check for empty value when Go language gets a working type switch for generics.
//	switch v := uid.(type) {
//	case string:
//		if len(v) == 0 {
//			return errors.New(ErrUidIsEmpty)
//		}
//	case uint:
//		if v == 0 {
//			return errors.New(ErrUidIsEmpty)
//		}
//	case int:
//		if v == 0 {
//			return errors.New(ErrUidIsEmpty)
//		}
//	}
//
//	return nil
//}

//func checkData[D DataType](data D) (err error) {
//	if getVariableSize(data) == 0 {
//		return errors.New(ErrDataIsEmpty)
//	}
//
//	return nil
//}

func getVariableSize(x any) (size int) {
	// unsafe.SizeOf() says any string takes 16 bytes, but how? [duplicate]
	// https://stackoverflow.com/questions/65878177/unsafe-sizeof-says-any-string-takes-16-bytes-but-how
	//return int(unsafe.Sizeof(x))

	// A small "hack" of this stupid shitty piece of garbage named "Golang".
	// https://stackoverflow.com/questions/44257522/how-to-get-total-referenced-memory-of-a-variable/44258164#44258164
	// https://stackoverflow.com/a/60508928
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(x)
	if err != nil {
		panic(err)
	}
	return b.Len()
}

func (r *Record[U, D]) getVolume() int {
	return r.volume
}

func (r *Record[U, D]) setVolume(volume int) {
	r.volume = volume
}

func (r *Record[U, D]) moveToTop() {
	if r == r.cache.top {
		return
	}

	if r != r.cache.bottom {
		r.upperRecord.lowerRecord = r.lowerRecord
		r.lowerRecord.upperRecord = r.upperRecord
		r.lowerRecord = r.cache.top
		r.cache.top.upperRecord = r
		r.upperRecord = nil
		r.cache.top = r
	} else {
		r.cache.bottom = r.upperRecord
		r.cache.bottom.lowerRecord = nil
		r.upperRecord = nil
		r.lowerRecord = r.cache.top
		r.cache.top.upperRecord = r
		r.cache.top = r
	}

	// Size is not changed.
	// Volume is not changed.
	// Map is not changed.
}

func (r *Record[U, D]) touch() {
	r.lastAccessTime = uint(time.Now().Unix())
}

func (r *Record[U, D]) isAlive() bool {
	return uint(time.Now().Unix()) < r.lastAccessTime+r.cache.recordTtl
}

func (r *Record[U, D]) update(data D, useVolume bool) {
	oldVolume := r.getVolume()
	r.data = data

	if useVolume {
		r.setVolume(getVariableSize(data))
	}

	r.touch()

	r.cache.volume += r.getVolume() - oldVolume
	// Size is not changed.
	// Map is not changed.
}

func (r *Record[U, D]) unlink() {
	if r.cache.size == 1 {
		r.cache.top = nil
		r.cache.bottom = nil
	} else if r.cache.top == r {
		r.cache.top = r.lowerRecord
		r.cache.top.upperRecord = nil
		r.upperRecord = nil
		r.lowerRecord = nil
	} else if r.cache.bottom == r {
		r.cache.bottom = r.upperRecord
		r.cache.bottom.lowerRecord = nil
		r.upperRecord = nil
		r.lowerRecord = nil
	} else { // size is >= 3 and r is a middle record.
		r.upperRecord.lowerRecord = r.lowerRecord
		r.lowerRecord.upperRecord = r.upperRecord
		r.upperRecord = nil
		r.lowerRecord = nil
	}

	r.cache.size--
	r.cache.volume -= r.getVolume()
	delete(r.cache.recordsByUid, r.uid)
}
