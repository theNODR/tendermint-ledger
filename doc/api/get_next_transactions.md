# Команда `get next transactions`
## Описание
Возвращает страницу с транзакциями по параметрам, заданными в предшествующем запросе `get first transactions`.
Задается только номер страницы с данными для возврата. Все остальные параметры запроса фиксируются предществующим запросом.

## Вход
Значение поля `cmd`: 2048
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
    // Page to return
    // Type: number (uint)
    // Value: greater than 0
    "page": 1,
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
  "message":"{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"page\":5,\"cursorId\":\"a3a11ded-2560-4343-8c47-76be655ce2d5\"}",
  "pk":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n",
  "sign":"txiqG74/ACMHZT6H8mw9kVjvShhfDh8iYCEc6whSW1n/LBvEaJUhJZisQwSvktK6noph37psdZlohordF3GMm4LP3u6rjk2FkqFYor01GQm79VofK6L72c3USpV9hh1qai8ihtsJlpk2Wrm6DPWOVtakq+4owp7zMo388V+dTFhi+xMsEOMAgNBDws+NaPzUElB82i+CQwZlsXkZrl81iolqK96vBssROfDlvQiPy6ZykqWZb4UX/QNFT5Hn0BjpzG2d6bdh8Z+GLpf4dFDbewPLalz5rPkp3klSp5cUFz5XpZ06JPdDbOxQa5attgF8jRKyJGVWkVrZs8cUTQ4TqQ=="
}
```
## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Array of transactions
    // Type: Array of JSON
    // Value: fields values and types are equal of database fields. Please watch ./database.md
    "items":[{
      "amount":0,
      "contract_id":"",
      "from":"",
      "tx_id":"",
      "parent_tx_id":"",
      "planned_quantum_count":0,
      "price_amount":0,
      "price_quantum_power":0,
      "pk":"",
      "quantum_count":0,
      "sign":"",
      "status":0,
      "timelock":0,
      "timestamp":0,
      "to":"",
      "tx_type": 0,
      "version":"",
      "volume":0
    }]
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
  "message":"{\"items\":[{\"amount\":1100,\"contract_id\":\"\",\"from\":\"0x0crikSMJ8swCDY-cHHvjyNND_7xsbr8Gc67rw26w==\",\"tx_id\":\"8l1-y7iBmJ4TcJ_HVsvd9BSjnx1nx85DtWNZBg==\",\"parent_tx_id\":\"PKwYUDLh7aRR-S_CazNqGM1xAeg3WCRMpE0o7g==\",\"planned_quantum_count\":0,\"price_amount\":0,\"price_quantum_power\":0,\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3Ji3VGgaMM3Yv5MA007w\\nPW7HO1w9pV6EnuqJSCrCWbYrDN+9XxIZ6+EU/SmiiMTVyJvlAIyDPXGlePIFw0cO\\n8Dwvl9IGJUgc2kw4KbCEgvF9qj9VYyZYJ3D7iX57zwXg8A9XQwG1Zz/AWO4/dtcb\\nbg6h6BddK/ouPtGs4jOJNiaxVhy4d3IOK5ofP/ot/6HNaPGYfatpOEv5exW7o9ib\\nmzhOz1J0C/PSuPQv+8TSlmKa8d6u9S7v7zEHYBJBcvOH80KIBxghoR5hQYtoX7fW\\nwhyUv/+gNyGXBMn1EvSlNeJvQjFV393QA0p4wqF5wwZXQOdZJKRJoEdLxRFJZHmM\\nVwIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"quantum_count\":0,\"sign\":\"a\\ufffd3\\ufffd\\u001b\\u0007\\u00081}4\\u0013\\ufffd\\ufffd}\\ufffd\\ufffdI=\\ufffd\\u0017^M\\ufffd\\ufffd\\ufffd\\ufffd\\u001a\\ufffdj\\u001b\\t\\\\A\\t\\u0008،j\\u0004\\ufffd\\ufffdH\\ufffd\\ufffd+@\\ufffdJ\\ufffd)\\ufffd\\ufffd\\ufffdg\\ufffdD\\ufffd\\ufffd\\ufffd^{\\ufffd\\ufffd\\ufffd\\ufffdo~\\u0010Ѱ\\ufffd\\u001c\\u0002\\ufffdlƴ\\\\;L\\ufffd\\ufffdB\\ufffd\\ufffdx\\ufffd\\ufffd\\u003c\\ufffdCv\\ufffd\\ufffd/\\u0019\\ufffd\\ufffd\\u0002i\\ufffdf2\\ufffd\\ufffd\\ufffd헊T4T\\ufffd^\\ufffd\\ufffd\\ufffd\\u000c\\ufffdf\\ufffd\\ufffd'\\ufffd%\\\\Àp0\\ufffdb\\ufffdMw\\u003e\\u0000n\\u0011w)s\\ufffd\\ufffd\\ufffd\\u001c\\ufffd\\ufffd쥖އ\\u001a\\u0026a]9\\ufffd]\\ufffd\\rm}R\\ufffd\\u0010\\u0014M\\ufffda\\ufffd\\u0010\\ufffd*=\\ufffd\\ufffd\\ufffd\\ufffdk\\u0007\\ufffd\\ufffdu\\u000f|\\ufffd;\\ufffd\\ufffd\\u0018鏦z|\\ufffd\\ufffd\\ufffdGW˧\\ufffd\\ufffdMfo\\ufffd\\u001e\\ufffd`Cva\\u0016Y{\\ufffd\\ufffd\\ufffd\\ufffd\\ufffd}\\ufffd\\ufffd\\ufffd\\ufffd\\ufffd\\u0007\\ufffd\\u001b5+ތ\\ufffd\\n\\ufffdƲ\\ufffd\\t[E\\ufffd,\\ufffd\\ufffd\\ufffd$yÈ\",\"status\":1,\"timelock\":0,\"timestamp\":1530784758082,\"to\":\"0x0biC3La0MGgsmvseGLTopIVUE8ujs-k-PcvzQepg==\",\"tx_type\":32,\"version\":\"\",\"volume\":0}]}",
  "sign":"z/P6bsYREq2/t93x7PiDw4XttLnc5Xj3mNknQxkyqCsdFKHeM1jmuSR69G0v/KbEtH6BhUOCGDw8JgK7YYz99M8K4m1A9rzV3drgTgIZiPugeo3cijaNOtrRlQl0QACPjPYgEglE8wMsqoXMD5PjjYB5STBTsVDPk3/zmZVqoWhsJNaz7s19lxroS6w6R8CRUw0ChzOsUy0nBq79Ervpa591+80xFdSNPBgqZsgw7q95IuWBTxMXlKY4tALsuKzWnfHr/uWi1ZbmqX6+pOrhUOGd53rDAnZWCJhP7JOhxe0hPHJgsOY6RuQCyW8VC/geg21G+Ya4KGlIiMtnycAd6w=="
}
```
