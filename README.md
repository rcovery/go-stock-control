# Detalhes da minha submissão

Stack:
- Go 
- Turso (SQLite)

---

## Instruções para rodar localmente

Requisitos: [mise](https://mise.jdx.dev) ou Go 1.26

```sh
# Com mise
mise install
mise exec go -- go build -o app ./cmd/main.go
./app

# Sem mise
go build -o app ./cmd/main.go
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

> A urgência de reposição é calculada no momento da inserção/atualização.

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
É necessário ter o ID da peça em mãos

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
É necessário ter o ID da peça em mãos

```sh
curl -vX DELETE 'localhost:9000/parts/:id'
```

---

## Testes automatizados

```sh
go test -v ./...
```

Testes criados:
- Calculo de prioridade
