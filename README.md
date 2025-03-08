# goexpert-ratelimiter
Projeto do Desafio Técnico "Rate Limiter" da Pós-GoExpert(FullCycle).


## O desafio
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.
- Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
- Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
- As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.



## Como rodar o projeto: manual
``` shell
## 1. Suba os containers
docker-compose up -d

## 2. Teste os cenários
## Limitação por IP
for i in {1..4}; do curl -is -w "Request $i: %{http_code}\n" -o /dev/null "http://localhost:8080/api/v1/zipcode/01001001"; done
echo "wait for block duration: 5s" && sleep 5
curl -is -w "status: %{http_code}\n" -o /dev/null http://localhost:8080/api/v1/zipcode/01001001

## Limitação por token
for i in {1..6}; do curl -is -w "Request $i: %{http_code}\n" -o /dev/null -H "API_KEY: my-token" http://localhost:8080/api/v1/zipcode/01001001; done
echo "wait for block duration: 5s" && sleep 5
curl -is -w "status: %{http_code}\n" -o /dev/null -H "API_KEY: my-token" "http://localhost:8080/api/v1/zipcode/01001001"
```



## Como rodar o projeto: make
``` shell
## 1. Suba os containers
make up

## 2. Teste os cenários
## Limitação por IP
make run SCENARIO=ip

## Limitação por token
make run SCENARIO=token
```



## Funcionalidades da Linguagem Utilizadas
- web-frameworks: go-chi
- envs: godotenv
- middlewares



## Requisitos: implementação
- [x] O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- [x] O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- [x] O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- [x] As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- [x] Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- [x] O sistema deve responder adequadamente quando o limite é excedido:
    - Código HTTP: 429
    - Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
-  [x] Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
-  [x] Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
-  [x] A lógica do limiter deve estar separada do middleware.

## Requisitos: entrega
- [x] O código-fonte completo da implementação.
- [x] Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
- [x] Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
- [x] Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- [x] O servidor web deve responder na porta 8080.

