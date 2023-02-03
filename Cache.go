package cache

import (
	"errors"
	"fmt"
	"sync"
)

type Cache[U UidType, D DataType] struct {
	top          *Record[U, D]
	bottom       *Record[U, D]
	size         int
	sizeLimit    int
	volume       int
	volumeLimit  int
	recordsByUid map[U]*Record[U, D]
	recordTtl    uint
	lock         *sync.RWMutex
}

func NewCache[U UidType, D DataType](
	sizeLimit int,
	volumeLimit int,
	recordTtl uint,
) (cache *Cache[U, D]) {
	if recordTtl == 0 {
		panic(ErrTtlIsZero)
	}

	cache = new(Cache[U, D])
	cache.initialize(sizeLimit, volumeLimit, recordTtl)
	return cache
}

func (c *Cache[U, D]) initialize(
	sizeLimit int,
	volumeLimit int,
	recordTtl uint,
) {
	c.top = nil
	c.bottom = nil
	c.size = 0
	c.sizeLimit = sizeLimit
	c.volume = 0
	c.volumeLimit = volumeLimit
	c.recordsByUid = make(map[U]*Record[U, D])
	c.recordTtl = recordTtl
	c.lock = new(sync.RWMutex)
}

func (c *Cache[U, D]) hasLimitedSize() bool {
	return c.sizeLimit > 0
}

func (c *Cache[U, D]) isEmpty() bool {
	return c.size == 0
}

func (c *Cache[U, D]) isNotEmpty() bool {
	return c.size > 0
}

func (c *Cache[U, D]) hasLimitedVolume() bool {
	return c.volumeLimit > 0
}

func (c *Cache[U, D]) linkNewTopRecord(rec *Record[U, D]) {
	if c.isNotEmpty() {
		rec.lowerRecord = c.top
		c.top.upperRecord = rec
		c.top = rec
	} else {
		c.top = rec
		c.bottom = rec
	}

	c.size++
	c.volume += rec.volume
	c.recordsByUid[rec.uid] = rec
}

func (c *Cache[U, D]) unlinkBottomRecord() (rec *Record[U, D], err error) {
	rec = c.bottom

	if c.size > 1 {
		c.bottom = rec.upperRecord
		c.bottom.lowerRecord = nil
	} else if c.size == 1 {
		c.top = nil
		c.bottom = nil
	} else { // size = 0.
		return nil, errors.New(ErrBottomRecordDoesNotExist)
	}

	rec.upperRecord = nil
	rec.lowerRecord = nil

	c.size--
	c.volume -= rec.volume
	delete(c.recordsByUid, rec.uid)

	return rec, nil
}

func (c *Cache[U, D]) getFreeVolume() int {
	return c.volumeLimit - c.volume
}

func (c *Cache[U, D]) GetVolume() (usedVolume int, volumeLimit int) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.volume, c.volumeLimit
}

func (c *Cache[U, D]) AddRecord(uid U, data D) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var rec *Record[U, D]
	var recExists bool
	rec, recExists = c.recordsByUid[uid]
	if recExists {
		// If the UID is already used,
		// we update data of the record having this UID.
		rec.moveToTop()
		rec.update(data)
	} else {
		// UID is not found,
		// we add a new record.
		rec, err = NewRecord(c, uid, data)
		if err != nil {
			return err
		}

		if c.hasLimitedVolume() && (rec.volume > c.volumeLimit) {
			return errors.New(ErrRecordIsTooBig)
		}

		c.linkNewTopRecord(rec)
	}

	// Now we need (or do not need) to apply various constraints.
	if c.hasLimitedSize() {
		if c.size > c.sizeLimit {
			n := c.size - c.sizeLimit
			for i := 1; i <= n; i++ {
				_, err = c.unlinkBottomRecord()
				if err != nil {
					return err
				}
			}
		}
	}

	if c.hasLimitedVolume() {
		for {
			if c.getFreeVolume() >= 0 {
				return nil
			}

			_, err = c.unlinkBottomRecord()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Cache[U, D]) GetRecord(uid U) (data D, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var rec *Record[U, D]
	var ok bool
	rec, ok = c.recordsByUid[uid]
	if !ok {
		return data, fmt.Errorf(ErrRecordIsNotFound, uid)
	}

	if !rec.isAlive() {
		rec.unlink()
		return data, fmt.Errorf(ErrRecordIsOutdated, uid)
	}

	rec.moveToTop()
	rec.touch()

	return rec.data, nil
}

func (c *Cache[U, D]) Clear() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for {
		if c.isEmpty() {
			break
		}

		_, err = c.unlinkBottomRecord()
		if err != nil {
			return err
		}
	}

	return nil
}
