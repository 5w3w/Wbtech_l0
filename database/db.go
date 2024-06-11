package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	models "wbtech_l0/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "swew"
	password = "1234"
	dbname   = "wbtech_l0"
)

func InitDB() (*sql.DB, error) {
	var err error
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode = disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	return db, nil
}

func (s *Storage) GetOrderById(id string) (*models.Order, error) {
	query := "SELECT uid, data FROM orders WHERE uid = $1"
	row := s.Db.QueryRow(query, id) //получаем данные из DB
	var uid string
	var data []byte
	var order models.Order
	err := row.Scan(&uid, &data) //Копируем данные в переменные uid, data
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no order found with id %s", id)
		}
		return nil, err
	}
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}
	order.ID = uid
	return &order, nil
}

func (s *Storage) SaveOrder(orderJSON []byte) error {

	var order models.Order
	err := json.Unmarshal(orderJSON, &order)
	if err != nil {
		return err
	}

	query := `INSERT INTO orders (data) VALUES ($1) RETURNING uid`
	var uid int
	err = s.Db.QueryRow(query, orderJSON).Scan(&uid)
	if err != nil {
		return err
	}

	log.Printf("Order added to DB with UID: %v", uid)
	return nil
}
