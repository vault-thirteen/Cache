package main

import (
	"fmt"
	"time"

	cache "github.com/vault-thirteen/Cache"
)

func main() {
	test_1()
	time.Sleep(time.Second)
	test_2()
	time.Sleep(time.Second)
	test_3()
	time.Sleep(time.Second)
	test_4()

	fmt.Println("Press 'Enter' to quit.")
	_, _ = fmt.Scanln()
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

// Test #1. Overheating with a single record.
// Cache has 2 records; the same single record is added multiple times.
func test_1() {
	c := cache.NewCache[string, string](100, 100_000, 3600)
	var err error

	err = c.AddRecord("B", "2")
	mustBeNoError(err)

	iMax := 100_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("A", "1")
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqPerSecond := float64(iMax) / durTotal.Seconds()
	fmt.Printf("Time elapsed: %f sec.; N=%d; KRPS=%.2f.\r\n",
		durTotal.Seconds(), iMax, reqPerSecond/1000)
}

// Test #2. Overheating with two switching records.
// Cache has 2 records; the same pair of records is added multiple times.
func test_2() {
	c := cache.NewCache[string, string](100, 100_000, 3600)
	var err error

	err = c.AddRecord("B", "2")
	mustBeNoError(err)
	err = c.AddRecord("A", "1")
	mustBeNoError(err)

	iMax := 100_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("B", "2")
		mustBeNoError(err)
		err = c.AddRecord("A", "1")
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqPerSecond := (2 * float64(iMax)) / durTotal.Seconds()
	fmt.Printf("Time elapsed: %f sec.; N=%d; KRPS=%.2f.\r\n",
		durTotal.Seconds(), iMax, reqPerSecond/1000)
}

// Test #3. Heating with 100 records each having 1000 bytes of data.
func test_3() {
	c := cache.NewCache[string, string](100, 100_000, 3600)
	var err error

	tmp := make([]byte, 0, 1000)
	for i := 1; i <= 1000; i++ {
		tmp = append(tmp, ' ')
	}
	data := string(tmp)

	var uids = make([]string, 100)
	for i := 0; i < 100; i++ {
		uids[i] = fmt.Sprintf("UID #%d Sample text is here bla-bla-bla hello world"+
			" ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"+
			" ppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp", i+1)
	}

	for i := 0; i < 100; i++ {
		err = c.AddRecord(uids[i], data)
		mustBeNoError(err)
	}

	iMax := 1_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		for j := 0; j < 100; j++ {
			err = c.AddRecord(uids[j], data)
			mustBeNoError(err)
		}
	}
	durTotal := time.Now().Sub(t1)
	reqPerSecond := (100 * float64(iMax)) / durTotal.Seconds()
	fmt.Printf("Time elapsed: %f sec.; N=%d; KRPS=%.2f.\r\n",
		durTotal.Seconds(), iMax, reqPerSecond/1000)
}

// Test #4. Heating with 1000 records each having 1'000'000 bytes of data.
func test_4() {
	c := cache.NewCache[string, string](1000, 1_000_000_000, 3600)
	var err error

	tmp := make([]byte, 0, 1_000_000)
	for i := 1; i <= 1_000_000; i++ {
		tmp = append(tmp, ' ')
	}
	data := string(tmp)

	var uids = make([]string, 1000)
	for i := 0; i < 1000; i++ {
		uids[i] = fmt.Sprintf("UID #%d Sample text is here bla-bla-bla hello world"+
			" ooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo"+
			" ppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp", i+1)
	}

	for i := 0; i < 1000; i++ {
		err = c.AddRecord(uids[i], data)
		mustBeNoError(err)
	}

	iMax := 100_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		for j := 0; j < 1000; j++ {
			err = c.AddRecord(uids[j], data)
			mustBeNoError(err)
		}
	}
	durTotal := time.Now().Sub(t1)
	reqPerSecond := (1000 * float64(iMax)) / durTotal.Seconds()
	fmt.Printf("Time elapsed: %f sec.; N=%d; KRPS=%.2f.\r\n",
		durTotal.Seconds(), iMax, reqPerSecond/1000)
}
