package cache

import (
	"errors"
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
func NewRecord[U UidType, D DataType](cache *Cache[U, D], uid U, data D) (rec *Record[U, D], err error) {
	err = checkUid(uid)
	if err != nil {
		return nil, err
	}

	err = checkData(data)
	if err != nil {
		return nil, err
	}

	rec = &Record[U, D]{
		uid:            uid,
		data:           data,
		volume:         len(data),
		lastAccessTime: 0, // See below.
		cache:          cache,
		upperRecord:    nil,
		lowerRecord:    nil,
	}

	rec.touch()

	return rec, nil
}

func checkUid[U UidType](uid U) (err error) {
	//TODO: Check for empty value when Go language gets a working type switch
	// for generics.
	//switch v := uid.(type) {
	//case string:
	//	if len(v) == 0 {
	//		return errors.New(ErrUidIsEmpty)
	//	}
	//case uint:
	//	if v == 0 {
	//		return errors.New(ErrUidIsEmpty)
	//	}
	//case int:
	//	if v == 0 {
	//		return errors.New(ErrUidIsEmpty)
	//	}
	//}

	return nil
}

func checkData[D DataType](data D) (err error) {
	if len(data) == 0 {
		return errors.New(ErrDataIsEmpty)
	}

	return nil
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

func (r *Record[U, D]) update(data D) {
	oldVolume := r.volume
	r.data = data
	r.volume = len(data)

	r.touch()

	r.cache.volume += r.volume - oldVolume
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
	r.cache.volume -= r.volume
	delete(r.cache.recordsByUid, r.uid)
}
