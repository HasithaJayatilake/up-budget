package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/charmbracelet/log"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db = setupDatabase()
	defer db.Close()

	// s := gocron.NewScheduler(time.UTC)
	// s.Every(1).Day().At("00:00").Do(fetchAndStoreData)
	// s.StartBlocking()

	fetchTransactions()
}

func setupDatabase() *sqlx.DB {
	dbConnection := os.Getenv("DB_CONNECTION")

	pgdb, err := sqlx.Open("postgres", dbConnection)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = pgdb.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Info("Connected to the database!")

	return pgdb
}

func fetchAndStoreData() {
	bearerToken := os.Getenv("UP_PAT")
	accountId := os.Getenv("UP_ACCOUNT")
	apiUrl := `https://api.up.com.au/api/v1/accounts/` + accountId + `/transactions`

	// resp, err := http.Get("https://api.up.com.au/api/v1/accounts")
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching data from API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching data from API: status code %d", resp.StatusCode)
	}

	var apiResponse TransactionAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Fatalf("Error decoding API response: %v", err)
	}

	for _, item := range apiResponse.Data {
		log.Info(item)

	}
}

func fetchTransactions() {
	apiUrl := "https://api.up.com.au/api/v1/"
	bearerToken := os.Getenv("UP_PAT")
	accountId := os.Getenv("UP_ACCOUNT_ID")
	transactionsUrl := fmt.Sprintf("%s/accounts/%s/transactions", apiUrl, accountId)

	// Parse the base URL
	baseUrl, err := url.Parse(transactionsUrl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Add query parameters
	params := url.Values{}
	params.Add("filter[status]", "SETTLED")
	params.Add("filter[since]", "2024-06-26T00:00:00Z")
	// params.Add("filter[until]", "2020-01-01T01:02:03+10:00")
	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching data from API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching data from API: status code %d", resp.StatusCode)
	}

	var apiResponse TransactionAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Fatalf("Error decoding API response: %v", err)
	}

	for _, transaction := range apiResponse.Data {
		log.Info(transaction)
		// _, err := db.NamedExec(`INSERT INTO transactions (id, type, status, raw_text, description, message, is_categorizable, hold_amount, hold_foreign_amount, round_up, cashback, amount, foreign_amount, card_purchase_method, settled_at, created_at)
		//                         VALUES (:id, :type, :status, :raw_text, :description, :message, :is_categorizable, :hold_info.amount, :hold_info.foreign_amount, :round_up, :cashback, :amount, :foreign_amount, :card_purchase_method, :settled_at, :created_at)
		//                         ON CONFLICT (id) DO UPDATE SET type = :type, status = :status, raw_text = :raw_text, description = :description, message = :message, is_categorizable = :is_categorizable, hold_amount = :hold_info.amount, hold_foreign_amount = :hold_info.foreign_amount, round_up = :round_up, cashback = :cashback, amount = :amount, foreign_amount = :foreign_amount, card_purchase_method = :card_purchase_method, settled_at = :settled_at, created_at = :created_at`,
		// 	map[string]interface{}{
		// 		"id":                       transaction.ID,
		// 		"type":                     transaction.Type,
		// 		"status":                   transaction.Attributes.Status,
		// 		"raw_text":                 transaction.Attributes.RawText,
		// 		"description":              transaction.Attributes.Description,
		// 		"message":                  transaction.Attributes.Message,
		// 		"is_categorizable":         transaction.Attributes.IsCategorizable,
		// 		"hold_info.amount":         transaction.Attributes.HoldInfo.Amount,
		// 		"hold_info.foreign_amount": transaction.Attributes.HoldInfo.ForeignAmount,
		// 		"round_up":                 transaction.Attributes.RoundUp,
		// 		"cashback":                 transaction.Attributes.Cashback,
		// 		"amount":                   transaction.Attributes.Amount,
		// 		"foreign_amount":           transaction.Attributes.ForeignAmount,
		// 		"card_purchase_method":     transaction.Attributes.CardPurchaseMethod,
		// 		"settled_at":               transaction.Attributes.SettledAt,
		// 		"created_at":               transaction.Attributes.CreatedAt,
		// 	})
		// if err != nil {
		// 	log.Fatalf("Error updating database: %v", err)
		// }
	}
}
