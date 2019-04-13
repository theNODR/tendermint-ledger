# Команда `get transfer channel state`
## Описание
Возвращает состояние трансфер канала по его идентификатору.
В общем случае, можно запросить состояние любого, не обязательно собственного трансфер канала.
Но необходимо иметь в виду, что даже полная нода может не знать состояние трансфер канала в случае, если и исходящий, и входящий адреса с ней не связаны.

Возможные состояния трансфер-канала - канала не существует, канал открыт, канал закрыт.
 
## Вход
Значение поля `cmd`: 512
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message":{
    // Transfer channel id
    // Type: string
    "channelId": "",
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
  "message":"{\"channelId\":\"test490411\",\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "pk":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n",
  "sign":"dapO35HByieNtgmelZ7Y/eJcactsEs/Y8aXQG/VTDlad2aYYsXnBF0tL2ImRIBkmtM0eL49IukAFmW7asw8u6ArjQN2BajbeuHZVZ+fc630U4+F9FsOfWvuK+ey8zwh6syDGknxEo2x5M4XO4Xm7BARY+T9woGHk2/NJaZr5U4HBkQu/cjYw0pw9AqsQj1NNoRRUplptJQK2OgweCSuA9t7l0rMtUcXj6+6yhSGaHFG5n/lyCeXixFx3QGNg5oEMtw5YZa4fXubksEUMIofIFDrkblv5ci25IUjJJcJRUXeZONUKuf047DPYYdBJmzkhU4/IF7InUMYeI0HqmkZUXw=="
}
```
## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Transfer channel id
    // Type: string
    "channelId": "",
    // Transfer channel state
    // 1 - transfer channel not exist
    // 2 - transfer channel is opened
    // 4 - transfer channel has closed already
    // Type: number (uint)
    "state": 0,
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
  "message": "{\"channelId\":\"test4904\",\"state\":4,\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8vQVhLPYuLSPGUHLgboK\\n9zqvaBi97vBSofH1ukz6r3zF2MpOWrMr6NKCxvItu32aEmjX/Pnh1UAZ0n7C/sx8\\nyGk8cJ5f11y+nv2odcZHbNmUhVdxKZO6ashpznJtGhL+vOEdp4Oha74DsyxRz6SX\\nSblJlNBSrxFYb9sZ+GjrEk6FmNYhm3J8vjkVvXtOrwX7d2/kU5tiCzSYBckBwJ2b\\nMW3g0Zt332N8Kj5Nc9xIBMSDGQn1+4DvgOoetLn/9lrE9Hr6RXv3RfLgkmf3HQTo\\n1Iux1Af94D899yYdkPnpNGwNwgcytJ+sJXYYt6FZp8b6SRLpS2ZYhky92Hi/21RQ\\nqwIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "sign":"YdulunCvwyB6GQqKThdgx6gl9ZkhtlxVzSVU7Pb3Ds3FW6B/5w9PVxLaBknx7RcCKGYiGJNisfq9jE6XD91jPsRe+EjYLNctccOE7xU7KzFYJhnTITpjCs1qV+NwCxN6hf/upqtC8zsUdhRpSacPLa+Q4xGaUWEtIZBYZzWxvwOI5NAp55jed8nHZ/W+V2ltyLq+HNv0Q8AHu7AYMxZaYYceP2/zsqWH7W9XHwoKOi6R78jkFh2OF2vrOWF7LYPyD3XiP96jy27bJVoXER3VJYzMNLh55xwYxbDMDkPBwSRHkRpm/5ScGY1xHx8GRZ08uvWP4r2R8EvwAYmAkcmPRQ=="
}
```
