# Команда `close transfer channel`
## Описание
Команда закрытия трансфер канала между раздающей и принимающей нодами.

При закрытии трансфер канала происходит разблокировка средств траффик канала на адресе раздающей ноды и перевод средств, 
по которым была произведена фактическая согласованная работа между нодами.

Фактически переведенный объем средств не может быть больше объема средств, которые резервировались при открытии трансфер-канала.

С точки зрения сервера получение дублирующих сообщений закрытия трансфер-канала является нормальным, не приводит к ошибке, но и не вызывает выполнения повторных действий, включая запись транзакций в журнал.
## Вход
Значение поля `cmd`: 16
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Spend peer public key
    // Type: string
    // Value: key.exportKey('pkcs8-public-pem')
    "fromPK": "",
    // Income peer public key
    // Type: string
    // Value: key.exportKey('pkcs8-public-pem')
    "toPK": "",
    // Channel id. Это значение обязательно должно быть уникальным
    // Type: string
    // Value: may be SHA256(from + to + timestamp + plannedTokens + plannerVolume + price + cryptomatch)
    "channelId": "",
    // Spend address. Адрес, с которого будут списаны токены.
    // Всегда адрес, открытый с помощью метода open spend channel
    // Type: string
    "fromAddress": "",
    // Income address. Адрес, на который будут начислены токены.
    // Это может быть либо адрес, открытый при помощи метода open income channel, либо адрес собсвенного кошелька пира.
    // Type: string
    "toAddress": "",
    // Стоимость еденицы траффика
    // Type: String (uint64)
    "priceTokens": "",
    // Величина 2 ** priceQuantumPower определяет размер еденицы траффика для передачи.
    // Это означает, что передача 1,5 * 2 ** priceQuantumPower едениц траффика стоит как передача 2 едениц траффике.
    // Type: Number (uint8)
    "priceQuantumPower": 0,
    // Максимальное число едениц траффика для передачи в открытом трансфер канале
    // Type: Number (uint8)
    "plannedQuantumCount": 0,
    // Фактически оплаченное число едениц траффика
    // Всегда на более plannedQuantumCount
    // Type: Number (uint8)
    "quantumCount": 0,
    // Фактическое число передаваемых токенов
    // Всегда равно priceTokens * quantumCount
    // Type: string (Uint64)
    "tokens": "",
    // Фактически переданный объем траффика, байт
    // Type: Number (uint64)
    "volume": 0
  },
  // Spend peer public key
  // Type: string
  // Value: key.exportKey('pkcs8-public-pem')
  "fromPK": "",
  // RSA signature of message field by spend peer in base64 encoding
  // Type: string
  // Value: key.sign(message, 'base64', 'utf8')
  "fromSign": "",
  // Income peer public key
  // Type: string
  // Value: key.exportKey('pkcs8-public-pem')
  "toPK": "",
  // RSA signature of message field by income peer in base64 encoding
  // Type: string
  // Value: key.sign(message, 'base64', 'utf8')
  "toSign": ""
}
```
## Пример входа
```json
{
  "message":"{\"fromPK\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\\nnQIDAQAB\\n-----END PUBLIC KEY-----\",\"toPK\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\",\"channelId\":\"test15307826733498129\",\"fromAddress\":\"0x0cLzRz0YCVTJPwI_YtA8wRs-pT2JJCOgQOV_xHyQ==\",\"toAddress\":\"0x0dCkromg5T-x4nUHssrvY7k9eFpepHTheDORJELQ==\",\"priceTokens\":\"1000\",\"priceQuantumPower\":20,\"plannedQuantumCount\":10,\"tokens\":\"2000\",\"quantumCount\":2,\"volume\":1572864}",
  "fromPK":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----",
  "fromSign":"Jt1Rcatr5oFmCkPFT8vfYDBtukfrma+D8ULUNXkEyC8R2e1QFthvN3a9P1N/7T6lK1wAWPh1Tw8CTQ1y3ONZ+IaRg3Wq5MR5DkjPmWHYaiINdfHMzrzMUcC+4ILpywAjMWfTB3dG6VWQFDdD/AcJ2VqJu1QwIy1NPxjVIkPdovhpuVVHTf9mIOV6HtaYo4l76zSV5nD/1JwJb6+WDXyBbTGspOu305pOTcKDSkWboMu0pGOfel84AfBIl83Z/yDiSOlh/c/euwL6x/8ao3PVkvZuszkAaBmNRhaTV19PI84ZgcVdo9luISREAu45iMKvrb8rhKWDTjBTHyAxona20Q==",
  "toPK":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----",
  "toSign":"i3fM6rty9wKoovKyRrx+JGsPrTxAc4Uo+aDJ9spPIBfFA9D+i8t1OWlkN5o8fjdpRV7wJ56kRsxAx7AXS5qZmchcwc2IMnxWj+VSZe08uPRE4RzyTpibQtPjISGVIqXJ17Ey0oOFnNpcUc2yfx52sxEwyNSD5qxCNQknxDJ3hgRdaNMVXcF+xPsMmeB7vHzZhf161yldfqnWO2CoRGl01lKL/fJMrcQN6kx02NbtQclUPboHasnKCzj++9HNwyq58LkAjuCLI6R0EhC7eHrFhUj6C03EqR1iOOZuHcxhmMZnuWddxl+v3loBa7f12yOgcTpM4hH8tvXXwLw6vJCMjg=="
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
  "message": "{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8vQVhLPYuLSPGUHLgboK\\n9zqvaBi97vBSofH1ukz6r3zF2MpOWrMr6NKCxvItu32aEmjX/Pnh1UAZ0n7C/sx8\\nyGk8cJ5f11y+nv2odcZHbNmUhVdxKZO6ashpznJtGhL+vOEdp4Oha74DsyxRz6SX\\nSblJlNBSrxFYb9sZ+GjrEk6FmNYhm3J8vjkVvXtOrwX7d2/kU5tiCzSYBckBwJ2b\\nMW3g0Zt332N8Kj5Nc9xIBMSDGQn1+4DvgOoetLn/9lrE9Hr6RXv3RfLgkmf3HQTo\\n1Iux1Af94D899yYdkPnpNGwNwgcytJ+sJXYYt6FZp8b6SRLpS2ZYhky92Hi/21RQ\\nqwIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "sign": "er0COvR1+BCcfTRq7BYz4OWTNCeCVIIjb8APSAkne260J+ABsLk4q1kfP0GXtq0Vv0b9VLGRp4Ykt7wdCTWHwXzZIe1buyo2iKwXAbpVYdGLMjub6XJwlGGSag2b0UGHzZLysKZpwkbUiNtbTMTwPaR8tRzdPVK048N+BBZqZEAy88Y8AcDwt/8Y40APu6L/d5A5GvUzcRau/GszG3ggH09mo0+/9TLZs1c9hVAhaXEiL6dNNEpgJObk7WxhC4wsqeD7AvlaRdOy/qicLS6qmw3GzJVdCRSj/Gj64hAHYOCpMkS1dwiDEFNzRGWCwHJY9n07jcgr8I6OFpQygJ7Crg=="
}
```
