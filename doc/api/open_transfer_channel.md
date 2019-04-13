# Команда `open transfer channel`
## Описание
Команда на открытие трансфер-канала между раздающей и принимающей нодами.

При открытии трансфер-канала происходит блокировка средств на адресе принимающей ноды.

С точки зрения сервера получение дублирующих сообщений открытия трансфер-канала является нормальным, не приводит к ошибке, но и не вызывает выполнения повторных действий, включая запись транзакций в журнал.
## Вход
Значение поля `cmd`: 8
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
    // Spend address. Адрес, с которого планируется списать токены.
    // Всегда адрес, открытый с помощью метода open spend channel
    // Type: string
    "fromAddress": "",
    // Income address. Адрес, на который планируется начислить токены.
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
    "plannedQuantumCount": 0
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
  "message":"{\"fromPK\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\\nnQIDAQAB\\n-----END PUBLICKEY-----\",\"toPK\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\",\"channelId\":\"test15307826733498129\",\"fromAddress\":\"0x0cLzRz0YCVTJPwI_YtA8wRs-pT2JJCOgQOV_xHyQ==\",\"toAddress\":\"0x0dCkromg5T-x4nUHssrvY7k9eFpepHTheDORJELQ==\",\"priceTokens\":\"1000\",\"priceQuantumPower\":20,\"plannedQuantumCount\":10}",
  "fromPK":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----",
  "fromSign":"GQxWxvUm9dah8wskUsUEodaaGOCghPU8/c4DD36wJtQMUXBrKs1+heIgmNDkPlhFqKZ6BV0PZXFZrgqGcrTyF24ATjLQJbVEwXyUwKebX0ddVoMVOyebUvYgrYAK1ng38ZFKdpyMRVTrRf9bh/UhPpGUMexTPzXsZK0a/sNluQt7b8pr23kwI8SjQeuxSdr+D6gdphkbJz4RqCyB2GMAQTHKl8xh5y9yOMM+B0kQ1JvVMA/UPXILrbeR3Qq/J/mgiUKTYNlgxoNoW6vKd9e3VqRoK7IigLvh2AbUQ1WXh8d0yNU/38essu+c8O0/w8PbW3jR9g02HHp0AVdP0mg58g==",
  "toPK":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----",
  "toSign":"QZnIYqSLsMSw6RL/YdPVSxfJtcCsvsntgT77QYqdPXucm84DQFYyGW4vqGwQNBnIk5lb2vWcb8q1wDeNTediWEiqBAWS6DYMXRCL5kSTu1u2d5LVZKK6K1JW6aMy1mytBMQDb3geiguofog+6i+Chwai+0hz/8Lj+ne8YTNzqGW16QApUywhVpHTQ3E8Nzo+b2hIaiE6lOQNvLrNUiQsCVGS1wS62QI9Rv1yl8HG/J2+h69PfH/dY4nnnRH8oLBUIUwY4NKeKpKNh9rNRWcq1W3AgWld7ChEx1uW5wmson3F76esK3Jj1UBgkl8UHzu4LTT61eaNH54ya8VM6tb56g=="
}
```
## <a name='open_transfer_channel-output'></a>Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Tracker public key
    // Type: string
    // Value: key = new NodeRSA(message)
    "pk": "",
    // Время (в миллисикундах UTC), после которого адрес будет автоматически закрыт
    // После этого момента никакие операции над адресом не доступны
    // Type: number (int64)
    "timeLock": 0,
    // Время (в миллисекундах, имеется некоторая погрешность в большую сторону), в течении которого канал будет доступен
    // Type: number (int64)
    "lifetime": 0
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
  "message": "{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtKP/V1c3VvcxoCk7LZbP\\nclf/y6x+qwtDjt1YbQulP6YrcDI6j0UCYhpYThn8SjohER5TmrgffroUJ0dfQqHo\\nMrXcrBKjxfuFBdypFX4wKcDBRgvy1lDE+jow7GCg3A0CWDyCqLBI0nt1TbRkb9OM\\nc3Ie0d/YArQHVmGMKebqdeMMwEASUBhCknAmHIJOQApwm4L2qiKPAnUUHEKMUQMK\\nuGWVkF01R+l9qalq2nRUpRfeXKj4S6vSo7VPWFDXERkXFGQc1Zvu6XFH6zzOuZgs\\nuVYAaP65Q35w06/j4PIP37l0ObaROJUeG56MOjTFsxidryO2R0Iw8+XCloio5Lu2\\noQIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"timelock\":1532517581249,\"lifetime\":299994}",
  "sign": "er0COvR1+BCcfTRq7BYz4OWTNCeCVIIjb8APSAkne260J+ABsLk4q1kfP0GXtq0Vv0b9VLGRp4Ykt7wdCTWHwXzZIe1buyo2iKwXAbpVYdGLMjub6XJwlGGSag2b0UGHzZLysKZpwkbUiNtbTMTwPaR8tRzdPVK048N+BBZqZEAy88Y8AcDwt/8Y40APu6L/d5A5GvUzcRau/GszG3ggH09mo0+/9TLZs1c9hVAhaXEiL6dNNEpgJObk7WxhC4wsqeD7AvlaRdOy/qicLS6qmw3GzJVdCRSj/Gj64hAHYOCpMkS1dwiDEFNzRGWCwHJY9n07jcgr8I6OFpQygJ7Crg=="
}
```
