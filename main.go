package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func main() {

	connectionString := "postgres://postgres:keshav@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createProductTable(db)

	// product := Product{"Shoes", 5000, true}
	// pk := insertProductData(db, product)

	// var name string
	// var available bool
	// var price float64

	// query := `SELECT name, price, available FROM products WHERE id = $1`
	// err = db.QueryRow(query, 11).Scan(&name, &price, &available)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Printf("No rows were found with id: %d\n", 11)
	// 		return
	// 	}
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Name: %s\nPrice: %f\nAvailable: %t\n", name, price, available)

	data := []Product{}
	rows, err := db.Query("SELECT name, available, price FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var name string
	var available bool
	var price float64

	for rows.Next() {
		err := rows.Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}

	fmt.Println(data)
}

func createProductTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		available BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertProductData(db *sql.DB, product Product) int {
	query := `INSERT INTO products (name, price, available) VALUES ($1, $2, $3) RETURNING id`

	var pk int
	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}

	return pk
}
