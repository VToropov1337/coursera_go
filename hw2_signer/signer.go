package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(hashSignJobs ...job) {
	in := make(chan interface{}) // можно забить
	var wgexec sync.WaitGroup

	for _, Job := range hashSignJobs {
		out := make(chan interface{})

		wgexec.Add(1)
		go Worker(&wgexec, Job, in, out)
		in = out
	}
	wgexec.Wait()
}

func Worker(wgexec *sync.WaitGroup, Job job, in, out chan interface{}) {
	defer wgexec.Done()
	defer close(out)
	Job(in, out)

}

func Md5(wg *sync.WaitGroup, data string) string {
	defer wg.Done()
	r := DataSignerMd5(data)
	return r
}

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for val := range in {
		data := fmt.Sprintf("%v", val)
		crcMd5 := DataSignerMd5(data)
		wg.Add(1)
		go SingleHashJob(&wg, data, crcMd5, out)
	}
	wg.Wait()
}

func SingleHashJob(wg *sync.WaitGroup, data string, md5 string, out chan interface{}) {
	defer wg.Done()
	fmt.Println("shj", data)
	crc32Chan := make(chan string)
	Md5Chan := make(chan string)

	go func(c chan string, v string) {
		r := DataSignerCrc32(v)
		c <- r
	}(crc32Chan, data)

	go func(c chan string, v string) {
		r := DataSignerCrc32(v)
		c <- r
	}(Md5Chan, md5)

	crc32DataHash := <-crc32Chan
	Md5DataHash := <-Md5Chan
	out <- crc32DataHash + "~" + Md5DataHash

}

func MultiHash(in, out chan interface{}) {
	var wgm sync.WaitGroup
	for v := range in {
		str, ok := v.(string)
		if !ok {
			return
		}
		wgm.Add(1)
		go MultiHashJob(&wgm, out, str)
	}
	wgm.Wait()
}

func MultiHashJob(wgm *sync.WaitGroup, out chan interface{}, value string) {
	defer wgm.Done()
	var wgMulti sync.WaitGroup
	arr := make([]string, 6)
	for i := 0; i < 6; i++ {
		fmt.Println("mhj", value)
		wgMulti.Add(1)
		go func(wg *sync.WaitGroup, value string, i int, arr *[]string) {
			v := DataSignerCrc32(strconv.Itoa(i) + value)
			(*arr)[i] = v
			wg.Done()
		}(&wgMulti, value, i, &arr)
	}
	wgMulti.Wait()
	result := strings.Join(arr, "")
	out <- result

}

func CombineResults(in, out chan interface{}) {
	var arr []string
	for v := range in {
		str, ok := v.(string)
		if !ok {
			return
		}

		arr = append(arr, str)
		fmt.Println("com", "======>", runtime.NumGoroutine())
	}
	sort.Strings(arr)
	result := strings.Join(arr, "_")
	fmt.Println(result)
	out <- result
}
