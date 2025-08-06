# lab-otel-zipkin

```
lab-otel-zipkin/
├── service-a/
│   ├── main.go
│   ├── handlers/
│   ├── tracing/
│   └── Dockerfile
├── service-b/
│   ├── main.go
│   ├── .env
│   ├── handlers/
│   ├── tracing/
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

<details>
<summary>Descrição</summary>
Precisamos desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

Requisitos - Serviço A (responsável pelo input):
- O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
- O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
  - Caso seja válido, será encaminhado para o Serviço B via HTTP
  - Caso não seja válido, deve retornar:
    - Código HTTP: 422
    - Mensagem: invalid zipcode

Requisitos - Serviço B (responsável pela orquestração):
- O sistema deve receber um CEP válido de 8 digitos
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
- O sistema deve responder adequadamente nos seguintes cenários: 
  - Em caso de sucesso: 
    - Código HTTP: 200
    - Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
  - Em caso de falha, caso o CEP não seja válido (com formato correto): 
    - Código HTTP: 422
    - Mensagem: invalid zipcode
  - Em caso de falha, caso o CEP não seja encontrado: 
    - Código HTTP: 404
    - Mensagem: can not find zipcode

Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:
- Implementar tracing distribuído entre Serviço A - Serviço B
- Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

Dicas:
- Utilize a API viaCEP para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/
- Utilize a API WeatherAPI para consultar as temperaturas desejadas: https://www.weatherapi.com/
- Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
  - Sendo F = Fahrenheit
  - Sendo C = Celsius
  - Sendo K = Kelvin
- Para dúvidas da implementação do OTEL, você pode consultar https://opentelemetry.io/docs/languages/go/getting-started/
- Para implementação de spans, você pode consultar https://opentelemetry.io/docs/languages/go/instrumentation/#creating-spans
- Você precisará utilizar um serviço de collector do OTEL pode consultar https://opentelemetry.io/docs/collector/quick-start/
- Para mais informações sobre Zipkin, você pode consultar https://zipkin.io/

Importante: 
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilizar docker/docker-compose para realização dos testes da aplicação.
</details>


## Teste Local

Entre na pasta **service-b** e renomeie o arquivo `.env.exemple` para `.env` e preencha a propriedade `WEATHERAPI_KEY` com a sua chave do serviço https://www.weatherapi.com/.

### Suba os serviços

Retorne para a raiz do projeto **lab-otel-zipkin**.

Rode o seguinte comando para compilar e subir os serviços:

```shell
docker compose up --build
```

Após todos os serviços estarem prontos, rode o seguinte comando no terminal, para que seja realizada uma requisição para o **servico-a**, onde o fluxo se inicia:

```shell
curl -X POST -H "Content-Type: application/json" -d '{"cep": "29330000"}' http://localhost:8080/cep
```

Uma resposta no formato JSON, exemplificada abaixo, deve ser retornada. Além dos traces serem enviados para o Zipkin (http://127.0.0.1:9411).

```json
{"city":"Itapemirim","temp_C":19.1,"temp_F":66.38,"temp_K":292.1}
```

### Derrubar os serviços

```shell
docker compose down
```

### Comandos auxiliares

```shell
# baixar as dependências do Go
go mod tidy

# limpar o cache do docker
docker builder prune

# compilar o service-a
cd service-a
docker build -t service-a .

# compilar o service-b
cd service-b
docker build -t service-b .
```