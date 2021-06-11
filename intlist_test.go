package concurrentorderdlist

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var benchmarkArr []int

func init() {
	rand.Seed(time.Now().UnixNano())
	size := 1000
	benchmarkArr = make([]int, size)
	for i := range benchmarkArr {
		benchmarkArr[i] = i
	}
	shuffle(benchmarkArr)
}

func BenchmarkIntList_General(b *testing.B) {
	size := len(benchmarkArr)
	l := NewInt()
	var wg sync.WaitGroup
	//wg.Add(size * 2)
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(idx int) {
			l.Insert(idx)
			wg.Done()
		}(i)
	}
	for i := 0; i < size; i++ {
		go func(idx int) {
			l.Delete(idx)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestCommonTest(t *testing.T) {
	size := 1000
	arr := make([]int, size)
	for i := range arr {
		arr[i] = i
	}
	shuffle(arr)
	l := NewInt()
	var wg sync.WaitGroup
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(idx int) {
			if ok := l.Insert(idx); !ok {
				t.Errorf("insert %d failed", idx)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	res2 := []int{}
	l.Range(func(v int) bool {
		res2 = append(res2, v)
		return true
	})
	sort.Ints(arr)
	if !reflect.DeepEqual(arr, res2) {
		t.Errorf("expect \n%v, \nbut got \n%v", arr, res2)
	}
	if !reflect.DeepEqual(len(arr), l.Len()) {
		t.Errorf("expect len is %d, but got %d", len(arr), l.Len())
	}

	var wg2 sync.WaitGroup
	wg2.Add(size)
	for i := 0; i < size; i++ {
		go func(idx int) {
			l.Delete(idx)
			wg2.Done()
		}(i)
	}
	wg2.Wait()
	res2 = []int{}
	l.Range(func(v int) bool {
		res2 = append(res2, v)
		return true
	})
	if !reflect.DeepEqual([]int{}, res2) {
		t.Errorf("expect [], but got %v", res2)
	}
	if !reflect.DeepEqual(0, l.Len()) {
		t.Errorf("expect len is %d, but got %d", 0, l.Len())
	}

}

func TestFunctional(t *testing.T) {
	l := NewInt()

	addItems := func(values ...int) {
		for _, v := range values {
			if ok := l.Insert(v); !ok {
				//t.Errorf("failed to insert value %d in array %v", v, values)
			}
		}
	}

	deleteItems := func(values ...int) {
		for _, v := range values {
			l.Delete(v)
			//	if ok := l.Delete(v); !ok {
			//		t.Errorf("failed to insert value %d in array %v", v, values)
			//	}
		}
	}

	rand.Seed(time.Now().Unix())
	arr1 := []int{1, 4, 4, 7, 10, 13, 16}
	shuffle(arr1)
	arr2 := []int{2, 5, 8, 11, 11, 14, 17}
	shuffle(arr2)
	arr3 := []int{3, 6, 9, 9, 12, 15, 18}
	shuffle(arr3)

	//arr1 := []int{10, 1, 13, 7, 4, 16, 4}
	//arr2 := []int{14, 11, 5, 11, 8, 17, 2}
	//arr3 := []int{9, 18, 3, 9, 6, 15, 12}

	var wg1 sync.WaitGroup
	wg1.Add(3)
	go func() {
		addItems(arr1...)
		wg1.Done()
	}()
	go func() {
		addItems(arr2...)
		wg1.Done()
	}()

	go func() {
		addItems(arr3...)
		wg1.Done()
	}()
	wg1.Wait()

	expectedArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}

	res := []int{}
	l.Range(func(v int) bool {
		res = append(res, v)
		return true
	})

	if !reflect.DeepEqual(expectedArr, res) {
		t.Errorf("expect %v \nbut got %v\n %v\n%v\n%v\n", expectedArr, res, arr1, arr2, arr3)
	}
	if !reflect.DeepEqual(len(expectedArr), l.Len()) {
		t.Errorf("expect len is %d, but got %d", len(expectedArr), l.Len())
	}

	_ = deleteItems
	// now delete something and add something at the same time
	//arr4 := []int{4, 6, 2, 1, 7, 10, 15}
	//arr5 := []int{4, 6, 2, 1}
	//shuffle(arr4)
	//shuffle(arr5)
	//
	//var wg2 sync.WaitGroup
	//wg2.Add(3)
	//go func() {
	//	deleteItems(arr4[:3]...)
	//	wg2.Done()
	//}()
	//go func() {
	//	deleteItems(arr4[3:]...)
	//	wg2.Done()
	//}()
	//go func() {
	//	addItems(arr5...)
	//	wg2.Done()
	//}()
	//wg2.Wait()
	//
	//expectedArr2 := []int{1, 2, 3, 4, 5, 6, 8, 9, 11, 12, 13, 14, 16, 17, 18}
	//res2 := []int{}
	//l.Range(func(v int) bool {
	//	res2 = append(res2, v)
	//	return true
	//})
	//if !reflect.DeepEqual(expectedArr2, res2) {
	//	t.Errorf("expect %v, but got %v", expectedArr2, res2)
	//}
	//if !reflect.DeepEqual(len(expectedArr2), l.Len()) {
	//	t.Errorf("expect len is %d, but got %d", len(expectedArr2), l.Len())
	//}

}

func shuffle(arr []int) {
	rand.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

func TestInsertConcurrnet(t *testing.T) {
	l := NewInt()
	var expectedArr2 []int
	size := 1000
	for i := 0; i < size; i++ {
		expectedArr2 = append(expectedArr2, i)
	}
	shuffle(expectedArr2)
	var wg sync.WaitGroup
	wg.Add(len(expectedArr2))
	for i := 0; i < len(expectedArr2); i++ {
		go func(idx int) {
			l.Insert(expectedArr2[idx])
			wg.Done()
		}(i)
	}
	wg.Wait()
	sort.Ints(expectedArr2)
	res2 := []int{}
	l.Range(func(v int) bool {
		res2 = append(res2, v)
		return true
	})
	if !reflect.DeepEqual(expectedArr2, res2) {
		t.Errorf("expect %v, but got %v", expectedArr2, res2)
	}
	if l.Len() != len(expectedArr2) {
		t.Errorf("expect len is %d, but got %d\n", len(expectedArr2), l.Len())
	}
}

//go:nosplit
func fastrandn(n uint32) uint32 {
	return uint32(rand.Int63n(int64(n)))
}

func TestIntSet(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// Correctness.
	l := NewInt()

	if l.Len() != 0 {
		t.Fatal("invalid length")
	}
	if l.Contains(0) {
		t.Fatal("invalid contains")
	}
	if l.Delete(0) {
		t.Fatal("invalid delete")
	}

	if !l.Insert(0) || l.Len() != 1 {
		t.Fatal("invalid insert")
	}
	if !l.Contains(0) {
		t.Fatal("invalid contains")
	}
	if !l.Delete(0) || l.Len() != 0 {
		t.Fatal("invalid delete")
	}

	if !l.Insert(20) || l.Len() != 1 {
		t.Fatal("invalid insert")
	}
	if !l.Insert(22) || l.Len() != 2 {
		t.Fatal("invalid insert")
	}
	if !l.Insert(21) || l.Len() != 3 {
		t.Fatal("invalid insert")
	}

	var i int
	l.Range(func(score int) bool {
		if i == 0 && score != 20 {
			t.Fatal("invalid range")
		}
		if i == 1 && score != 21 {
			t.Fatal("invalid range")
		}
		if i == 2 && score != 22 {
			t.Fatal("invalid range")
		}
		i++
		return true
	})

	i = 0
	l.Range(func(_ int) bool {
		i++
		return i != 2
	})
	if i != 2 {
		t.Fatal("invalid range")
	}

	if !l.Delete(21) || l.Len() != 2 {
		t.Fatal("invalid delete")
	}

	i = 0
	l.Range(func(score int) bool {
		if i == 0 && score != 20 {
			t.Fatal("invalid range")
		}
		if i == 1 && score != 22 {
			t.Fatal("invalid range")
		}
		i++
		return true
	})

	const num = 10000
	// Make rand shuffle array.
	// The testArray contains [1,num]
	testArray := make([]int, num)
	testArray[0] = num + 1
	for i := 1; i < num; i++ {
		// We left 0, because it is the default score for head and tail.
		// If we check the skiplist contains 0, there must be something wrong.
		testArray[i] = int(i)
	}
	for i := len(testArray) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := fastrandn(uint32(i + 1))
		testArray[i], testArray[j] = testArray[j], testArray[i]
	}

	// Concurrent insert.
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			l.Insert(testArray[i])
			wg.Done()
		}()
	}
	wg.Wait()
	if l.Len() != num {
		t.Fatalf("invalid length expected %d, got %d", num, l.Len())
	}

	//printArr(l)

	// Don't contains 0 after concurrent insertion.
	if l.Contains(0) {
		t.Fatal("contains 0 after concurrent insertion")
	}

	// Concurrent contains.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Contains(testArray[i]) {
				wg.Done()
				panic(fmt.Sprintf("insert doesn't contains %d", i))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Concurrent delete.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Delete(testArray[i]) {
				wg.Done()
				panic(fmt.Sprintf("can't delete %d", i))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if l.Len() != 0 {
		t.Fatalf("invalid length expected %d, got %d", 0, l.Len())
	}

	// Test all methods.
	const smallRndN = 1 << 8
	for i := 0; i < 1<<16; i++ {
		wg.Add(1)
		go func() {
			r := fastrandn(num)
			if r < 333 {
				l.Insert(int(fastrandn(smallRndN)) + 1)
			} else if r < 666 {
				l.Contains(int(fastrandn(smallRndN)) + 1)
			} else if r != 999 {
				l.Delete(int(fastrandn(smallRndN)) + 1)
			} else {
				var pre int
				l.Range(func(score int) bool {
					if score <= pre { // 0 is the default value for header and tail score
						panic("invalid content")
					}
					pre = score
					return true
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Correctness 2.
	var (
		x     = NewInt()
		y     = NewInt()
		count = 10000
	)

	for i := 0; i < count; i++ {
		x.Insert(i)
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			x.Range(func(score int) bool {
				if x.Delete(score) {
					if !y.Insert(score) {
						panic("invalid insert")
					}
				}
				return true
			})
			wg.Done()
		}()
	}
	wg.Wait()
	if x.Len() != 0 || y.Len() != count {
		t.Fatal("invalid length")
	}

	// Concurrent Insert and Delete in small zone.
	x = NewInt()
	var (
		insertcount uint64 = 0
		deletecount uint64 = 0
	)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				if fastrandn(2) == 0 {
					if x.Delete(int(fastrandn(10))) {
						atomic.AddUint64(&deletecount, 1)
					}
				} else {
					if x.Insert(int(fastrandn(10))) {
						atomic.AddUint64(&insertcount, 1)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if insertcount < deletecount {
		panic("invalid count")
	}
	if insertcount-deletecount != uint64(x.Len()) {
		panic("invalid count")
	}
}

func printArr(i List) {
	count := 0
	i.Range(func(value int) bool {
		fmt.Printf("%d ", value)
		count++
		if count%20 == 0 {
			fmt.Println()
		}
		return true
	})
}
