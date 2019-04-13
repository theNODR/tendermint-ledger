# Команда `get first transactions`
## Описание
Возвращает список транзакций и дополнительную информацию по запросу по ним.
В дополнительную информацию входит общее число транзакций по запросу, общее число страниц, число транзакций на странице, номер возвращенной страницы и уникальный идентификатор запроса, с помощью которого можно запращивать другие страницы по этому же запросу.
Уникальный идентификатор запроса обладает следующими ключевыми моментами:
* фиксирует список транзакий. Это означает, что в последующие запросы по этому идентификатору не попадут более новые транзакции.
* время жизни идентификатора - 5 минут с момента его появления или последнего испольльзования, если таковые были.

Список транзакций можно фильтровать (не сделано), сортировать (не сделано), можно задать размер страницы и номер возвращаемой страницы.

## Вход
Значение поля `cmd`: 1024
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Row on page count and page number parameters
    // Type: JSON
    "limits": {
      // Row on page count
      // Type: number (uint)
      // Value: greater than 0
      "pageSize":1,
      // Returned page number. Starts with 1
      // Type: number (uint)
      // Value: greater than 0
      "page":1
    },
    // Peer public key
    // Type: string
    // Value: key.exportKey('pkcs8-public-pem')
    "pk": ""
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
  "message":"{\"limits\":{\"pageSize\":1,\"page\":1},\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "pk":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n",
  "sign":"MBv1egcSVapet1RuHA/k5FnQI7gEY7P0OBINO1+GPHAtAQCeGFnLrctpBRwrY+htYVInTeXDVM2lv/H3Edb16SLY2hbT3wMkaznORbeKEZaLyX7NFjQ3I6DNkRN2szv3+bqrGm3ptWUE/q1EkbjqNz3pB6QKXipUMFxe0znlDjL3imlEhELvuDCIho7SgWzeIXpLoWZnMB6luD63g9+/lOFstpwT+KhA4I+BXWIfg90K49j4Vz61y5cKNarUzJ/fe4khgCAgCS2dXfsiPluKAI8+HohwUQ9vWcw36KtT1LOgGl/hnw+ite6AlZTtG8tzHKQ+7VKlpvyOe1TgwXK+nA=="
}
```
## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Metadata info
    // Type: JSON
    "_meta": {
      // Cursor id to `get next transactions` method. Fix start point, order, limits and other query consts
      // Type: string
      // Value: uuid
      "cursorId": "",
      // Total count of transactions
      // Type: number (uint)
      "totalCount": 0,
      // Total count of pages
      // Type: number (uint)
      "pageCount": 0,
      // Current page
      // Type: number (uint)
      "currentPage": 0,
      // Transactions on page
      // Type: number (uint)
      "perPage": 0
    },
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
  "message":"{\"_meta\":{\"cursorId\":\"f654a892-98cd-48a2-8209-cbd2a0f77401\",\"totalCount\":5,\"pageCount\":5,\"currentPage\":1,\"perPage\":1},\"items\":[{\"amount\":0,\"contract_id\":\"\",\"from\":\"0x0dE7ZvUJS7d71iR-xfDbqn6hYbxV6mqVnrviEVqQ==\",\"tx_id\":\"rUMyljclu4Rn-dattVnfNfRxZzfPE-GAJs6eDg==\",\"parent_tx_id\":\"W3aPv605nANJ4YyszCBqt_MnTQsNb5EcQ_3qNQ==\",\"planned_quantum_count\":0,\"price_amount\":0,\"price_quantum_power\":0,\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3Ji3VGgaMM3Yv5MA007w\\nPW7HO1w9pV6EnuqJSCrCWbYrDN+9XxIZ6+EU/SmiiMTVyJvlAIyDPXGlePIFw0cO\\n8Dwvl9IGJUgc2kw4KbCEgvF9qj9VYyZYJ3D7iX57zwXg8A9XQwG1Zz/AWO4/dtcb\\nbg6h6BddK/ouPtGs4jOJNiaxVhy4d3IOK5ofP/ot/6HNaPGYfatpOEv5exW7o9ib\\nmzhOz1J0C/PSuPQv+8TSlmKa8d6u9S7v7zEHYBJBcvOH80KIBxghoR5hQYtoX7fW\\nwhyUv/+gNyGXBMn1EvSlNeJvQjFV393QA0p4wqF5wwZXQOdZJKRJoEdLxRFJZHmM\\nVwIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"quantum_count\":0,\"sign\":\"\\u0015/\\u0017\\ufffdc\\ufffd5_\\ufffdw\\ufffd\\ufffd_\\u0012$6?\\ufffd\\ufffdXs\\u0010\\rY\\ufffdN\\ufffd\\ufffd\\ufffd\\u0011\\ufffd\\u0010N\\ufffdn\\ufffd?\\ufffdǼ\\ufffd*\\ufffd\\ufffd\\ufffd\\ufffd\\u0006\\u0011vR\\ufffdr\\ufffdѩ-\\\"\\ufffdQ;\\ufffdχǧ\\ufffd\\u0014n\\ufffd-\\ufffd\\ufffd\\ufffd\\ufffd%E\\ufffd\\ufffd+u\\ufffd\\ufffd\\ufffdw\\ufffd\\t\\ufffdc\\ufffd\\u0007\\ufffd2\\ufffd\\u001b\\ufffd\\u0004@\\ufffdF\\ufffd\\u0018\\ufffd\\u000e\\ufffd\\ufffd-\\ufffd\\u000b\\u003e\\ufffd\\ufffd|4\\ufffd\\n]\\u001e\\ufffd\\ufffdd\\u0019\\ufffd\\u003eI\\ufffd\\ufffdF\\ufffd\\ufffd\\ufffdQ_\\u001e\\ufffd\\ufffd\\ufffd\\ufffd\\u0007\\ufffd͙\\ufffd\\ufffd\\ufffd0\\ufffd\\ufffdBZ`\\ufffd\\ufffd\\u0015$Y\\\"*fLԄw\\ufffd\\ufffdy\\ufffd:\\u0019\\ufffde\\ufffdʜ՞Hdޑ\\\\\\ufffd\\ufffdv\\ufffdxXn\\u0014H\\ufffd\\u00084\\ufffd\\ufffd\\u00058q\\ufffdw\\ufffd\\ufffd\\u000c\\ufffd\\ufffd=\\ufffd\\ufffd[^!\\ufffd\\ufffd\\u001c\\u0007\\ufffd\\ufffd\\ufffdl\\ufffd\\ufffd\\ufffdYX\\u003c6\\ufffd\\u0018\\ufffd]\\ufffdJ\\ufffd5\\ufffd\\u0012^$\\u0002\\ufffd\\ufffd\\u001d\\ufffd\\ufffdB\\ufffd\\ufffd\\ufffdRF\\ufffd\\ufffdֽ\",\"status\":1,\"timelock\":0,\"timestamp\":1530784920867,\"to\":\"0x0biC3La0MGgsmvseGLTopIVUE8ujs-k-PcvzQepg==\",\"tx_type\":64,\"version\":\"\",\"volume\":0}]}",
  "sign":"kf35WhGO5kxRm2NaokeGcKhb0/xOSo7yYdMqlR1lqkQXhESR/90vzAzsisMutqzwgXB2XDal/8P4cVYD6lcm3NyWfNKgWZYwuWOpIFPShzkrH22qy7BpDDufYjhQvKjveakV0dXZD7wskOhrAQFhlv69mlA+IlZweF/2Lua1E1RcHhZfx/rOXUqHNAHHnhbqaQXNeZc658hUVLkxHgqlThy9UVtMtzueXcV1UUJUL0j6dbRSpCxIirTgi6UzXAQHsA+DjV+zt6zD0oB8yaCZsl5BldK7cJVlhICevn3q01MZkumEg1qnAAgtkup8FM4VKXXiFPuDJn2mjCirmbjKKw=="
}
```
