package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Conn *pgxpool.Pool
}

var connectionDb *Database

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func connectDb() *Database {
	conn, err := pgxpool.New(context.Background(), getEnv("DATABASE_UR", "postgres://postgres:prixpto@189.28.180.37:32768/arquivo"))
	if err != nil {
		panic("Erro ao conectar ao banco de dados: " + err.Error())
	}
	println("	=> Conexão com o banco de dados estabelecida com sucesso!")
	connectionDb = &Database{Conn: conn}
	return connectionDb
}

func selectArquivos(date string) (pgx.Rows, error) {
	rows, err := connectionDb.Conn.Query(context.Background(), "SELECT \"Binario\",\"MimeType\", \"Hash\" FROM \"Arquivo\" WHERE \"DataCadastro\" BETWEEN '"+date+" 00:00:00' AND '"+date+" 23:59:59.999999' AND \"UrlAws\" IS NULL")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func updateArquivo(hash string) error {
	// CORREÇÃO: Usar parâmetros preparados e Exec
	sql := `UPDATE "Arquivo" 
            SET "BucketAws" = $1, 
                "RegionAws" = $2, 
                "EndPointAws" = $3, 
                "UrlAws" = $4,
                "DataAlteracao" = CURRENT_TIMESTAMP
            WHERE "Hash" = $5`

	bucketAws := "sync-legado"
	regionAws := "us-east-1"
	endPointAws := "amazonaws.com"
	urlAws := "https://sync-legado.s3.amazonaws.com/" + hash

	// Usar Exec para UPDATE, não Query
	result, err := connectionDb.Conn.Exec(context.Background(), sql,
		bucketAws, regionAws, endPointAws, urlAws, hash)

	if err != nil {
		return fmt.Errorf("erro ao atualizar arquivo %s: %w", hash, err)
	}

	// Verificar se alguma linha foi atualizada
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("nenhuma linha atualizada para hash: %s", hash)
	}

	fmt.Printf("Arquivo %s atualizado com sucesso (%d linhas)\n", hash, rowsAffected)
	return nil
}
