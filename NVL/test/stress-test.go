package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vault-thirteen/Cache/NVL"
)

func main() {
	test_1()
	test_2()
	test_3()
	test_4()
	test_5()
	test_6()

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
	c := nvl.NewCache[string, string](100, 3600)
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
	reqCount := iMax
	showSummary(durTotal, reqCount)
}

// Test #2. Overheating with two switching records.
// Cache has 2 records; the same pair of records is added multiple times.
func test_2() {
	c := nvl.NewCache[string, string](100, 3600)
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
	reqCount := iMax * 2
	showSummary(durTotal, reqCount)
}

// Test #3. Heating with 100 records each having 1000 bytes of data.
func test_3() {
	c := nvl.NewCache[string, string](100, 3600)
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
	reqCount := iMax * 100
	showSummary(durTotal, reqCount)
}

// Test #4. Heating with 1000 records each having 1'000'000 bytes of data.
func test_4() {
	c := nvl.NewCache[string, string](1000, 3600)
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
	reqCount := iMax * 1000
	showSummary(durTotal, reqCount)
}

// Test #5. Heating with 1000 pointers to object.
func test_5() {
	var obj MyClassA
	{
		obj = MyClassA{
			Name:        "John",
			Age:         123,
			IsSomething: true,
			X:           1.1,
			Y:           2.2,
			Z:           3.3,

			wg: sync.WaitGroup{},
			ba: []byte("TestArrayOfBytes"),
			m:  new(sync.Mutex),
		}
		obj.neighbour = &obj
		obj.counter.Store(12345)
		obj.m.Lock()
		obj.m.Unlock()
	}

	c := nvl.NewCache[string, *MyClassA](1000, 3600)
	var err error

	err = c.AddRecord("J", &obj)
	mustBeNoError(err)

	iMax := 100_000_000
	t1 := time.Now()
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("X", &obj)
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax
	showSummary(durTotal, reqCount)
}

// Test #6. Heating with 1000 interfaces.
func test_6() {
	fmt.Println("Writing and reading an interface, casting its type ...")

	var obj MyClassA
	{
		obj = MyClassA{
			Name:        "John",
			Age:         123,
			IsSomething: true,
			X:           1.1,
			Y:           2.2,
			Z:           3.3,

			wg: sync.WaitGroup{},
			ba: []byte("TestArrayOfBytes"),
			m:  new(sync.Mutex),
		}
		obj.neighbour = &obj
		obj.counter.Store(12345)
		obj.m.Lock()
		obj.m.Unlock()
	}

	c := nvl.NewCache[string, IMyInterface](256, 3600) // Emulation of small memory size.
	var err error

	err = c.AddRecord("J", &obj)
	mustBeNoError(err)

	iMax := 100_000_000
	t1 := time.Now()
	var x IMyInterface
	var w *MyClassA
	for i := 1; i <= iMax; i++ {
		err = c.AddRecord("X", &obj)
		mustBeNoError(err)

		x, err = c.GetRecord("X")
		mustBeNoError(err)
		w = FromAny(x)
		if w.Age != 123 {
			err = errors.New("failure")
		}
		mustBeNoError(err)
	}
	durTotal := time.Now().Sub(t1)
	reqCount := iMax * 2 // W+R
	showSummary(durTotal, reqCount)
}

func showSummary(timeElapsed time.Duration, requestsCount int) {
	reqPerSecond := float64(requestsCount) / timeElapsed.Seconds()
	fmt.Printf("Time elapsed: %f sec.; N=%d; KRPS=%.2f.\r\n", timeElapsed.Seconds(), requestsCount, reqPerSecond/1000)
}
