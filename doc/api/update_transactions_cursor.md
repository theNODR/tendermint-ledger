# Команда `update transactions cursor`
## Описание
Продляет время жизни курсора на сервере.
* Время жизни курсора продляется на 5 минут с текущего момента времени.

### Вход
Значение поля `cmd`: 4096
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Peer public key
    // Type: string
    // Value: key.exportKey('pkcs8-public-pem')
    "pk": "",
    // Cursor id from `get first transactions` query answer
    // Should be valid. If invalid return error
    // Type: string 
    "cursorId": ""
  },
  // Peer public key
  // Type: string
  // Value: key.exportKey('pkcs8-public-pem')
  "pk": "",
  // RSA signature of message field in base64 encoding
  // Type: string
  // Value: key.sign(message, 'base64', 'utf8')
  "sign": ""
}
```

## Пример входа
```json
{
  "message": "{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"cursorId\":\"7c9b24aa-9a84-4c68-ab64-e01c4048ac51\"}",
  "sign":"qD+cMbrWhLQjlEzzkog/vemn8sNJm4t73p/E2QTdgCnZRTgrrHhNTujR9T6/t5nc34Iu91LAtt9sLFTT/kzu7kWaaR+xmdI60uBPi8Pw65UO+a+QYMhGwtJd7UM7qtAT5O582ypp/1q7iZVscO63RMewSbzPHs/skjxOSNRn8LH8Iht+TaPKqToy9Czy0VHtWcFs2CfPSqITDDYR67HtmWefY2zEUIlns0ZJjHGE/K20+yrtW8pGrQVEHmqXYZHy+rV/3oOxoRyZcV/wgJauUEFR/YK/AxrHZZ4QCAJPRge+Z0CFThveMF2ltdjKocdCvFaBi8MU8mXwfLU3/Q2AuQ==",
  "pk":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n"
}
```

## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Tracker public key
    // Type: string
    // Value: key = new NodeRSA(message)
    "pk": ""
  },
  // RSA signature of message field in base64 encoding
  // Type: string
  // Value: key.sign(message, 'base64', 'utf8')
  "sign": ""
}
```

## Пример выхода
```json
{
  "message":"{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8vQVhLPYuLSPGUHLgboK\\n9zqvaBi97vBSofH1ukz6r3zF2MpOWrMr6NKCxvItu32aEmjX/Pnh1UAZ0n7C/sx8\\nyGk8cJ5f11y+nv2odcZHbNmUhVdxKZO6ashpznJtGhL+vOEdp4Oha74DsyxRz6SX\\nSblJlNBSrxFYb9sZ+GjrEk6FmNYhm3J8vjkVvXtOrwX7d2/kU5tiCzSYBckBwJ2b\\nMW3g0Zt332N8Kj5Nc9xIBMSDGQn1+4DvgOoetLn/9lrE9Hr6RXv3RfLgkmf3HQTo\\n1Iux1Af94D899yYdkPnpNGwNwgcytJ+sJXYYt6FZp8b6SRLpS2ZYhky92Hi/21RQ\\nqwIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "sign":"er0COvR1+BCcfTRq7BYz4OWTNCeCVIIjb8APSAkne260J+ABsLk4q1kfP0GXtq0Vv0b9VLGRp4Ykt7wdCTWHwXzZIe1buyo2iKwXAbpVYdGLMjub6XJwlGGSag2b0UGHzZLysKZpwkbUiNtbTMTwPaR8tRzdPVK048N+BBZqZEAy88Y8AcDwt/8Y40APu6L/d5A5GvUzcRau/GszG3ggH09mo0+/9TLZs1c9hVAhaXEiL6dNNEpgJObk7WxhC4wsqeD7AvlaRdOy/qicLS6qmw3GzJVdCRSj/Gj64hAHYOCpMkS1dwiDEFNzRGWCwHJY9n07jcgr8I6OFpQygJ7Crg=="
}
```
