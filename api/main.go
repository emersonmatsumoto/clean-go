package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/emersonmatsumoto/clean-go/orders"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
	"github.com/emersonmatsumoto/clean-go/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoUri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("Erro ao conectar no Mongo:", err)
	}
	defer client.Disconnect(ctx)

	stripeKey := os.Getenv("STRIPE_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_KEY não foi configurada nas variáveis de ambiente")
	}

	userComp := users.NewComponent(client)
	productComp := products.NewComponent(client)
	paymentComp := payments.NewComponent(stripeKey)

	orderComp := orders.NewComponent(client, productComp, paymentComp, userComp)

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var input orders.PlaceOrderInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		output, err := orderComp.PlaceOrder(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	})

	log.Println("Servidor rodando na porta :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
