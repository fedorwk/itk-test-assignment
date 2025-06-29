# ITK Wallet   
Тестовое задание - разработка приложения кошелька с возможностью запроса баланса и проведения операций пополнения и списания.    
API спроектирован в соответствии с требованиями задания   
# Установка и запуск   
```
docker compose up -d
```
После выполнения команд сервер запускается на порте 8080    
**Зависимости**:   
- [Google UUID](https://github.com/google/uuid) - Генерация уникальных UUID   
- [Gin](https://github.com/gin-gonic/gin) - Роутер   
- [SQLX](https://github.com/jmoiron/sqlx) - расширение стандартной библиотеки database/sql   
- [PQ](https://github.com/lib/pq) - драйвер для PostgreSQL   
   
# Описание API
   
## Создание кошелька \| POST /api/v1/wallet/create    
**Пример запроса:**   
```
curl --request POST \
  --url http://localhost:8080/api/v1/wallet/create

```
**Ответ:**   
```
{
	"balance": "0",
	"walletId": "01fcec70-c3c3-40cf-8a9b-4976e29a50fe"
}
```
## Получение баланса кошелька \| GET /api/v1/wallets/{id}
   
**Пример запроса:**   
```
curl --request GET \
  --url http://localhost:8080/api/v1/wallets/00000000-0000-0000-0000-000000000000

```
Пример ответа:   
```
{
	"balance": "1500",
	"walletId": "00000000-0000-0000-0000-000000000000"
}

```
## Изменение баланса \| POST /api/v1/wallet 
   
Параметры задаются в теле запроса в виде JSON.    
Для "operationType" поддерживается 2 значения:   
- "WITHDRAW" — списание суммы с баланса   
- "DEPOSIT" — пополнение баланса   
   
Успешный запрос возвращает пустой ответ с кодом 204   
**Пример запроса на пополнение:**   
```
curl --request POST \
  --url http://localhost:8080/api/v1/wallet \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.2.0' \
  --data '{
	"walletId": "00000000-0000-0000-0000-000000000000",
	"operationType": "DEPOSIT",
	"amount": 1000
}'
```
**Пример запроса на списание:**   
```
curl --request POST \
  --url http://localhost:8080/api/v1/wallet \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.2.0' \
  --data '{
	"walletId": "00000000-0000-0000-0000-000000000000",
	"operationType": "WITHDRAW",
	"amount": 750
}'

```
   
В случае недостатка средств возвращается ошибка с кодом 402 (Payment Required):   
```
{
	"error": "insufficient funds to withdraw"
}
```
В случае неверного указания операции возвращается ошибка с кодом 400 (Bad Request):   
```
{
	"error": "unsupported operation"
}
```
   
