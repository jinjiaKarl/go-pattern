//build a program that reports the disk usage of one or more directories specified on the command line
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// concurrent walkDir
var verbose = flag.Bool("v", false, "show verbose progress messages")

func walkDir(dir string, fileSize chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			wg.Add(1)
			go walkDir(subdir, fileSize, wg)
		} else {
			fileSize <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 10) //控制一个函数可以达到的最大并发量

func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du:%v\n", err)
		return nil
	}
	return entries
}
func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	var wg sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go walkDir(root, fileSizes, &wg)
	}

	go func() {
		wg.Wait() //等待全部协程运行完，关闭channel
		close(fileSizes)
	}()
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(100 * time.Millisecond) //使用time.Tick() 会造成goroutine leak
	}
	var nfiles, nbytes int64
	//loop 退出for和select循环
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop //fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}
func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}
