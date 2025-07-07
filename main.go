package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	files := make(chan int, 200)
	filesProcessed := make(chan int, 200)
	filesCompleted := make(chan int, 200)

	go addNumber(files, 100, 250) // Simula select no banco de dados

	for worker := 0; worker < runtime.GOMAXPROCS(0); worker++ {
		wg.Add(1)
		go processFiles(files, filesProcessed, &wg)
	}
	for worker := 0; worker < runtime.GOMAXPROCS(0); worker++ {
		wg.Add(1)
		go processComplete(filesProcessed, filesCompleted, &wg)
	}
	go func() {
		wg.Wait()
		close(filesProcessed)
		close(filesCompleted)
	}()
	for v := range filesCompleted {
		fmt.Println("Arquivo migrado com sucesso:", v)
	}
	fmt.Println("Fim: ", time.Since(start))
}

func addNumber(numberList chan<- int, repeticoes int, carga int) {
	var i = 1
	for r := 1; r < repeticoes; r++ {
		for j := 1; j <= carga; j++ {
			numberList <- i
			i = i + 1
		}
		println("Adicionado ao canal mais ", carga, " arquivos")
		time.Sleep(20 * time.Second) // Simula tempo de processamento
	}
	println(i)
	close(numberList)
	fmt.Println("Finalizado de adicionar arquivo", i)
}

func processFiles(files <-chan int, filesProcessed chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range files {
		fmt.Println(" * Processando arquivo:", v)
		filesProcessed <- v
		time.Sleep(250 * time.Millisecond)

	}
	fmt.Println("Arquivos processados")
}

func processComplete(filesProcessed <-chan int, filesCompleted chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range filesProcessed {
		fmt.Println(" > Atualizando arquivo no BD do Sync:", v)
		filesCompleted <- v
		time.Sleep(600 * time.Millisecond)
	}
}
