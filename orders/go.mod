module github.com/emersonmatsumoto/clean-go/orders

go 1.25.3

replace (
	github.com/emersonmatsumoto/clean-go/payments => ../payments
	github.com/emersonmatsumoto/clean-go/products => ../products
	github.com/emersonmatsumoto/clean-go/users => ../users
)

require (
	github.com/emersonmatsumoto/clean-go/payments v0.0.0
	github.com/emersonmatsumoto/clean-go/products v0.0.0
	github.com/emersonmatsumoto/clean-go/users v0.0.0
	go.mongodb.org/mongo-driver/v2 v2.4.1
	go.opentelemetry.io/otel v1.39.0
	go.opentelemetry.io/otel/metric v1.39.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/stripe/stripe-go/v84 v84.1.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
