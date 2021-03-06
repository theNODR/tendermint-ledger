# Описание формата хранилища

## Прямо сейчас
В качестве хранилища используется Clickhouse.

Сервер: `http://clickhouse:8123`

Имя базы данных: `ledger`

Название таблицы: `transactions`

## Описание структуры таблицы
Название столбца | Тип данных столбца | Описание
---|---|---
create_at | DateTime | Дата/Время добавления записи транзакции в Clickhouse
server_date | UInt64 | Дата/Время (UTC, в миллисекундах) добавления добавления транзакции в Clickhouse
event_date | Date | Дата добавления транзакции транзакции в Clickhouse, вычисляется по формуле toDate(server_date/1000)
timelock | UInt64 | Дата/время (UTC, в миллисекундах). Используется для транзакций TX_OPEN_INCOME, TX_OPEN_SPEND, TX_OPEN_TRANSFER. После того, как текущий момент перейдет указанный в timelock, OPEN-транзакция считается закрытой с нулевыми результатами. Спустя какое-то врям закрывающая транзакция появится в ledger.
tx_id | string | Уникальный идентификатор транзакции
parent_tx_id | string | Уникальный идентификатор родительской транзакции, используется для закрывающих тразакций (типов TX_CLOSE_INCOME, TX_CLOSE_SPEND, TX_CLOSE_TRANSFER) и содержит уникальный идентификатор соответсвующей открывающей транзакции (типа TX_OPEN_INCOME, TX_OPEN_SPEND, TX_OPEN_TRANSFER)
version | string | Не используется
tx_type | Enum8 | Тип транзакции: <p>TX_INIT (значение 1) - начальная транзакции первода средств на адрес трекера кошелька</p><p>TX_OPEN_INCOME (2) - транзакция получения входящего адреса, на которую будут перечисляться средства за розданный траффик</p><p>TX_OPEN_SPEND (4) - транзакция получения исходящего адреса, на который трекером перечисляется фиксированная сумма, и который используется для расчетов с раздающей нодой</p><p>TX_OPEN_TRANSFER (8) - транзакция открытия канала передачи данных между пирами</p><p>TX_CLOSE_TRANSFER (16) - транзакция закрытия канала передачи между пирами, подтверждающая перевод монет между входящим и исходящим адресом</p><p>TX_CLOSE_SPEND (32) - транзакция закрытия исходящего адреса. Остаток средств на адресе переводится на адрес трекера</p><p>TX_CLOSE_INCOME (64) - транзакция закрытия входящего адреса. Остаток средств на адресе переводится на адрес трекера</p><p></p>
from | string | Адрес, с которого списываются телепорты
to | string | Адрес, на который начисляются телепорты
amount | UInt64 | Фактическое число телепортов, переведенных между адресами. Для TX_OPEN_INCOME, TX_OPEN_TRANSFER - всегда 0
volume | UInt64 | Фактическое число переданного траффика в байтах. Для TX_INIT, TX_OPEN_INCOME, TX_OPEN_TRANSFER - всегда 0. Является справочным значением.
quantum_count | UInt64 | Фактически переданное число квантов (порций) траффика. Для TX_INIT, TX_OPEN_SPEND, TX_OPEN_INCOME, TX_OPEN_TRANSFER - всегда 0.
contract_id | string | Для TX_OPEN_TRANSFER, TX_CLOSE_TRANSFER уникальный идентификатор открытого канала. Для остальных типов транзакций всегда пустая строка
price_amount | UInt64 | <p>Стоимость еденицы траффика. Используется в TX_OPEN_INCOME, TX_OPEN_SPEND, TX_OPEN_TRANSFER.</p> <p>В закрывающих транзакциях значение должно совпадать со значение из открывающей транзакции.</p>
price_quantum_power | UInt8 | <p>2 ** price_quantum_power является величиной еденицы траффика. В TX_OPEN_INCOME совместно с price_amount определяет минимальную цену, за которую пир готов раздовать траффик, в TX_OPEN_TRANSFER с price_amount определяет предложение по цене траффика, TX_OPEN_SPEND совместно с price_amount определяет максимальную цену, за которую пир готов качать траффик.</p> <p>В закрывающих транзакциях значение должно совпадать со значение из открывающей транзакции.</p>
planned_quantum_count | UInt64 | <p>В TX_OPEN_TRANSFER определяет запланированное к передаче число едениц траффика.</p> <p>В TX_CLOSE_TRANSFER должно совпадать со значение из открывающей транзакции.</p> <p> Значаение должно быть не меньше значения quantum_count.</p>
data | String | Не используется
status | UInt8 | <p>0 - транзакция не исполнена. Транзакция сохранена для истории, ее следует игнорировать</p><p>1 - транзакция исполнена.</p>
pk | string | Публичный ключ трекера, инициировавшего запись
sign | string | Base64-сигнатура транзакции. Совместно с публичным ключом используется для подтверждения истинности записи
