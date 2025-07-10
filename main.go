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
	filesInAws := make(chan GetArquivo, 100000)
	filesCompleted := make(chan string, 100000)

	connectDb()
	connectAws()

	go getFiles(files)

	for worker := 0; worker < runtime.GOMAXPROCS(0); worker++ {
		wg.Add(1)
		go sendToAws(files, filesInAws, &wg)
	}
	go func() {
		wg.Wait()
		close(filesInAws)
	}()

	for worker := 0; worker < runtime.GOMAXPROCS(0); worker++ {
		wg2.Add(1)
		go updateFilesBd(filesInAws, filesCompleted, &wg2)
	}
	go func() {
		wg2.Wait()
		close(filesCompleted)
	}()

	for v := range filesCompleted {
		fmt.Println("Arquivo migrado com sucesso: " + v)
	}
	fmt.Println("Fim: ", time.Since(start))
}
