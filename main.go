package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type GetArquivo struct {
	Binario  []byte `json:"binario"`
	MimeType string `json:"mimeType"`
	Hash     string `json:"hash"`
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	files := make(chan GetArquivo, 100000)
	filesProcessed := make(chan GetArquivo, 100000)
	filesCompleted := make(chan GetArquivo, 100000)

	connectDb()

	// go addNumber(files, 100, 250) // Simula select no banco de dados
	go getFiles(files)
	for worker := 0; worker < runtime.GOMAXPROCS(0); worker++ {
		wg.Add(1)
		go sendToAws(files, filesProcessed, &wg)
	}
	go func() {
		wg.Wait()
		close(filesProcessed)
	}
	for worker := 0; worker < runtime.GOMAXPROCS(0); worker++ {
		wg.Add(1)
		go processComplete(filesProcessed, filesCompleted, &wg2)
	}
	go func(){
		wg2.Wait()
		close(filesCompleted)
	}
	// go func() {
	// 	wg.Wait()
	// 	close(filesProcessed)
	// 	close(filesCompleted)
	// }()
	go func() {
		for v := range files {
			fmt.Println(v.Hash)
		}
	}()

	fmt.Println("Fim: ", time.Since(start))
}
