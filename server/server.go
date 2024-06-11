package server

import (
	"html/template"
	"log"
	"net/http"
	"wbtech_l0/cache"
	db "wbtech_l0/database"
	"wbtech_l0/models"

	"github.com/gorilla/mux"
)

var (
	tmpl    = template.Must(template.ParseFiles("../static/index.html"))
	storage *db.Storage
)

func getOrderFormHandler(w http.ResponseWriter, r *http.Request) { //Обработка html - шаблона
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) { //Обработка заказа на получение информации из него
	id := r.URL.Query().Get("id")

	data := struct {
		Order models.Order
		Error string
	}{
		Error: "",
	}

	if cachedOrder, found := cache.GetOrderFromCache(id); found { //Проверка на то, найден ли заказ
		if order, ok := cachedOrder.(*models.Order); ok { //Проверка на соответствие типов кэшированного заказа
			data.Order = *order //Если корректный тип данные копируются в data.Order
		} else {
			data.Error = "Invalid order data in cache"
		}
	} else { // Если заказ не найден в кэшэ, находим его в БД и сохраняем в кэш
		order, err := storage.GetOrderById(id)
		if err != nil {
			log.Printf("Error retrieving order from DB: %v", err)
			data.Error = "Order not found"
		} else {
			data.Order = *order
			cache.SaveToCache(*order)
		}
	}

	w.Header().Set("Content-Type", "text/html")   // Получаем HTML - контент
	if err := tmpl.Execute(w, data); err != nil { //Записываем данные
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartServer(s *db.Storage) {
	storage = s
	r := mux.NewRouter()
	r.HandleFunc("/", getOrderFormHandler).Methods("GET") // Обработка URL
	r.HandleFunc("/order", getOrderHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
