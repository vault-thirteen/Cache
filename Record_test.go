package cache

import (
	"testing"
	"time"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewRecord(t *testing.T) {
	aTest := tester.New(t)
	var err error
	var r *Record[string, string]

	// Test #1. checkUid fails.
	// TODO: Wait for Go language update for generics.
	//r, err = NewRecord[string, string](nil, "", "data")
	//aTest.MustBeAnError(err)
	//aTest.MustBeEqual(err.Error(), ErrUidIsEmpty)

	// Test #2. checkData fails.
	r, err = NewRecord[string, string](nil, "uid", "")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)

	// Test #3. OK.
	r, err = NewRecord[string, string](nil, "uid", "data")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(r.uid, "uid")
	aTest.MustBeEqual(r.data, "data")
	aTest.MustBeEqual(r.volume, 4)
	aTest.MustBeEqual(r.cache, (*Cache[string, string])(nil))
	aTest.MustBeEqual(r.upperRecord, (*Record[string, string])(nil))
	aTest.MustBeEqual(r.lowerRecord, (*Record[string, string])(nil))
}

func Test_checkUid(t *testing.T) {
	// TODO: Wait for Go language update for generics.
}

func Test_checkData(t *testing.T) {
	aTest := tester.New(t)
	var err error

	// Test #1. Bad data.
	err = checkData[string]("")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)

	// Test #2. OK.
	err = checkData[string]("ok")
	aTest.MustBeNoError(err)
}

func Test_moveToTop(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var ok bool

	c = _test_prepare_ABC_cache(aTest)

	// Test #1. 3R, Top.
	c.top.moveToTop() // ABC -> ABC.
	ok = _test_ensure_order_3_records(c, [3]string{"A", "B", "C"}, [3]string{"1", "2", "3"})
	aTest.MustBeEqual(ok, true)

	// Test #2. 3R, Bottom.
	c.bottom.moveToTop() // ABC -> CAB.
	ok = _test_ensure_order_3_records(c, [3]string{"C", "A", "B"}, [3]string{"3", "1", "2"})
	aTest.MustBeEqual(ok, true)

	// Test #3. 3R, Middle.
	c.top.lowerRecord.moveToTop() // CAB -> ACB.
	ok = _test_ensure_order_3_records(c, [3]string{"A", "C", "B"}, [3]string{"1", "3", "2"})
	aTest.MustBeEqual(ok, true)

	c = _test_prepare_AB_cache(aTest)

	// Test #4. 2R, Top.
	c.top.moveToTop() // AB -> AB.
	ok = _test_ensure_order_2_records(c, [2]string{"A", "B"}, [2]string{"1", "2"})
	aTest.MustBeEqual(ok, true)

	// Test #5. 2R, Bottom.
	c.bottom.moveToTop() // AB -> BA.
	ok = _test_ensure_order_2_records(c, [2]string{"B", "A"}, [2]string{"2", "1"})
	aTest.MustBeEqual(ok, true)

	c = _test_prepare_A_cache(aTest)

	// Test #6. 1R, Top.
	c.top.moveToTop() // A -> A.
	ok = _test_ensure_order_1_record(c, "A", "1")
	aTest.MustBeEqual(ok, true)
}

func Test_isAlive(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var err error
	c = NewCache[string, string](0, 0, 1)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)

	// Test #1. Fresh record.
	aTest.MustBeEqual(c.top.isAlive(), true)

	// Test #2. Stale record.
	time.Sleep(time.Second * time.Duration(2))
	aTest.MustBeEqual(c.top.isAlive(), false)
}

func Test_update(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var err error
	c = NewCache[string, string](0, 0, 60)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(c.volume, 1)

	// Test #1.
	c.top.update("333")
	aTest.MustBeEqual(c.top.data, "333")
	aTest.MustBeEqual(c.volume, 3)
}

func Test_unlink(t *testing.T) {
	aTest := tester.New(t)
	var c *Cache[string, string]
	var ok bool

	// Test. 1R -> 0R.
	c = _test_prepare_A_cache(aTest)
	c.top.unlink() // A -> {}.
	ok = _test_ensure_order_0_records(c)
	aTest.MustBeEqual(ok, true)

	// Test. 2R, Top.
	c = _test_prepare_AB_cache(aTest)
	c.top.unlink() // AB -> B.
	ok = _test_ensure_order_1_record(c, "B", "2")
	aTest.MustBeEqual(ok, true)

	// Test. 2R, Bottom.
	c = _test_prepare_AB_cache(aTest)
	c.bottom.unlink() // AB -> A.
	ok = _test_ensure_order_1_record(c, "A", "1")
	aTest.MustBeEqual(ok, true)

	// Test. 3R, Top.
	c = _test_prepare_ABC_cache(aTest)
	c.top.unlink() // ABC -> BC.
	ok = _test_ensure_order_2_records(c, [2]string{"B", "C"}, [2]string{"2", "3"})
	aTest.MustBeEqual(ok, true)

	// Test. 3R, Middle.
	c = _test_prepare_ABC_cache(aTest)
	c.top.lowerRecord.unlink() // ABC -> AC.
	ok = _test_ensure_order_2_records(c, [2]string{"A", "C"}, [2]string{"1", "3"})
	aTest.MustBeEqual(ok, true)

	// Test. 3R, Bottom.
	c = _test_prepare_ABC_cache(aTest)
	c.bottom.unlink() // ABC -> AB.
	ok = _test_ensure_order_2_records(c, [2]string{"A", "B"}, [2]string{"1", "2"})
	aTest.MustBeEqual(ok, true)
}
