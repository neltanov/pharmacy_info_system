package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

type QueryRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params"`
}

type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type Customer struct {
	ID          int    `json:"id"`
	Surname     string `json:"surname"`
	Name        string `json:"name"`
	MiddleName  string `json:"middle_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

type Patient struct {
	ID         int    `json:"id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	MiddleName string `json:"middle_name"`
	Age        int    `json:"age"`
	Diagnosis  string `json:"diagnosis"`
}

type Doctor struct {
	ID         int    `json:"id"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	MiddleName string `json:"middle_name"`
}

type Receipt struct {
	ID        int `json:"id"`
	DoctorID  int `json:"doctor_id"`
	PatientID int `json:"patient_id"`
}

type Order struct {
	ID             int    `json:"id"`
	CustomerID     int    `json:"customer_id"`
	ReceiptID      int    `json:"receipt_id"`
	OrderDate      string `json:"order_date"`
	ProductionDate string `json:"production_date"`
	Status         string `json:"status"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := getEnv("DATABASE_USER", "")
	if dbUser == "" {
		log.Fatal("DATABASE_USER environment variable is required")
	}
	dbPasswordRaw := getEnv("DATABASE_PASSWORD", "")
	if dbPasswordRaw == "" {
		log.Fatal("DATABASE_PASSWORD environment variable is required")
	}
	dbHost := getEnv("DATABASE_HOST", "localhost")
	dbPort := getEnv("DATABASE_PORT", "5432")
	dbName := getEnv("DATABASE_NAME", "postgres")

	dbPassword := url.QueryEscape(dbPasswordRaw)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/query_names", queryNamesHandler).Methods("GET")

	r.HandleFunc("/queries/{query}", executeQuery).Methods("GET")
	r.HandleFunc("/query", queryHandler).Methods("POST")
	r.HandleFunc("/orders", getOrdersHandler).Methods("GET")
	r.HandleFunc("/orders", createOrderHandler).Methods("POST")
	r.HandleFunc("/orders/{id}", updateOrderHandler).Methods("PUT")
	r.HandleFunc("/orders/{id}", deleteOrderHandler).Methods("DELETE")

	r.HandleFunc("/create_customer", createCustomer).Methods("POST")
	r.HandleFunc("/create_patient", createPatient).Methods("POST")
	r.HandleFunc("/create_doctor", createDoctor).Methods("POST")
	r.HandleFunc("/create_receipt", createReceipt).Methods("POST")
	r.HandleFunc("/create_order", createOrder).Methods("POST")

	port := getEnv("SERVER_PORT", "8000")

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var medicineID int
	var medicineType string
	var productionDate time.Time
	var status string

	rows, err := db.Query(context.Background(), `
		SELECT m.id, m.type 
		FROM medicine_list ml
		JOIN medicine m ON ml.medicine_id = m.id
		WHERE ml.receipt_id = $1`, order.ReceiptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orderDate, err := time.Parse("2006-01-02", order.OrderDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allMedicinesAvailable := true
	for rows.Next() {
		if err := rows.Scan(&medicineID, &medicineType); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var totalAmount int
		err := db.QueryRow(context.Background(), `SELECT total_amount FROM medicine_warehouse WHERE medicine_id = $1`, medicineID).Scan(&totalAmount)
		if err != nil || totalAmount <= 0 {
			allMedicinesAvailable = false
			if medicineType == "local_medicine" {
				var productionTime string
				err = db.QueryRow(context.Background(), `SELECT pt.time_to_product FROM local_medicine lm JOIN production_techonology pt ON lm.production_techology = pt.id WHERE lm.medicine_id = $1`, medicineID).Scan(&productionTime)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				duration, err := time.ParseDuration(productionTime)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				productionDate = orderDate.Add(duration)
			} else {
				productionDate = orderDate.AddDate(0, 0, 7) // неделя после начала заказа
			}
		}
	}

	if allMedicinesAvailable {
		status = "done"
		productionDate = orderDate
	} else {
		status = "in_production"
	}

	err = db.QueryRow(context.Background(),
		`INSERT INTO orders (customer_id, receipt_id, order_date, production_date, status)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		order.CustomerID, order.ReceiptID, order.OrderDate, productionDate, status,
	).Scan(&order.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func createReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(context.Background(),
		`INSERT INTO receipt (doctor_id, patient_id)
		 VALUES ($1, $2) RETURNING id`,
		receipt.DoctorID, receipt.PatientID,
	).Scan(&receipt.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

func createDoctor(w http.ResponseWriter, r *http.Request) {
	var doctor Doctor
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(context.Background(),
		`INSERT INTO doctor (surname, name, middle_name)
		 VALUES ($1, $2, $3) RETURNING id`,
		doctor.Surname, doctor.Name, doctor.MiddleName,
	).Scan(&doctor.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctor)
}

func createPatient(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(context.Background(),
		`INSERT INTO patient (surname, name, middle_name, age, diagnosis)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		patient.Surname, patient.Name, patient.MiddleName, patient.Age, patient.Diagnosis,
	).Scan(&patient.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(context.Background(),
		`INSERT INTO customer (surname, name, middle_name, phone_number, address)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		customer.Surname, customer.Name, customer.MiddleName, customer.PhoneNumber, customer.Address,
	).Scan(&customer.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func queryNamesHandler(w http.ResponseWriter, r *http.Request) {
	names, err := loadQueryNamesFromFile("queries/query_names.txt")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading query names: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(names); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func executeQuery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := vars["query"]

	params := r.URL.Query()
	result, err := performQuery(query, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return
	}
}

func performQuery(queryID string, params url.Values) (*QueryResult, error) {
	query, err := loadQueryFromFile(fmt.Sprintf("queries/%s.sql", queryID))
	if err != nil {
		return nil, err
	}

	var rows pgx.Rows

	if len(params) > 0 {
		var args []interface{}
		args = append(args, params.Get("Тип"))
		rows, err = db.Query(context.Background(), query, args...)
	} else {
		rows, err = db.Query(context.Background(), query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()
	result := &QueryResult{
		Columns: make([]string, len(columns)),
		Rows:    make([][]interface{}, 0),
	}

	for i, col := range columns {
		result.Columns[i] = string(col.Name)
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		result.Rows = append(result.Rows, values)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func loadQueryFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func loadQueryNamesFromFile(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := req.Query
	params := req.Params

	var args []interface{}
	for _, param := range params {
		args = append(args, param)
	}

	rows, err := db.Query(context.Background(), query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query execution error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()

	results := QueryResult{
		Columns: make([]string, len(columns)),
	}

	for i, col := range columns {
		results.Columns[i] = string(col.Name)
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching row: %v", err), http.StatusInternalServerError)
			return
		}

		results.Rows = append(results.Rows, values)
	}

	if rows.Err() != nil {
		http.Error(w, fmt.Sprintf("Row iteration error: %v", rows.Err()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	result, err := performQuery("get_orders", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return
	}
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO orders (customer_id, receipt_id, order_date, production_date, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := db.QueryRow(context.Background(), query, order["customer_id"], order["receipt_id"], order["order_date"], order["production_date"], order["status"])

	var id int
	if err := row.Scan(&id); err != nil {
		http.Error(w, fmt.Sprintf("Query execution error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"id": id}); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var updatedOrder Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Формирование SQL запроса для обновления данных заказа в базе данных
	query := `
		UPDATE orders 
		SET customer_id = $1, receipt_id = $2, order_date = $3, production_date = $4, status = $5 
		WHERE id = $6
	`

	// Выполнение SQL запроса к базе данных
	_, err := db.Exec(context.Background(), query, updatedOrder.CustomerID, updatedOrder.ReceiptID, updatedOrder.OrderDate, updatedOrder.ProductionDate, updatedOrder.Status, orderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query execution error: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправка ответа с кодом статуса 204 (No Content) после успешного обновления заказа
	w.WriteHeader(http.StatusNoContent)
}

func deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	query := `DELETE FROM orders WHERE id = $1`
	_, err := db.Exec(context.Background(), query, orderID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query execution error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
