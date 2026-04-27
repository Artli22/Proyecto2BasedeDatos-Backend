package main

import (
    "database/sql"
    "fmt"
	"log" 
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func conectarDB(){
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error al configurar la conexion:", err)
    }

	if err = DB.Ping(); err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

    fmt.Println("Conexion a la base de datos exitosa")
}