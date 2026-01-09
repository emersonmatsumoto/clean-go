package usecases

import (
	"errors"
	"testing"

	"github.com/emersonmatsumoto/clean-go/products/internal/entities"
	"github.com/emersonmatsumoto/clean-go/products/internal/ports"
)

type mockRepo struct {
	findByID func(id string) (*entities.Product, error)
}

func (m *mockRepo) FindByID(id string) (*entities.Product, error) {
	return m.findByID(id)
}

func TestGetProductUseCase_Execute(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		repoFn     func(id string) (*entities.Product, error)
		wantErr    bool
		errMessage string
		wantProd   *entities.Product
	}{
		{
			name:       "empty id",
			id:         "",
			repoFn:     nil,
			wantErr:    true,
			errMessage: "id do produto é obrigatório",
		},
		{
			name: "repo error",
			id:   "123",
			repoFn: func(id string) (*entities.Product, error) {
				return nil, errors.New("db failure")
			},
			wantErr:    true,
			errMessage: "db failure",
		},
		{
			name: "not found",
			id:   "404",
			repoFn: func(id string) (*entities.Product, error) {
				return nil, nil
			},
			wantErr:    true,
			errMessage: "produto não encontrado",
		},
		{
			name: "success",
			id:   "1",
			repoFn: func(id string) (*entities.Product, error) {
				return &entities.Product{ID: "1", Name: "Product 1", Price: 9.99}, nil
			},
			wantErr:  false,
			wantProd: &entities.Product{ID: "1", Name: "Product 1", Price: 9.99},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo ports.ProductRepository
			repo = &mockRepo{findByID: tt.repoFn}

			uc := NewGetProductUseCase(repo)
			got, err := uc.Execute(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tt.errMessage {
					t.Fatalf("expected error message %q, got %q", tt.errMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got == nil {
				t.Fatalf("expected product, got nil")
			}

			if got.ID != tt.wantProd.ID || got.Name != tt.wantProd.Name || got.Price != tt.wantProd.Price {
				t.Fatalf("expected product %+v, got %+v", tt.wantProd, got)
			}
		})
	}
}
