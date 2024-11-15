package cache

import "github.com/vault-thirteen/auxie/tester"

func _test_ensure_order_3_records[U UidType](
	cache *Cache[U, string],
	uids [3]U, // [0] is top, [2] is bottom.
	datas [3]string, // [0] is top, [2] is bottom.
) (ok bool) {
	var r1, r2, r3 *Record[U, string]

	r1 = cache.top
	r2 = r1.lowerRecord
	r3 = cache.bottom

	// Check links & order.
	if (r1.upperRecord != nil) || (r1.lowerRecord != r2) ||
		(r2.upperRecord != r1) || (r2.lowerRecord != r3) ||
		(r3.upperRecord != r2) || (r3.lowerRecord != nil) {
		return false
	}

	// Check UIDs.
	if (r1.uid != uids[0]) ||
		(r2.uid != uids[1]) ||
		(r3.uid != uids[2]) {
		return false
	}

	// Check Data.
	if (r1.data != datas[0]) ||
		(r2.data != datas[1]) ||
		(r3.data != datas[2]) {
		return false
	}

	// Check Map.
	if (cache.recordsByUid[uids[0]].data != datas[0]) ||
		(cache.recordsByUid[uids[1]].data != datas[1]) ||
		(cache.recordsByUid[uids[2]].data != datas[2]) ||
		(len(cache.recordsByUid) != 3) {
		return false
	}

	return true
}

func _test_ensure_order_2_records[U UidType](
	cache *Cache[U, string],
	uids [2]U, // [0] is top, [1] is bottom.
	datas [2]string, // [0] is top, [1] is bottom.
) (ok bool) {
	var r1, r2 *Record[U, string]

	r1 = cache.top
	r2 = cache.bottom

	// Check links & order.
	if (r1.upperRecord != nil) || (r1.lowerRecord != r2) ||
		(r2.upperRecord != r1) || (r2.lowerRecord != nil) {
		return false
	}

	// Check UIDs.
	if (r1.uid != uids[0]) ||
		(r2.uid != uids[1]) {
		return false
	}

	// Check Data.
	if (r1.data != datas[0]) ||
		(r2.data != datas[1]) {
		return false
	}

	// Check Map.
	if (cache.recordsByUid[uids[0]].data != datas[0]) ||
		(cache.recordsByUid[uids[1]].data != datas[1]) ||
		(len(cache.recordsByUid) != 2) {
		return false
	}

	return true
}

func _test_ensure_order_1_record[U UidType](
	cache *Cache[U, string],
	uid U,
	data string,
) (ok bool) {
	var r1 *Record[U, string]

	r1 = cache.top

	if cache.top != cache.bottom {
		return false
	}

	// Check links & order.
	if (r1.upperRecord != nil) || (r1.lowerRecord != nil) {
		return false
	}

	// Check UIDs.
	if r1.uid != uid {
		return false
	}

	// Check Data.
	if r1.data != data {
		return false
	}

	// Check Map.
	if (cache.recordsByUid[uid].data != data) ||
		(len(cache.recordsByUid) != 1) {
		return false
	}

	return true
}

func _test_ensure_order_0_records[U UidType](
	cache *Cache[U, string],
) (ok bool) {

	if (cache.top != nil) ||
		(cache.bottom != nil) {
		return false
	}

	// Check Map.
	if len(cache.recordsByUid) != 0 {
		return false
	}

	return true
}

func _test_prepare_ABC_cache(aTest *tester.Test) (c *Cache[string, string]) {
	var err error
	c = NewCache[string, string](0, 0, 60)
	err = c.AddRecord("C", "3")
	aTest.MustBeNoError(err)
	err = c.AddRecord("B", "2")
	aTest.MustBeNoError(err)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)
	return c
}

func _test_prepare_ABC_cache_with_low_ttl(aTest *tester.Test) (c *Cache[string, string]) {
	var err error
	c = NewCache[string, string](0, 0, 3)
	err = c.AddRecord("C", "3")
	aTest.MustBeNoError(err)
	err = c.AddRecord("B", "2")
	aTest.MustBeNoError(err)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)
	return c
}

func _test_prepare_AB_cache(aTest *tester.Test) (c *Cache[string, string]) {
	var err error
	c = NewCache[string, string](0, 1000, 60)
	err = c.AddRecord("B", "2")
	aTest.MustBeNoError(err)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)
	return c
}

func _test_prepare_A_cache(aTest *tester.Test) (c *Cache[string, string]) {
	var err error
	c = NewCache[string, string](0, 0, 60)
	err = c.AddRecord("A", "1")
	aTest.MustBeNoError(err)
	return c
}

func _test_prepare_0_cache() (c *Cache[string, string]) {
	c = NewCache[string, string](0, 0, 60)
	return c
}

func _test_prepare_record_Q(aTest *tester.Test, c *Cache[string, string]) (r *Record[string, string]) {
	var err error
	r, err = NewRecord(c, "Q", "W", true)
	aTest.MustBeNoError(err)
	return r
}
