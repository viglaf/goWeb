package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

// Сюда писать код
func main() {
	inputData := []int{7, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}
	ExecutePipeline(hashSignJobs...)

}
func ExecutePipeline(jobs ...job) {

	pipe := make(chan interface{})
	wg := new(sync.WaitGroup)
	wg.Add(len(jobs))
	for _, job := range jobs {
		out := make(chan interface{})
		go func() {
			job(pipe, out)
			wg.Done()
		}()
	}
	wg.Wait()

}
func SingleHash(in chan interface{}, out chan interface{}) {
	var v string
	out = make(chan interface{})
	for i := range in {
		switch s := i.(type) {
		case string:
			v = s
		case int:
			v = strconv.Itoa(s)
		default:
			fmt.Println("dataRaw is not a string or int1")
		}
		crc1 := DataSignerCrc32(v)
		crc2 := DataSignerCrc32(DataSignerMd5(v))
		result := fmt.Sprintf("%s~%s", crc1, crc2)
		out <- result
	}
	close(out)
}

func MultiHash(in chan interface{}, out chan interface{}) {
	var result, v string
	for i := range in {
		switch s := i.(type) {
		case string:
			v = s
		case int:
			v = strconv.Itoa(s)
		default:
			fmt.Println("dataRaw is not a string or int2")
		}
		for i := 0; i < 6; i++ {
			result += DataSignerCrc32(strconv.Itoa(i) + v)
		}
		out <- result
	}
	close(out)
}
func CombineResults(in chan interface{}, out chan interface{}) {
	sl := make([]string, 0, 100)
	var result string
	for i := range in {
		switch v := i.(type) {
		case string:
			sl = append(sl, v)
		default:
			fmt.Println("dataRaw is not a string or int3")
		}
		sort.Strings(sl)

		for _, s := range sl {
			result += s + "_"
		}
		out <- result
	}
	close(out)
}
