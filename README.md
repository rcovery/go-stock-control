# Índice

- [Detalhes e decisões da minha submissão](#detalhes-e-decisões-da-minha-submissão)
  - [Estrutura](#estrutura)
  - [Cálculo de urgência e prioridade](#cálculo-de-urgência-e-prioridade)
- [Instruções para rodar localmente](#instruções-para-rodar-localmente)
- [Exemplos de requisição](#exemplos-de-requisição)
  - [Criar peça](#criar-peça)
  - [Listar peças](#listar-peças)
  - [Listar peças por prioridade de reposição](#listar-peças-por-prioridade-de-reposição)
  - [Atualizar peça](#atualizar-peça)
  - [Remover peça](#remover-peça)
- [Testes automatizados](#testes-automatizados)

## Detalhes e decisões da minha submissão

Quando comecei a desenvolver a solução para o desafio, pensei em várias estratégias e  arquiteturas para criar o serviço, optei por jogar fora algumas delas para focar em uma abordagem mais pragmática (explico isso mais pra baixo).

#### Estrutura
A estrutura do projeto está dividida da seguinte maneira:
- **Controller**: `internal/http/handlers`
- **Service**: `internal/part/service/`
- **Domain** (entidade, regras de negócio):
  `internal/part/part.go`, `restock.go`, `restock_priority.go`
- **Repository** (interface):
  `internal/part/repository.go`
- **Implementação do repository**: `internal/part/repository/*`

Criei um módulo "part" (dentro de internal) de uma maneira idiomática, que representa o que é uma "peça" neste projeto, então toda a lógica necessária para gerenciar as peças está isolada (seguindo a essência de DDD) dentro desse módulo.

Optei por essa estrutura justamente por ser direto ao ponto e por ter somente 1 recurso nesse serviço. Em aplicações maiores, provavelmente eu seguiria uma variação mais escalável disso.

#### Cálculo de urgência e prioridade
No início da minha implementação eu cogitei fazer o recurso de cálculo e lista de prioridade de uma maneira mais robusta, porém, no final acabei indo por um caminho mais simples para seguir o escopo do desafio.

Para tomar a decisão, eu pensei em coisas como "A aplicação tem um fluxo grande de vendas?", "A lista de reposição será consultada em qual momento?", "Concorrência seria um problema importante?", "Consigo fazer algo robusto do modo mais econômico?", "O cálculo pode mudar?", então nas minhas ideias eu considerei coisas como custo (R$) e complexidade.

#### As ideas que eu tive foram:
- **Calcular e ordenar as peças ao fazer a listagem**: é uma abordagem mais simples (acabei seguindo com essa), mas, em um ambiente de produção com alta demanda não é o ideal por questões de performance e precisão;
- **Calcular a urgência e necessidade de reposição ao inserir/atualizar as peças**: Essa é uma abordagem interessante, mas dependendo do volume de vendas teria problema com concorrência, aí teríamos que usar locks para garantir que um dado não sobrescreva outro. Eu penso que um microsserviço deste tipo deveria ser rápido, e não travar o processo com locks;
- **Lançar um job em background para fazer o cálculo de maneira assíncrona**: A ideia é subir um worker leve, algo como um [River](https://riverqueue.com/) para gerenciar a fila. É a ideia que eu quase segui até o final, o problema é que ele dependeria do Postgres para isso, aí quebraria a regra do banco trocável;
- **RabbitMQ**: Aí já seria um canhão para matar uma formiga, isso depende muito do contexto e nesse caso poderia ficar muito mais complexo para dar manutenção;

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

O servidor deve subir em `localhost:9000` e o banco Turso (SQLite) é criado em `app.db` na raiz do projeto.

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

Arquivos dos testes:
- [Lista de prioridade](internal/part/restock_priority_test.go)
- [Cálculo de urgência](internal/part/restock_test.go) (Casos extremos incluídos)
- [Validações da API](internal/part/validation_test.go) (Obs: estoque negativo é permitido)