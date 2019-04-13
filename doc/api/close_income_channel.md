# Команда `close income channel`
## Описание
Команда закрытия закрывает входящий канал и переводит оставшиеся на входящем адресе кошелька средства на адрес кошелька трекера.

Для того, чтобы адрес и канал были закрыты необходимо отсутствие открытых трансфер-каналов.
## Вход
Значение поля `cmd`: 32
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message":{
    // Spend address
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
	"message": "{\"address\":\"0x0d6j_nIp8gSf7xJ4IFe_igMqYcK99vSi-oL7GMcQ==\",\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
	"pk": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n",
	"sign": "GRGcAYTdw+r5sGHt3LyelMleaZiF+YyjO98LT8h+Qeq/zHkUIawQOYkbH7kTGVbW3rqUetG3JIZCF8EwspoBfz5AhJNvGaYyB6sFvaHjh9Fn2hLLk8oR9kdVsX1mp782shM7gLR+n8/w2rf9RC6ne86tduxUHWcWY6BFPePQ4N7fhRK/BoPZLf49IUKfZjr1o1dcgCP5dy5VBLdyEHqd/TKxvV3KTzsdmviAMVyTOCo9uYn2U2tZre5LUg5jyW8BPktJ0pFUMfIptuIKKAZXqlG+yhcFSCn/w9zOVJrdPQusZ83TRQb0H+EAJRKi1HvDEcqRZWmWbCR7hkXb+x7k6g=="
}
```
## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Address state before close
    // Type: Object
    "state": {
      // Количество токенов на адресе на момент закрытия канала
      // Type: string
      "amount":"",
      // Число байт, за которые начисленялись токены.
      // Type: number (uint64)
      "quantumVolume":0,
      // Фактическое число переданных байт
      // Type: number (uint64)
      "volume":0
    },
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
  "message":"{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAn996Li8/jWBEdXRQyjXw\\nOAUdXTm/nLe3qKWBAitLHSeDKoh7vez2TOrFi3l96FChOIQI77YdNq3mYNjM1LNv\\nqXbHdP6IxTOIME5F1mxZqhVrsAjdDbSOJY6/gLE99GMsv6VuyLdf7PXbpMZrOLDV\\nj3p2nzhH0uaXQGWCuXYbxQkRp/qlj3bqxhU7EQpp94dRv06LpLUf6KAJzb0onBL8\\nGTkDpIwjg4b34L5989yY6eJqQUoysczSQPgHwvjZH5lJvTSE9LTENS6ZhhuBzVHh\\nw280UVBZLKt5vlwd6VnN+fR9n0G6XbeWU1gZ+moJLySNpHTgJdFZFnn8NlNFqRJC\\nSwIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"state\":{\"amount\":\"0\",\"quantumVolume\":0,\"volume\":0}}",
  "sign":"RC6mFfFa3QIv2QCHi9Sa10qwbIAOKqUYMkjo0aVlY8FT48AmgATXWnRf9C3BRhD8SfV/M4LgTwqWHBYGT3DPFaxTSJpYKjnNP6HF8Am3j6B4uhRTgFKd5A+93BBCvd+U6Wd1iIVIIZ+n8zt3MxE6hC4UJMQdo3Ibh13rAz35aUHZ6VvzkbxFHl9CX8NT4W5UR4dFGFWhBT9Ex5ofxHbyHhH57wYFQarW5AnSRjrlI1UCf+GezNARodJlHQdFbVcgO+mnwWqwmL/RzIBYL1FexK8Nq3qJ7PjDObat9OlRwWySwR5sU/evzjptYptC/V/9vUfsU7NiciUuloyod6I07w=="
}
```
