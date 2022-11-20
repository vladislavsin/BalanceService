# BalanceService

## Api
### Основная часть
#### GET balance/user/id - Получение баланса пользователя
параметры: id - uint
#### POST balance/add - Начисляет средства пользователю на баланс 
параметры: user_id - uint, amount - uint
#### POST reservation - Резервирует средства за услугу с основного баланса
параметры: user_id - uint, service_id - uint, order_id - uint, amount - uint
#### POST reservation/accept - Подтверждение о предоставленной услуги (списание из резерва)
параметры: user_id - uint, service_id - uint, order_id - uint, amount - uint

### Дополнительные задания
#### GET /transactions/user/:userId - Получение истории транзакции пользователя
параметры: id - uint, page - uint, sort - string (created_at / amount), order - string (desc / asc)
пример: /transactions/user/1?page=2&sort=amount&order=asc
