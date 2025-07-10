package main

import (
	"fmt"
	"sync"
	"time"
)

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
			return
		}
	}
}

func sendToAws(files <-chan GetArquivo, filesProcessed chan<- GetArquivo, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range files {
		println("Enviando arquivo para o S3:", v.Hash)
		err := sendDataS3(v)
		if err != nil {
			continue
		}
		filesProcessed <- v
	}
}

func updateFilesBd(filesProcessed <-chan GetArquivo, filesCompleted chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range filesProcessed {
		updated := updateArquivo(v.Hash)
		if updated != nil {
			fmt.Print("Erro ao atualizar arquivo pno BD:", updated)
			continue
		}
		filesCompleted <- v.Hash
	}
}
