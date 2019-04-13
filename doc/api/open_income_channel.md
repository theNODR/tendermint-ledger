# Команда `open income channel`
## Описание
Команда возвращает новый доступный для использования в качестве входящего адреса в `transfer`-каналах адрес кошелька.
## Вход
Значение поля `cmd`: 2
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
  "message": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n",
  "pk": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\nfQIDAQAB\n-----END PUBLIC KEY-----\n",
  "sign": "a2ZtyYhCJHhW68C8+PGtF1PlAPjqaLdNQIw9xSlzGvpRaXAehXtvGmAvu4m2d/d/sCbSVqOoJYzXmBVu5XflsZKi4HIo1gc0LUIWXKpqwAxktfc1QBDeviquwqHd+HQGfsO3hghegMM02zpHZ8OAdLm9vEYxetTCFPVbcUhVdQwWhYev7GgUxZFSZI2CnQwQVr4S/MORHfr8H36k7Lu96rrTq1Th+tDDAR/UQdBamrt81/t/s/CE2jJvYJUg26PfyUlnOGV7/xU7FwUbvjAXmDCJLsBIaAvlmXX0nwsSwdoRExA0E5W/GGVq8EoE/FvI2vJOa469nBBrENxPq901BA=="
}
```
## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Income address
    // Type: string
    "address": "",
    // Текущий состояние адреса
    // Type: Object
    "current": {
      // Число заработанных токенов
      // Type: string (uint64)
      "amount": "",
      // Число байт за которые, заработаны токены
      // Type: Number (uint64
      "quantumVolume": 0,
      // Число отданных байт
      // Всегда не более, чем значение quantumVolume
      // Type: Number (uint64)
      "volume": 0
    },
    // Минимальная цена за еденицу траффика, которую можно дать в рамках открываемого канала.
    // Type: Object
    "price": {
      // Число токенов
      // Type: string (UInt64)
      "amount": "",
      // до 2 ** quantumPower байт траффика необходимо отдать за токены
      // Type: Number (uint)
      "quantumPower": 0
    },
    // Фактический лимит средств адреса. Учитывает также средства, которые пир может заработать в открытых траффик-каналах.
    // Для раздающей ноды данная информация является справочной, и никак не используется.
    // Type: Object
    "limit": {
      // Токены, которые пир может заработать, с учетом в траффик-каналов
      // Type: String (uint64)
      "amount": "",
      // Объем траффика, который пир может отдать, если все открытые траффик-каналы будут выполнены полностью.
      // Type: Number (uint64)
      "quantumVolume": 0
    },
    // Tracker public key
    // Type: string
    // Value: key = new NodeRSA(message)
    "pk": "",
    // Время (в миллисикундах UTC), после которого адресс будет автоматически закрыт
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
  "message":"{\"address\":\"0x0dgKhkmT3064jXHRuAfm1YFW09DeFidrPEsLOe2A==\",\"current\":{\"amount\":\"0\",\"quantumVolume\":0,\"volume\":0},\"limit\":{\"amount\":\"0\",\"quantumVolume\":0},\"price\":{\"amount\":\"500\",\"quantumPower\":25},\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtKP/V1c3VvcxoCk7LZbP\\nclf/y6x+qwtDjt1YbQulP6YrcDI6j0UCYhpYThn8SjohER5TmrgffroUJ0dfQqHo\\nMrXcrBKjxfuFBdypFX4wKcDBRgvy1lDE+jow7GCg3A0CWDyCqLBI0nt1TbRkb9OM\\nc3Ie0d/YArQHVmGMKebqdeMMwEASUBhCknAmHIJOQApwm4L2qiKPAnUUHEKMUQMK\\nuGWVkF01R+l9qalq2nRUpRfeXKj4S6vSo7VPWFDXERkXFGQc1Zvu6XFH6zzOuZgs\\nuVYAaP65Q35w06/j4PIP37l0ObaROJUeG56MOjTFsxidryO2R0Iw8+XCloio5Lu2\\noQIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"timelock\":1532519081134,\"lifetime\":1799996}",
  "sign":"sgK+28lTIeWm+2SJov5Q+yVnfnDHIo99Ymq6K9IvBjqSSwfhqMYAGtAKR7D0xR1tlfRoyfLtTLKpUA3e6Khf5MrNNR6zNEcyv0MJzbDvlJudV3eRoiGxIt0HHuKSZq9txJotiAAGz9+2/uSoU/JYV83P+OUrHLsjc1kOa7HjWxLeEgdFX8/LdXRfXml6swe0XwcedPf/oYAhs6uhJe2V9uj+y3Vxh5leVYbA8uPKWj1piLTpahQJeA+D2NydjoZ0HIL+SnRUeFzbWlBL1LE+QQWbN+xT6flqQMMCsaFH1Sarp0m+qrBfZC9iUxVar0OJ/6FBwNftr1GC1C1RWr05pQ=="
}
```
