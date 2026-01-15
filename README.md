# clean-go

Projeto de exemplo em Go organizado com princípios de Clean Architecture e módulos por domínio
(orders, payments, products, users). O objetivo é demonstrar separação de responsabilidades,
componentização e fluxo simples de um e-commerce (colocar pedido e processar pagamento).

## Visão geral

Este repositório contém um servidor HTTP minimalista com os seguintes módulos:

- `api/` - ponto de entrada da aplicação (`main.go`) que monta os componentes e expõe endpoints HTTP.
- `orders/` - lógica de criação de pedidos (handler + usecase + entities + repositório).
- `payments/` - integração com gateway de pagamento (p.ex. Stripe) e usecase de pagamento.
- `products/` - leitura de catálogo de produtos.
- `users/` - leitura de dados de usuário.
 - `contracts/` - módulo com definições e contratos compartilhados (tipos, DTOs e interfaces) usados por múltiplos módulos. Projetado para expor definições comuns entre `orders`, `payments`, `products` e `users`.

Cada módulo expõe um componente (factory) que é instanciado em `api/main.go`.

## Arquitetura

O projeto segue uma arquitetura limpa (Clean Architecture / Hexagonal):

- Entidades (domain) ficam em `*/internal/entities`.
- Casos de uso (application logic) em `*/usecases` (alguns módulos representam as APIs que chamam esses usecases).
- Handlers de cada módulo expõem as operações do componente.
- `api/main.go` faz a composição: cria clientes externos (MongoDB, Stripe) e instancia componentes.
 - Contratos e tipos compartilhados ficam em `contracts/` (módulo Go separado). Isso centraliza DTOs/contratos entre serviços e evita dependências cíclicas entre módulos.

## Principais arquivos e responsabilidades

- `api/main.go` - configura MongoDB usando `MONGO_URI`, carrega `STRIPE_KEY`, instancia componentes e expõe `/orders`.
- `orders/handler_place_order.go` - entrada para criar um pedido. Valida payload e chama `PlaceOrder` do componente de orders.
- `orders/internal/entities/order.go` - entidade `Order` com métodos como `calculateTotal()` e `MarkAsPaid()`.
- `payments/handler_payment.go` - handler que expõe processamento de pagamento via componente de pagamentos.
- `products/handler_catalog.go` - handler para obter produto por ID.
- `users/handler_user.go` - handler para obter dados de usuário.
 - `contracts/` - pacote com tipos/contratos compartilhados entre módulos (por exemplo: DTOs, IDs, interfaces que definem port/adapter shapes). Cada módulo consome esses contratos para manter compatibilidade entre componentes.

## Endpoints principais

- POST /orders
´´´json
{
    "user_id": "696028aff2ba343bf6310796", // temporário até implementar autenticação
    "items": [
        {
            "product_id": "6960282df2ba343bf6310793",
            "quantity": 10
        }
    ],
    "card_token": "tok_visa"
}
´´´

## Variáveis de ambiente

Defina as variáveis abaixo antes de rodar a aplicação:

- `MONGO_URI` — URI de conexão com o MongoDB (ex: `mongodb://localhost:27017`).
- `STRIPE_KEY` — chave da API Stripe (o código atual exige que esteja presente).

## Como rodar (desenvolvimento)

Rode a API localmente com Go (assumindo Go 1.20+ e `go.work` configurado):

```sh
cd clean-go
go run ./api
```

A aplicação escuta por padrão na porta `:8080` (veja `api/main.go`).

Também existe `docker-compose.yml` e `Dockerfile` no repositório para rodar via containers. Para subir com o docker-compose:

```sh
docker-compose up --build
```


