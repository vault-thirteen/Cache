package cache

import (
	"testing"
	"time"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewRecord(t *testing.T) {
	aTest := tester.New(t)
	var err error
	var rS *Record[string, string]
	var rBA *Record[string, []byte]

	type MyClassA struct {
		Name string
		Age  int
	}
	var rMCA *Record[string, MyClassA]

	// Test #1. checkUid fails.
	// TODO: Wait for Go language update for generics.
	//r, err = NewRecord[string, string](nil, "", "data")
	//aTest.MustBeAnError(err)
	//aTest.MustBeEqual(err.Error(), ErrUidIsEmpty)

	// Test #2. checkData fails.
	// TODO: Wait until Go language becomes a normal programming language.
	// Golang can not get real size of an object or a variable.
	//r, err = NewRecord[string, string](nil, "uid", "")
	//aTest.MustBeAnError(err)
	//aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)

	// Test #3. Data type is string.
	rS, err = NewRecord[string, string](nil, "uid", "data777", true)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(rS.uid, "uid")
	aTest.MustBeEqual(rS.data, "data777")
	aTest.MustBeEqual(rS.volume, 4+7) // This may be changed in next versions of Go language.
	aTest.MustBeEqual(rS.cache, (*Cache[string, string])(nil))
	aTest.MustBeEqual(rS.upperRecord, (*Record[string, string])(nil))
	aTest.MustBeEqual(rS.lowerRecord, (*Record[string, string])(nil))

	// Test #4. Data type is []byte.
	rBA, err = NewRecord[string, []byte](nil, "uid", []byte{0x00, 0xFF, 0x00}, true)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(rBA.volume, 4+3) // See comments above.

	// Test #5. Data type is MyClassA.
	rMCA, err = NewRecord[string, MyClassA](nil, "uid", MyClassA{Name: "John", Age: 123}, true)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(rMCA.volume, 53) // See comments above.
}

//func Test_checkUid(t *testing.T) {
//	// TODO: Wait for Go language update for generics.
//}

//func Test_checkData(t *testing.T) {
//	aTest := tester.New(t)
//	var err error
//
//	// Test #1. Bad data.
//	err = checkData[string]("")
//	aTest.MustBeAnError(err)
//	aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)
//
//	// Test #2. OK.
//	err = checkData[string]("ok")
//	aTest.MustBeNoError(err)
//}

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
	c = NewCache[string, string](0, 1000, 60)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)
	//aTest.MustBeEqual(c.volume, 1)
	aTest.MustBeEqual(c.volume, (4 + 1))

	// Test #1.
	c.top.update("333", true)
	aTest.MustBeEqual(c.top.data, "333")
	//aTest.MustBeEqual(c.volume, 3)
	aTest.MustBeEqual(c.volume, (4 + 3))
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
