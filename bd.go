package main

import (
	"context"
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
	conn, err := pgxpool.New(context.Background(), getEnv("", ""))
	if err != nil {
		panic("Erro ao conectar ao banco de dados: " + err.Error())
	}
	println("	=> Conexão com o banco de dados estabelecida com sucesso!")
	connectionDb = &Database{Conn: conn}
	return connectionDb
}

// func selectArquivos(date string) (pgx.Rows, error) {
//     // Removido a linha que sobrescrevia o parâmetro
//     // date = "2022-08-18"

//     // Corrigido BEETWEN para BETWEEN e a sintaxe
//     sql := "SELECT * FROM Arquivo WHERE data_coluna BETWEEN '" + date + " 00:00:00' AND '" + date + " 23:59:59.999999'"
//     println(sql)
//     rows, err := connectionDb.Conn.Query(context.Background(), sql)
//     if err != nil {
//         return nil, err
//     }
//     return rows, nil
// }

func selectArquivos(date string) (pgx.Rows, error) {
	rows, err := connectionDb.Conn.Query(context.Background(), "SELECT \"Binario\",\"MimeType\", \"Hash\" FROM \"Arquivo\" WHERE \"DataCadastro\" BETWEEN '"+date+" 00:00:00' AND '"+date+" 23:59:59.999999' AND \"UrlAws\" IS NULL")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func updateArquivo(idArquivo string, bucketAws string, regionAws string, endPointAws string) (pgx.Rows, error) {
	rows, err := connectionDb.Conn.Query(context.Background(), "SELECT * FROM Arquivo LIMIT 250")
	if err != nil {
		return nil, err
	}
	return rows, nil
}
