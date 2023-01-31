package Cache

import (
	"sync"
	"testing"
	"time"

	"github.com/vault-thirteen/tester"
)

func Test_NewCache(t *testing.T) {
	aTest := tester.New(t)

	// Test.
	var c *Cache[string, string] = nil
	aTest.MustBeEqual(c, (*Cache[string, string])(nil))

	c = NewCache[string, string](0, 0, 60)
	aTest.MustBeDifferent(c, (*Cache[string, string])(nil))

	aTest.MustBeEqual(c.recordTtl, uint(60))
}

func Test_initialize(t *testing.T) {
	aTest := tester.New(t)

	// Test.
	var c = new(Cache[string, string])
	aTest.MustBeEqual(c.lock, (*sync.RWMutex)(nil))
	c.initialize(1, 2, 3)

	aTest.MustBeEqual(c.top, (*Record[string, string])(nil))
	aTest.MustBeEqual(c.bottom, (*Record[string, string])(nil))
	aTest.MustBeEqual(c.size, 0)
	aTest.MustBeEqual(c.sizeLimit, 1)
	aTest.MustBeEqual(c.volume, 0)
	aTest.MustBeEqual(c.volumeLimit, 2)
	aTest.MustBeEqual(c.recordTtl, uint(3))
	aTest.MustBeDifferent(c.lock, (*sync.RWMutex)(nil))
}

func Test_hasLimitedSize(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]

	// Test #1.
	c = NewCache[string, string](0, 0, 60)
	aTest.MustBeEqual(c.hasLimitedSize(), false)

	// Test #2.
	c = NewCache[string, string](1, 0, 60)
	aTest.MustBeEqual(c.hasLimitedSize(), true)
}

func Test_isEmpty(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]

	// Test #1.
	c = NewCache[string, string](0, 0, 60)
	aTest.MustBeEqual(c.isEmpty(), true)

	// Test #2.
	c.size++
	aTest.MustBeEqual(c.isEmpty(), false)
}

func Test_isNotEmpty(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]

	// Test #1.
	c = NewCache[string, string](0, 0, 60)
	aTest.MustBeEqual(c.isNotEmpty(), false)

	// Test #2.
	c.size++
	aTest.MustBeEqual(c.isNotEmpty(), true)
}

func Test_hasLimitedVolume(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]

	// Test #1.
	c = NewCache[string, string](0, 0, 60)
	aTest.MustBeEqual(c.hasLimitedVolume(), false)

	// Test #2.
	c = NewCache[string, string](0, 1, 60)
	aTest.MustBeEqual(c.hasLimitedVolume(), true)
}

func Test_linkNewTopRecord(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var r *Record[string, string]
	var ok bool

	// Test. 0R -> 1R.
	c = _test_prepare_0_cache()
	r = _test_prepare_record_Q(aTest, c)
	c.linkNewTopRecord(r) // {} -> Q.
	ok = _test_ensure_order_1_record(c, "Q", "W")
	aTest.MustBeEqual(ok, true)

	// Test. 1R -> 2R.
	c = _test_prepare_A_cache(aTest)
	r = _test_prepare_record_Q(aTest, c)
	c.linkNewTopRecord(r) // A -> QA.
	ok = _test_ensure_order_2_records(c, [2]string{"Q", "A"}, [2]string{"W", "1"})
	aTest.MustBeEqual(ok, true)

	// Test. 2R -> 3R.
	c = _test_prepare_AB_cache(aTest)
	r = _test_prepare_record_Q(aTest, c)
	c.linkNewTopRecord(r) // AB -> QAB.
	ok = _test_ensure_order_3_records(c, [3]string{"Q", "A", "B"}, [3]string{"W", "1", "2"})
	aTest.MustBeEqual(ok, true)
}

func Test_unlinkBottomRecord(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var r *Record[string, string]
	var ok bool
	var err error

	// Test. 0R -> 0R.
	c = _test_prepare_0_cache()
	r, err = c.unlinkBottomRecord() // {} -> {}.
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrBottomRecordDoesNotExist)
	aTest.MustBeEqual(r, (*Record[string, string])(nil))

	// Test. 1R -> 0R.
	c = _test_prepare_A_cache(aTest)
	r, err = c.unlinkBottomRecord() // A -> {}.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_0_records(c)
	aTest.MustBeEqual(ok, true)
	aTest.MustBeEqual(r.uid, "A")

	// Test. 2R -> 1R.
	c = _test_prepare_AB_cache(aTest)
	r, err = c.unlinkBottomRecord() // AB -> A.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_1_record(c, "A", "1")
	aTest.MustBeEqual(ok, true)
	aTest.MustBeEqual(r.uid, "B")

	// Test. 3R -> 2R.
	c = _test_prepare_ABC_cache(aTest)
	r, err = c.unlinkBottomRecord() // ABC -> AB.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_2_records(c, [2]string{"A", "B"}, [2]string{"1", "2"})
	aTest.MustBeEqual(ok, true)
	aTest.MustBeEqual(r.uid, "C")
}

func Test_getFreeVolume(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]

	// Test #1.
	c = NewCache[string, string](0, 10, 60)
	aTest.MustBeEqual(c.getFreeVolume(), 10)

	// Test #2.
	c.volume++
	aTest.MustBeEqual(c.getFreeVolume(), 9)
}

func Test_GetVolume(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]

	// Test.
	c = NewCache[string, string](0, 10, 60)
	c.volume = 3
	var usedVolume int
	var volumeLimit int
	usedVolume, volumeLimit = c.GetVolume()
	aTest.MustBeEqual(usedVolume, 3)
	aTest.MustBeEqual(volumeLimit, 10)
}

func Test_AddRecord(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var ok bool
	var err error
	var oldLatOfRecordB, newLatOfRecordB uint

	// Preparation for Test #1.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)
	oldLatOfRecordB = c.recordsByUid["B"].lastAccessTime

	// Test #1. Record already exists.
	time.Sleep(time.Second * 1)    // LAT++
	err = c.AddRecord("B", "test") // AB -> BA.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_2_records(c, [2]string{"B", "A"}, [2]string{"test", "1"})
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 1+4)
	newLatOfRecordB = c.recordsByUid["B"].lastAccessTime
	aTest.MustBeEqual(newLatOfRecordB-oldLatOfRecordB > 0, true)

	// Preparation for Test #2.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)

	// Test #2. Record is new. Records is bad.
	err = c.AddRecord("Q", "") // AB -> AB.
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)
	ok = _test_ensure_order_2_records(c, [2]string{"A", "B"}, [2]string{"1", "2"})
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)

	// Preparation for Test #3.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)
	c.volumeLimit = 3

	// Test #3. Record is new. Record is too big.
	err = c.AddRecord("Q", "test") // AB -> AB.
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrRecordIsTooBig)
	ok = _test_ensure_order_2_records(c, [2]string{"A", "B"}, [2]string{"1", "2"})
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)

	// Preparation for Test #4.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)

	// Test #4. Record is new. No limits.
	err = c.AddRecord("Q", "test") // AB -> QAB.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_3_records(c, [3]string{"Q", "A", "B"}, [3]string{"test", "1", "2"})
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 3)
	aTest.MustBeEqual(c.volume, 4+1+1)

	// Preparation for Test #5.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)
	c.sizeLimit = 2

	// Test #5. Record is new. Size is limited.
	err = c.AddRecord("Q", "test") // AB -> QA.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_2_records(c, [2]string{"Q", "A"}, [2]string{"test", "1"})
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 4+1)

	// Preparation for Test #6.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)
	c.volumeLimit = 4

	// Test #6. Record is new. Volume is limited, 1 record is deleted.
	err = c.AddRecord("Q", "xxx") // AB -> QAB -> QA.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_2_records(c, [2]string{"Q", "A"}, [2]string{"xxx", "1"})
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 3+1)

	// Preparation for Test #7.
	c = _test_prepare_AB_cache(aTest) // AB.
	aTest.MustBeEqual(c.size, 2)
	aTest.MustBeEqual(c.volume, 2)
	c.volumeLimit = 3

	// Test #7. Record is new. Volume is limited, 2 records are deleted.
	err = c.AddRecord("Q", "xxx") // AB -> QAB -> Q.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_1_record(c, "Q", "xxx")
	aTest.MustBeEqual(ok, true)
	// Also check new values of size, volume and record's TTL.
	aTest.MustBeEqual(c.size, 1)
	aTest.MustBeEqual(c.volume, 3)
}

func Test_GetRecord(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var ok bool
	var data string
	var err error

	c = _test_prepare_ABC_cache_with_low_ttl(aTest) // ABC.

	// Test #1. 3R, Record is not found.
	data, err = c.GetRecord("Junk") // ABC -> ABC.
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), `record is not found, uid=Junk`)
	ok = _test_ensure_order_3_records(c, [3]string{"A", "B", "C"}, [3]string{"1", "2", "3"})
	aTest.MustBeEqual(ok, true)

	// Test #2. 3R, Record is outdated.
	// Wait for the record to become outdated. N.B.: TTL is 3 Seconds.
	time.Sleep(time.Second * (3 + 1))
	data, err = c.GetRecord("B") // ABC -> AC.
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), `record is outdated, uid=B`)
	ok = _test_ensure_order_2_records(c, [2]string{"A", "C"}, [2]string{"1", "3"})
	aTest.MustBeEqual(ok, true)

	c = _test_prepare_ABC_cache_with_low_ttl(aTest) // ABC.
	oldLatOfRecordB := c.recordsByUid["B"].lastAccessTime

	// Test #3. 3R, Record is alive.
	// Wait a bit, but not more than TTL period. N.B.: TTL is 3 Seconds.
	time.Sleep(time.Second * 1)  // 1 Sec.
	data, err = c.GetRecord("B") // ABC -> BAC.
	newLatOfRecordB := c.recordsByUid["B"].lastAccessTime
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(data, "2")
	// Ensure that the requested record has been moved to top and its LAT has
	// been updated.
	ok = _test_ensure_order_3_records(c, [3]string{"B", "A", "C"}, [3]string{"2", "1", "3"})
	aTest.MustBeEqual(ok, true)
	aTest.MustBeEqual(newLatOfRecordB-oldLatOfRecordB > 0, true)
}

func Test_Clear(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var ok bool
	var err error

	// Test. 0R -> 0R.
	c = _test_prepare_0_cache()
	err = c.Clear() // {} -> {}.
	aTest.MustBeNoError(err)

	// Test. 1R -> 0R.
	c = _test_prepare_A_cache(aTest)
	err = c.Clear() // A -> {}.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_0_records(c)
	aTest.MustBeEqual(ok, true)

	// Test. 2R -> 0R.
	c = _test_prepare_AB_cache(aTest)
	err = c.Clear() // AB -> {}.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_0_records(c)
	aTest.MustBeEqual(ok, true)

	// Test. 3R -> 0R.
	c = _test_prepare_ABC_cache(aTest)
	err = c.Clear() // ABC -> {}.
	aTest.MustBeNoError(err)
	ok = _test_ensure_order_0_records(c)
	aTest.MustBeEqual(ok, true)
}
