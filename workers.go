package main

import (
	"fmt"
	"sync"
	"time"
)

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

func sendToAws(files <-chan GetArquivo, filesProcessed chan<- GetArquivo, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range files {
		err := sendDataS3(v)
		if err != nil {
			continue
		}
		filesProcessed <- v
	}
}

func processComplete(filesProcessed <-chan int, filesCompleted chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range filesProcessed {
		fmt.Println(" > Atualizando arquivo no BD do Sync:", v)
		filesCompleted <- v
		time.Sleep(600 * time.Millisecond)
	}
}

func getFiles(filesList chan<- GetArquivo) {
	date := time.Now()
	for {
		date = date.AddDate(0, 0, -1)

		if date.Month() == time.March {
			// panic("Valores de março não são permitidos, fim da execução")
		}
		rows, err := selectArquivos(date.Format("2006-01-02"))
		if err != nil {
			fmt.Println("Erro ao selecionar arquivos:", err)
			continue
		}
		for rows.Next() {
			var arquivo GetArquivo
			rows.Scan(&arquivo.Binario, &arquivo.MimeType, &arquivo.Hash)
			filesList <- arquivo
		}
	}
}
