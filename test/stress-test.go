package main

import (
	"fmt"
	"time"

	cache "github.com/vault-thirteen/Cache"
)

func main() {
	test_1A()
	test_1B()
	test_2()
	test_3()
	test_4A()
	test_4B()
	test_5()

	fmt.Println("Press 'Enter' to quit.")
	_, _ = fmt.Scanln()
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

// Test #1A. Overheating with a single record.
// Cache has 2 records; the same single record is added multiple times.
// Volume counting is enabled.
func test_1A() {
	c := cache.NewCache[string, string](100, 100_000, 3600)
	var err error

	err = c.AddRecord("B", "2")
	mustBeNoError(err)

	iMax := 10_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("A", "1")
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax
	showSummary(durTotal, reqCount)
}

// Test #1B. Overheating with a single record.
// Cache has 2 records; the same single record is added multiple times.
// Volume counting is disabled.
func test_1B() {
	c := cache.NewCache[string, string](100, 0, 3600)
	var err error

	err = c.AddRecord("B", "2")
	mustBeNoError(err)

	iMax := 10_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("A", "1")
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax
	showSummary(durTotal, reqCount)
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

	iMax := 10_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("B", "2")
		mustBeNoError(err)
		err = c.AddRecord("A", "1")
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax * 2 // A+B.
	showSummary(durTotal, reqCount)
}

// Test #3. Heating with 100 records each having 1000 bytes of data.
func test_3() {
	c := cache.NewCache[string, string](100, 100_000, 3600)
	var err error

	tmp := make([]byte, 0, 100)
	for i := 1; i <= 100; i++ {
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

	iMax := 100_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		for j := 0; j < 100; j++ {
			err = c.AddRecord(uids[j], data)
			mustBeNoError(err)
		}
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax * 100 // j=100.
	showSummary(durTotal, reqCount)
}

// Test #4A. Heating with 100 records each having 1'000'000 bytes of data.
// Volume counting is enabled.
func test_4A() {
	c := cache.NewCache[string, string](100, 1_000_000_000, 3600)
	var err error

	tmp := make([]byte, 0, 1_000_000)
	for i := 1; i <= 1_000_000; i++ {
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

	iMax := 1_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		for j := 0; j < 100; j++ {
			err = c.AddRecord(uids[j], data)
			mustBeNoError(err)
		}
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax * 100 // j=100.
	showSummary(durTotal, reqCount)
}

// Test #4B. Heating with 100 records each having 1'000'000 bytes of data.
// Volume counting is disabled.
func test_4B() {
	c := cache.NewCache[string, string](100, 0, 3600)
	var err error

	tmp := make([]byte, 0, 1_000_000)
	for i := 1; i <= 1_000_000; i++ {
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

	iMax := 1_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		for j := 0; j < 100; j++ {
			err = c.AddRecord(uids[j], data)
			mustBeNoError(err)
		}
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax * 100 // j=100.
	showSummary(durTotal, reqCount)
}

func test_5() {
	type MyClass struct {
		Name      string
		Age       int
		BirthDate time.Time
	}

	c := cache.NewCache[string, MyClass](1000, 1_000_000, 3600)
	var err error

	var x = MyClass{Name: "John", Age: 123, BirthDate: time.Now()}

	for i := 0; i < 1000; i++ {
		err = c.AddRecord("UID", x)
		mustBeNoError(err)
	}

	iMax := 1_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("UID", x)
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax
	showSummary(durTotal, reqCount)
}

func showSummary(timeElapsed time.Duration, requestsCount int) {
	reqPerSecond := float64(requestsCount) / timeElapsed.Seconds()
	fmt.Printf("Time elapsed: %f sec.; N=%d; KRPS=%.2f.\r\n", timeElapsed.Seconds(), requestsCount, reqPerSecond/1000)
}
