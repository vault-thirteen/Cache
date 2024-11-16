package nvl

import (
	"errors"
	"fmt"
	"sync"
)

// Cache is cache. Surprisingly, but it is true.
type Cache[U UidType, D DataType] struct {
	top          *Record[U, D]
	bottom       *Record[U, D]
	size         int
	sizeLimit    int
	recordsByUid map[U]*Record[U, D]
	recordTtl    uint
	lock         *sync.RWMutex
}

// NewCache creates a new cache.
func NewCache[U UidType, D DataType](sizeLimit int, recordTtl uint) (cache *Cache[U, D]) {
	if recordTtl == 0 {
		panic(ErrTtlIsZero)
	}

	cache = new(Cache[U, D])
	cache.initialize(sizeLimit, recordTtl)
	return cache
}

func (c *Cache[U, D]) initialize(sizeLimit int, recordTtl uint) {
	c.top = nil
	c.bottom = nil
	c.size = 0
	c.sizeLimit = sizeLimit
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
	delete(c.recordsByUid, rec.uid)

	return rec, nil
}

// RecordExists checks whether the specified record exists or not. If the
// record is outdated, it is removed from the cache.
func (c *Cache[U, D]) RecordExists(uid U) (recordExists bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var rec *Record[U, D]
	rec, recordExists = c.recordsByUid[uid]
	if !recordExists {
		return false
	}

	if !rec.isAlive() {
		rec.unlink()
		return false
	}

	return true
}

// AddRecord either adds a new record to the top of the cache or moves an
// existing record to the top of the cache. If the record already exists, its
// data and LAT are updated.
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

	return nil
}

// GetRecord reads a record from the cache. If the record is outdated, it is
// removed from the cache and is not returned.
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

// RemoveRecord safely removes a record from the cache.
func (c *Cache[U, D]) RemoveRecord(uid U) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var rec *Record[U, D]
	var recExists bool
	rec, recExists = c.recordsByUid[uid]
	if !recExists {
		return
	}

	rec.unlink()

	return
}

// RemoveExistingRecord removes an existing record from the cache.
func (c *Cache[U, D]) RemoveExistingRecord(uid U) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var rec *Record[U, D]
	var recExists bool
	rec, recExists = c.recordsByUid[uid]
	if !recExists {
		return fmt.Errorf(ErrRecordIsNotFound, uid)
	}

	rec.unlink()

	return nil
}

// Clear removes all records from the cache.
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
