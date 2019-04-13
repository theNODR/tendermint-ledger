# Команда `handshake`
## Описание
Команда предназначена для получения публичного ключа трекера,
который в дальнейшем будет использоваться для подтверждения истинности полученных от трекера сообщений 
## Вход
Значение поля `cmd`: 1
```json
{
  // Peer public key
  // Type: string
  // Value: key.exportKey('pkcs8-public-pem')
  "message": "",
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
  "message": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAk9Eit8JCsn/EoMX1HJos\nZW6wItW15DVuyrgABRe4Lx7vTj5v6H8tlV+/Cfkcp8yO96ohWyxdIQN9J38McyIK\neNZLxHEV+8IBeA4D5vLSz0zjHhGHmv9nI5Sr7bShHGAvNpppbL/k5Gd+yek5uJrr\nldF41OJHs4z58oQvDKXI+6csHPhzuv1ReD8JiL/NgmVKVsHIDMlIurM4OrY7E1XE\ney4yMgOPQB8GmxbLGtGwzhfWzl0FWWokwP07SBTFwILhKRYzjbjENl9zQdlgYfhO\n0xPwjTbwoW6p/LCb4kDx55RkDIjI7JpuqUIVkfcTWco5DPUrVANIfNQeWRMgkShF\ntQIDAQAB\n-----END PUBLIC KEY-----",
  "pk": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAk9Eit8JCsn/EoMX1HJos\nZW6wItW15DVuyrgABRe4Lx7vTj5v6H8tlV+/Cfkcp8yO96ohWyxdIQN9J38McyIK\neNZLxHEV+8IBeA4D5vLSz0zjHhGHmv9nI5Sr7bShHGAvNpppbL/k5Gd+yek5uJrr\nldF41OJHs4z58oQvDKXI+6csHPhzuv1ReD8JiL/NgmVKVsHIDMlIurM4OrY7E1XE\ney4yMgOPQB8GmxbLGtGwzhfWzl0FWWokwP07SBTFwILhKRYzjbjENl9zQdlgYfhO\n0xPwjTbwoW6p/LCb4kDx55RkDIjI7JpuqUIVkfcTWco5DPUrVANIfNQeWRMgkShF\ntQIDAQAB\n-----END PUBLIC KEY-----",
  "sign": "SbHkvIYxciVy/SaH9U+q36W4xDFG380TSOLWCzeBAELemr1zarOS1VL6OH85x6RJW8RZGqwM6sYA7y1h3ly9dGaz0BaA2KBJMcrUCCPZkT7EMSQROCX3pn1Okv2liAs9WSLRzCWXtv1wlYTecNZAT/dX42NssDy86vrbpEvP+CXgaZHh2QXlmO2epZ1Zoqk7xRzP18WGEfH4PWgLYXay+z1fo3jux5gJvHR7kEuDwMJw3Bm8KD/NV/0WFkrzhKxGHt1TRg56xQSpPsN2z9RJLtbcCiJdG4PICjAZIoBk2D3CW3P18xcXNSUoJh7CFNkRy5X7yFgkqmuumwz7ofNkAw=="
}
```
## Выход
```json
{
  // Tracker public key
  // Type: string
  // Value: key = new NodeRSA(message)
  "message": "",
  // RSA signature of message field in base64 encoding
  // Type: string
  // Value: key.verify(message, sign, 'utf8', 'base64')
  "sign": ""
}
```
## Пример выхода
```json
{
  "message": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8vQVhLPYuLSPGUHLgboK\n9zqvaBi97vBSofH1ukz6r3zF2MpOWrMr6NKCxvItu32aEmjX/Pnh1UAZ0n7C/sx8\nyGk8cJ5f11y+nv2odcZHbNmUhVdxKZO6ashpznJtGhL+vOEdp4Oha74DsyxRz6SX\nSblJlNBSrxFYb9sZ+GjrEk6FmNYhm3J8vjkVvXtOrwX7d2/kU5tiCzSYBckBwJ2b\nMW3g0Zt332N8Kj5Nc9xIBMSDGQn1+4DvgOoetLn/9lrE9Hr6RXv3RfLgkmf3HQTo\n1Iux1Af94D899yYdkPnpNGwNwgcytJ+sJXYYt6FZp8b6SRLpS2ZYhky92Hi/21RQ\nqwIDAQAB\n-----END PUBLIC KEY-----\n",
  "sign": "uzFmf7AIIHqB30ZR+YeCZfqmFqEO4mYatbIXwi0zh8KPFfpkZaf0sEpQCES+Kecn6qAAlQyH4er6faobtBO4QrZ0uIUzxDjn0S9ckDBluzvQR1+6R9SjR8VX9mjsW1eg+SvnkJDPyg5rU9kTUN3ZLGmfrD2uBaAXbg01UXrZQz4iK93lrZd4Q+nRwq81z2OXL549TjCyxUqU1midMjyS1fIKfwXOf3oDNmErEfTpzggd1L80oiX9+8CPrhYrxnf0Mekz/uaH1N+b0Y7HtPMitxtuadWlq8yphG+D3jeLJ5AbORYujhQtYzuJRqZZ+8NgQeOLa2L5qvMmeLbu7ZusZw=="
}
```
