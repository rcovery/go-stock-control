# Índice

- [Detalhes da minha submissão](#detalhes-da-minha-submissão)
- [Instruções para rodar localmente](#instruções-para-rodar-localmente)
- [Exemplos de requisição](#exemplos-de-requisição)
  - [Criar peça](#criar-peça)
  - [Listar peças](#listar-peças)
  - [Listar peças por prioridade de reposição](#listar-peças-por-prioridade-de-reposição)
  - [Atualizar peça](#atualizar-peça)
  - [Remover peça](#remover-peça)
- [Testes automatizados](#testes-automatizados)

## Detalhes da minha submissão

Criei um módulo "part" (dentro de internal) que serve para representar o que é uma "peça" neste projeto, então toda a lógica necessária para gerenciar as peças está isolada de uma maneira simples dentro desse módulo.

A arquitetura do projeto está dividida da seguinte maneira:
- **Controller**: `internal/http/handlers`
- **Service**: `internal/part/service/`
- **Domain** (entidade, regras de negócio):
  `internal/part/part.go`, `restock.go`, `restock_priority.go`
- **Repository** (interface):
  `internal/part/repository.go`
- **Implementação do repository**: `internal/part/repository/*`

No início da minha implementação eu cogitei fazer o recurso de cálculo e lista de prioridade de uma maneira mais robusta para ambientes em produção, porém, no final acabei indo por um caminho mais simples para seguir o escopo do desafio.

---

## Instruções para rodar localmente

Requisitos: [mise](https://mise.jdx.dev) ou Go 1.26

```sh
# Com mise
mise install
mise exec go -- go build -o app ./cmd/api/main.go
./app

# Sem mise
go build -o app ./cmd/api/main.go
./app
```

O servidor deve subir em `0.0.0.0:9000` e o banco Turso (SQLite) é criado em `app.db` na raiz do projeto.

---

## Exemplos de requisição

### Criar peça

```sh
curl -vL 'localhost:9000/parts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Filtro de Óleo X",
    "category": "engine",
    "currentStock": 15,
    "minimumStock": 20,
    "averageDailySales": 4,
    "leadTimeDays": 5,
    "unitCost": 18.50,
    "criticalityLevel": 3
}'
```

### Listar peças

```sh
curl 'localhost:9000/parts'
```

Filtrando por categoria:

```sh
curl 'localhost:9000/parts?category=engine'
```

### Listar peças por prioridade de reposição

```sh
curl 'localhost:9000/restock/priorities'
```

### Atualizar peça

```sh
curl -vX PUT 'localhost:9000/parts/:id' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Filtro de Óleo X",
    "category": "engine",
    "currentStock": 30,
    "minimumStock": 20,
    "averageDailySales": 4,
    "leadTimeDays": 5,
    "unitCost": 18.50,
    "criticalityLevel": 3
}'
```

### Remover peça

```sh
curl -vX DELETE 'localhost:9000/parts/:id'
```

---

## Testes automatizados

```sh
go test -v ./...
```

Criei testes unitários para o cálculo de prioridade e os seguintes edge cases:
- Não precisa de reposição
  - Estoque projetado exatamente no mínimo
  - Vendas diárias zero
  - Lead time zero
- Precisa de reposição
  - Estoque negativo
  - Lead time alto
  - Criticidade alta

Também testei casos de desempate na prioridade de acordo com as regras.
