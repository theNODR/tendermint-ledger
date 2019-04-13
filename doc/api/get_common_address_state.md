# Команда `get common address state`
## Описание
Команда возращает текущее количество заработанных токенов.

Число заработанных токенов возвращается в разрезе заработанных на этом леджере и заработанных на других леджерах, но связанных с текущим леджером через исходящий адрес (spend address).

## Вход
Значение поля `cmd`: 8192
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message":{
    // Common address
    // Type: string
    "address": "",
    // Peer public key
    // Type: string
    // Value: key.exportKey('pkcs8-public-pem')
    "pk": ""
  },
  // Peer public key
  // Type: string
  // Value: key.exportKey('pkcs8-public-pem')
  "pk":"",
  // RSA signature of message field in base64 encoding
  // Type: string
  // Value: key.sign(message, 'base64', 'utf8')
  "sign":""
}
```

## Пример входа
```json
{
  "message":"{\"address\":\"0x0b117d3c25485e5b28nqzlQrrZPlorgNHU/ReD5OwzBUc0AGj9NZxxhQ==\",\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\\nnQIDAQAB\\n-----END PUBLIC KEY-----\"}",
  "pk":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----",
  "sign":"G3TCie3KV3k+LGKDpw7YUcNix3kbpwYKD2hVt/CokWhBpiP1RjhtGC8STiXOKxf7vu50KV5csBYl4E0BMkrp1nSuPDDAfCW2ftP04gQpb0GEZ1bMQOA9o/8yd1mwKLWdQrEZHZCaiJJUIqBosh/qXZJRSHsJUDZTCCFztXR3i8R9fCwdGP44vUAQXmbXfzTjN8AerVIBE38shmx7uZO4VY1FIyK9kItvr/91q1vD4dl5XUOGduI+QsVfzG12Qnl9w2oXe0Hl6sVstt7N5Duy45zJgYyM2Nf9UvVXfayWtnd5CB9B3XUtJ1TufFJARPbKcdx7NDZ7nVCq3LbvbcInUg=="
}
```

## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Common address
    // Type: string
    "address": "",
    // Common address state
    // Type: Object
    "state": {
        // Число токенов, которые заработал или может заработать адрес на этом леджере
        // Type: Object
        "own": {
          // Число фактически заработанных токенов на этом адресе
          // Type: string (uint64.toString())
          "fact": "0",
          // Максимальное число токено, которые может заработать адрес в с учетом всех открытых трансфер каналов
          // Type: string (uint64.toString())
          "plan": "0"
        },
        // Число токенов, которые были или будут отслежены у адреса на этом леджере.
        // Type: Object
        "foreign": {
          // Число токенов в закрытых трансфер каналах с указанным общим адресом для чужого леджера
          // Type: string (uint64.toString())
          "fact": "0",
          // Число токенов в открытых трансфер каналах с указанным общим адресом для чежого леджера
          // Type: string (uint64.toString())
          "plan": "0"
        },
        // Число токенов, которое было фактически перечислено на адрес через финансовые транзакции
        // Type: string (uint64.toString())
        "profit": "0"
    },
    // Tracker public key
    // Type: string
    // Value: key = new NodeRSA(message)
   "pk": ""
  },
  // RSA signature of message field in base64 encoding
  // Type: string
  // Value: key.sign(message, 'base64', 'utf8')
  "sign":""
}
```

## Пример выхода
```json
{
  "message":"{\"address\":\"0x0b117d3c25485e5b28nqzlQrrZPlorgNHU/ReD5OwzBUc0AGj9NZxxhQ==\",\"state\":{\"profit\":\"0\",\"own\":{\"fact\":\"2000\",\"plan\":\"0\"},\"foreign\":{\"fact\":\"0\",\"plan\":\"0\"}},\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0BFCTcX2Va3PNFnIQfld\\ncBZjDIdyeVR4Y03Wma98sv+8FvWaRhYWCclmVBjrGQe01Mmo5p6o8jvxCvWc/qqh\\ncVsKToCHuqauRtSAcJWTvuUtbaGi0Y6Bk0E5jPXCi/PeN4kLwvCTqBwWmcASrUx1\\nR8mETbPXlXb8kctZmL4/RYc4XWiM9uu+gOJgs1fqmvpQqOw7yj3aCR7ZqLRfiwiM\\nmPT9j0uU08YNiAqUw7DDDnqAvCQ6q4bAsf4rDx46aHvLbo5DJSbhkPQrEg4pI4BI\\nlczuSGD5EB91qAs5JTQnhYrXfLrZLCKGcqer+yGwXwjAjzqL96XWd+KE0mBq+bRE\\nLwIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "sign":"swr+cNgq2dJULoN5Jmc7pw4uxdL+Ivj95ulpdsUsFqJRqlRSPHyvNLSqcrE0P7/KcPtu4EpoOrFERWhLnnEErAqNbtresiOdLKoeANYA9uTBf0oZFi+rE7vFY0KtG9lmjxDjFF+vBXe1LNEr7IK6EaQPuDiPjg+5f/q5Prmr+6mV6u3Bsyy7OiZiw91UVZCFtyDm57V+Hxy0ChAGhabODOIXy/6BvaV3elzAIAzIyqPklJ0jD92k6swvmYWCf4/tpZtqvaLxQThxEeZkrN8u3uqbiWZvE3865cL3Kr4NS1+SchL8rNp0YwmhW6lyI7dMKX+lJkHzphGnWgsaiXYbiQ=="
}
```
