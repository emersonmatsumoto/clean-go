package entities

import (
	"testing"
	"time"
)

func TestNewOrder_SetsFieldsAndTotal(t *testing.T) {
	items := []OrderItem{
		{ProductID: "p1", Price: 10.0, Quantity: 2},
		{ProductID: "p2", Price: 5.5, Quantity: 1},
	}

	before := time.Now()
	o := NewOrder("user123", items, "Rua Falsa 123")
	after := time.Now()

	if o.UserID != "user123" {
		t.Fatalf("expected UserID to be %q, got %q", "user123", o.UserID)
	}

	if len(o.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(o.Items))
	}

	if o.ShippingAddress != "Rua Falsa 123" {
		t.Fatalf("unexpected shipping address: %q", o.ShippingAddress)
	}

	if o.Status != "PENDING" {
		t.Fatalf("expected Status PENDING, got %q", o.Status)
	}

	expectedTotal := 10.0*2 + 5.5*1
	if o.Total != expectedTotal {
		t.Fatalf("expected total %v, got %v", expectedTotal, o.Total)
	}

	if o.CreatedAt.Before(before) || o.CreatedAt.After(after) {
		t.Fatalf("CreatedAt not set to now; got %v (before=%v after=%v)", o.CreatedAt, before, after)
	}
}

func TestMarkAsPaid_SetsStatusAndTransaction(t *testing.T) {
	o := &Order{Status: "PENDING"}
	o.MarkAsPaid("tx-789")
	if o.Status != "PAID" {
		t.Fatalf("expected status PAID, got %q", o.Status)
	}
	if o.TransactionID != "tx-789" {
		t.Fatalf("expected transaction id tx-789, got %q", o.TransactionID)
	}
}
