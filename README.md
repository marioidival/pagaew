# Pagaew

## Stack

- Golang (version 1.9.2))
- MySQL (version 8.0)
- Docker e Docker compose


## Set up

### Migration

Pode executar a migração do banco de dados com o seguinte comando:

```sh
$ go run cmd/dbmigrate/main.go
```

ou com Docker Compose

```sh
$ docker compose up -d server-migrate
```

### Docker (Docker compose)

Para subir todo o projeto, execute o seguinte comando:

```sh
$ docker compose up -d
```

Isso fará o build da imagem da aplicação e baixará a imagem do MySQL versão 8 para o seu computador. Caso queria utilizar seu proprio server do MySQL, por favor adicione a string de conexão na variavel de ambiente `DATABASE_URL`.

### Environments Variables

No arquivo `docker-compose.yml` você pode ver os valores esperados para as seguintes variaveis nos serviços:

- `DATABASE_URL` -> é esperado o DSN MySQL do servidor. O valor padrão esta apontando para o servidor MySQL gerenciado pelo Docker.
- `ENVIRONMENT` ->  aqui é um valor simulando ambientes. No inicio da aplicação, simplesmente checamos se o valor é `prod`. O efeito dessa variavel está ligada a que tipo de repositorio será usado na aplicação: A implementação em Memoria ou usando MySQL.
- `INTERVAL` -> é o valor em segundos que o job irá buscar no banco de dados todos os `log_invoice` com status `PENDING`.

## API spec


### POST /load

Para enviar uma requisição neste endpoint, o corpo da mensagem deve ser do tipo CSV com a seguinte estrutura:

```text/csv
name,governmentId,email,debtAmount,debtDueDate,debtId
John Doe,11131111111,johndoe@kanastra.com.br,1000000.00,2022-10-12,8291
```

Pontos importantes:

- O cabeçalho desse CSV deve seguir a mesma ordem, caso contrário o endpoint retornará 400 BadRequest
- Os valores desse CSV devem ser válidos:
    - `governmentID` tem uma lógica básica onde consideramos qualquer número não repetido como válido.
    - `email` com um valor válido
    - `name` com uma string não vazia
    - `debtAmount` como um decimal válido não menor que 0.
    - `debtDueDate` com uma data válida.
    - `debtId` como uma string não vazia

### POST /webhook

Para enviar uma requisição neste endpoint, o corpo da mensagem deve ser do tipo JSON com a seguinte estrutura:

```json
{
	"debtId": "8291",
	"paidAt": "2022-06-09 10:00:00",
	"paidAmount": 100000.00,
	"paidBy": "John Doe"
}
```

Pontos importantes:
- `debtId` não existente causa um retorno 404 Not Found
- É permitido apenas o pagamento de uma vez do boleto. Mas sei que a logica aqui poderia ficar melhor, como por exemplo, entender se foi enviado um valor diferente do esperado.

