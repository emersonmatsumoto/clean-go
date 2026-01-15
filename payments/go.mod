module github.com/emersonmatsumoto/clean-go/payments

go 1.25.3

replace github.com/emersonmatsumoto/clean-go/contracts => ../contracts

require (
	github.com/emersonmatsumoto/clean-go/contracts v0.0.0
	github.com/stripe/stripe-go/v84 v84.1.0
	go.opentelemetry.io/otel v1.39.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
