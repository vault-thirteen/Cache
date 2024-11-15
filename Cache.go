package cache

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
	volume       int
	volumeLimit  int
	recordsByUid map[U]*Record[U, D]
	recordTtl    uint
	lock         *sync.RWMutex
}

// NewCache creates a new cache.
// If volumeLimit is > 0, volume calculation is enabled.
// Volume calculation consumes a lot of CPU time while Go language does not
// provide a method for fast calculation of real memory usage of a variable.
// The first attempt of Go language developers was to use reflection, but it
// was terribly slow. The second attempt is the 'unsafe.Sizeof' method, but
// this function does not show real size of a variable. Hand-made "hacks" work
// very slowly and, thus, their usage is optional. Unfortunately, in the year
// 2024 there is still no type switch in Go language for generic variables
// which could make this problem a little bit less annoying.
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

func (c *Cache[U, D]) initialize(sizeLimit int, volumeLimit int, recordTtl uint) {
	c.top = nil
	c.bottom = nil
	c.size = 0
	c.sizeLimit = sizeLimit
	c.volumeLimit = volumeLimit
	c.recordsByUid = make(map[U]*Record[U, D])
	c.recordTtl = recordTtl
	c.lock = new(sync.RWMutex)

	if c.hasLimitedVolume() {
		c.volume = 0
	}
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

	if c.hasLimitedVolume() {
		c.volume += rec.getVolume()
	}

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

	if c.hasLimitedVolume() {
		c.volume -= rec.getVolume()
	}

	delete(c.recordsByUid, rec.uid)

	return rec, nil
}

func (c *Cache[U, D]) getFreeVolume() int {
	return c.volumeLimit - c.volume
}

// GetVolume returns current volume of the cache.
func (c *Cache[U, D]) GetVolume() (usedVolume int, volumeLimit int) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.volume, c.volumeLimit
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
		rec.update(data, true)
	} else {
		// UID is not found,
		// we add a new record.
		rec, err = NewRecord(c, uid, data, c.hasLimitedVolume())
		if err != nil {
			return err
		}

		if c.hasLimitedVolume() && (rec.getVolume() > c.volumeLimit) {
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

// RemoveRecord removes a record from the cache.
func (c *Cache[U, D]) RemoveRecord(uid U) (recExists bool, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var rec *Record[U, D]
	rec, recExists = c.recordsByUid[uid]
	if !recExists {
		return recExists, fmt.Errorf(ErrRecordIsNotFound, uid)
	}

	rec.unlink()

	return recExists, nil
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
