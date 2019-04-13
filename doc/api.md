# Описание WebSocket

Взаимодействие с `API трекера кошелька` осуществляется при помощи WebSocket.
В среде разработки WebSocket API расположен [по адресу wss://ledger.teleport.media](https://ledger.teleport.media).

## Формат сообщения запроса
Сообщение с запросом - это JSON вида:
```json
{
  // Command code
  // Type: number
  "cmd": 0,
  // Unique identity used for connect outgoing and income messages
  // Type: string
  "corrToken": "",
  // Command-specific JSON. JSON format is specified for each command below
  // Type: JSON
  "payload": {}
}
```
Поле `corrToken` используется для связи сообщений запроса и ответа. Подразумевается, что это значение должно быть уникально в рамках клиента.
Сервер никак не проверяет это значение. То есть работа с этим значением - полностью обязанность клиента.

## Формат метаданных выполнения запроса
JSON следующего вида:
```json
// Contains timings in nanoseconds from January 1, 1970.
// It is values for determine, which part of ledger is low.
// Some of values mey be out. It means that error crash handle earlier.
{
  // Time is intermidiatially after webscoket get and handle message
  // Type: number (uint64)
  "beforeCreate": 0,
  // Time is after create handle message object
  // Type: number (uint64)
  "created": 0,
  // Time is after create response message
  // Type: number (uint64)
  "handled": 0
}
```

## Формат сообщения ответа
Сообщение ответа - это JSON вида:
```json
{
  // Correlation token value from corrToken field value of outgoing message
  // Type: string
  "corrToken": "",
  // Command-specific JSON. JSON format is specified for each command below.
  // Type: JSON
  "data": {},
  // Debug-time JSON. So it may exist, or it may unexist
  // Смотри в раздел Формат метаданных выполнения запроса
  "exData": {}
}
```

## Формат сообщения ответа с ошибкой
Сообщение ответа с ошибкой приходит в случае невозможности корреткной обработки сообщения с запросом.
Сообщение ответа с ошибкой - это JSON вида:
```json
{
  // Data for connect income and outgoing messages
  // If command is correct represent value of corrToken field in command
  // If command is incorrect represent stringify of outgoing message json
  // Type: string
  "data": "",
  // Error text
  // Type: string
  "error": "",
  // Debug-time JSON. So it may exist, or it may unexist
  // Смотри в раздел Формат метаданных выполнения запроса
  "exData": {}
}
``` 

# API трекера кошелька

1. [handshake;](./api/handshake.md)
1. [open spend channel;](./api/open_spend_channel.md)
1. [open income channel;](./api/open_income_channel.md)
1. [open transfer channel;](./api/open_transfer_channel.md)
1. [close transfer channel;](./api/close_transfer_channel.md)
1. [close income channel;](./api/close_income_channel.md)
1. [close spend channel;](./api/close_spend_channel.md)
1. [get spend address state;](./api/get_spend_address_state.md)
1. [get income address state;](./api/get_income_address_state.md)
1. [get transfer channel state;](./api/get_transfer_channel_state.md)
1. [get common address state;](./api/get_common_address_state.md)
1. [get first transactions;](./api/get_first_transactions.md)
1. [get next transactions;](./api/get_next_transactions.md)
1. [update transactions cursor.](./api/update_transactions_cursor.md)
