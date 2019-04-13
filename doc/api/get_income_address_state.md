# Команда `get income address state`
## Описание
Команда возвращает текущий баланс открытого входящего адреса кошелька.

Можно посмотреть баланс любого, не обязательно собственного, адреса.
## Вход
Значение поля `cmd`: 128
```json
{
  // Request content
  // Type: string
  // Value: JSON.stringify(message-object)
  "message":{
    // Income address
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
  "message":"{\"address\":\"0x0dk3IICvD5O0ZNmrYazB7A0-RA4ApJKJuAb_8sfA==\",\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvJAZV5BQfIq7gJ5y+MXR\\nGu8rCsrYm5aVQ3zRebGITOYRnw2hpVCKGJcJqOHCBsOEFkWNBmAgyKdZoROGlpn4\\nHnfI8KMZ0bVLx0Ldx5ciPoXDxAv0unVvgPdXgjPsR8eUMbOg9ioND1ul/Gv2WoNo\\nSFbWnMZuJPtSGWcc/l23jLfzB6ziIohunN5ggRckeoitgFokoOQcudN2wlGsBwWc\\nlR3U0L7pAUKS1u2BM1GqpPAqFCzZhOJzWE/rKrnAm6RFJ2rNUCob7yxzzgVt5rh6\\niJlc9zMOI0JTGB4v/UBlvwWBjdp/NI9jlmzOaiuRvSQIczy7MA9I1+U6520/5EoG\\nfQIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "sign":"UEvNHxm4VpyRqjFndaPk1461miQyeBabDgLL74FWerPGYEesq24b8uPrVxKdkoTqDS7DOM3RIzLuytiAhrXqdK2Or0J13pfL9rLjfWrZu9BBnzC0Q/XnNkb/HXy0NWP6y1/6an9iiw3GMaCoiwxtyJ/UiuAv9yXmzuriEOgEzTt4DFDBazhvJEwxzmVkjf8K3XniBjTqmypgwnayZ9MSONFO5rg0lOE1ZgU068ci8r0cftZHS7qeVFj+ADccaUE9GE66wmSk1+uNWWIGuogdTssmBB3dLW2PIn1VG7kBsHcJbfyLo2IMykyv4uwUD3q0t6Eycfik6KwGDom3p5d2Dw==",
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
    "timeLock": 0
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
  "message":"{\"address\":\"0x0dE7ZvUJS7d71iR-xfDbqn6hYbxV6mqVnrviEVqQ==\",\"current\":{\"amount\":\"0\",\"quantumVolume\":0,\"volume\":0},\"limit\":{\"amount\":\"0\",\"quantumVolume\":0},\"price\":{\"amount\":\"900\",\"quantumPower\":20},\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3Ji3VGgaMM3Yv5MA007w\\nPW7HO1w9pV6EnuqJSCrCWbYrDN+9XxIZ6+EU/SmiiMTVyJvlAIyDPXGlePIFw0cO\\n8Dwvl9IGJUgc2kw4KbCEgvF9qj9VYyZYJ3D7iX57zwXg8A9XQwG1Zz/AWO4/dtcb\\nbg6h6BddK/ouPtGs4jOJNiaxVhy4d3IOK5ofP/ot/6HNaPGYfatpOEv5exW7o9ib\\nmzhOz1J0C/PSuPQv+8TSlmKa8d6u9S7v7zEHYBJBcvOH80KIBxghoR5hQYtoX7fW\\nwhyUv/+gNyGXBMn1EvSlNeJvQjFV393QA0p4wqF5wwZXQOdZJKRJoEdLxRFJZHmM\\nVwIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"timelock\":1530786720750}",
  "sign":"cg3fHialsgnnVOj75jgnCkP/jy0PjHvMPjm3rlJ/KwjULVVSCiDfY1uoYeCihA2DVM7hFFbM/U9gvayTFMPxrJfDskHGz82FhSBOrCAJaSC5A1y6WyTrk2iPkU9okoup9wpnh+Skxy4joU1ZZ06PZCS+bmrVn/KUMzqjRCVNFolKOI+xVjTwAtuhGA1Cj94/kwgPI8fA2/7ftQohT95TB871kDqTQAtRwfEABtKBfy6cVxYyHh7Ytj0tdTEW/9oQC2HbfK69L7urEXzrgf4Nd6K2UPY8hT3bWD1sC/lzyZ1QkOqPMKn+7xjXpEDrGqpg2MOQIOFqPyQq6wpQ5fv2Dw=="
}
```
