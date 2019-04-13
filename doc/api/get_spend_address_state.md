# Команда `get spend address state`
## Описание
Команда возвращает текущий баланс открытого исходящего адреса кошелька.

Можно посмотреть баланс любого, не обязательно собственного, адреса.
## Вход
Значение поля `cmd`: 256
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
  "message":"{\"address\":\"0x0c0QGlyaLVC1T3CMuUnNtLyBdvXiFckELoTFIIgQ==\",\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\\nnQIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
  "sign":"BI2qeWSXppIsJOdbHa5JJluiGfItzLOmpLihUP6KQ8d88bmwBX12R1Xrvde+jiUGYiEN0FbDwLAjBR/zNjgIRbs/DxXSfLwIXADAuP8bf+W2mDC6RwRVLcmOl3ZUgkDxXwBn3xQFOkX9RpudJnDNBT5iiVGgO2HGlm6ckq1TBi6iANAHRXBX96zpq2q2DrzzP6Ve59KboRoA2YVtG7NRIOfBF1MT9v5rEpsCcneJU9kPlRvx45EkmTL7w/MIKgGxV/jaT5U/vNFYlfUj0ivTL2tu2tv6LyqWQcUnbeQvCDQW1TfXz9uuOusAtB9Cl79LZGO2EivRhoB0hUP9s9l8uA==",
  "pk":"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----\n"
}
```
## Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Spend address
    // Type: string
    "address": "",
    // Текущее состояние адреса
    // Type: Object
    "current": {
      // Число доступных токенов
      // Type: string (uint64.toString())
      "amount": "",
      // Число ОПЛАЧЕННЫХ полученных байт
      // Type: Number (uint8)
      "quantumVolume": 0,
      // Число полученных байт
      // Всегда не более, чем значение quantumVolume
      // Type: Number (uint64)
      "volume": 0 
    },
    // Максимальная цена за еденицу траффика, которую можно дать в рамках открываемого канала.
    // Type: Object
    "price": {
      // Число токенов
      // Type: string (uint64.toString())
      "amount": "",
      // до 2 ** quantumPower байт траффика можно получить за токены
      // Type: Number (uint8)
      "quantumPower": 0
    },
    // Фактический лимит средств адреса. Учитывает также средства, которые пир пообещал в открытых траффик-каналах.
    // Именно исходя из лимита определяется текущее предложение пира и возможность открытия новых траффик-каналов пиром.
    // Type: Object
    "limit": {
      // Токены, которые пир может пообещать в траффик-каналах
      // Type: string (uint64)
      "amount": 0.0,
      // Объем траффика, который пир может скачать, если все открытые траффик-каналы будут выполнены полностью.
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
  "message":"{\"address\":\"0x0crikSMJ8swCDY-cHHvjyNND_7xsbr8Gc67rw26w==\",\"current\":{\"amount\":\"1100\",\"quantumVolume\":0,\"volume\":0},\"limit\":{\"amount\":\"1100\",\"quantumVolume\":0},\"price\":{\"amount\":\"1100\",\"quantumPower\":20},\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3Ji3VGgaMM3Yv5MA007w\\nPW7HO1w9pV6EnuqJSCrCWbYrDN+9XxIZ6+EU/SmiiMTVyJvlAIyDPXGlePIFw0cO\\n8Dwvl9IGJUgc2kw4KbCEgvF9qj9VYyZYJ3D7iX57zwXg8A9XQwG1Zz/AWO4/dtcb\\nbg6h6BddK/ouPtGs4jOJNiaxVhy4d3IOK5ofP/ot/6HNaPGYfatpOEv5exW7o9ib\\nmzhOz1J0C/PSuPQv+8TSlmKa8d6u9S7v7zEHYBJBcvOH80KIBxghoR5hQYtoX7fW\\nwhyUv/+gNyGXBMn1EvSlNeJvQjFV393QA0p4wqF5wwZXQOdZJKRJoEdLxRFJZHmM\\nVwIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"timelock\":1530786557972}",
  "sign":"AEaK+Ubc0h3xCJwB2VgC5wF4hR42hJAzD7idgRORYovYApRTxPPWDvIdYzOp1vB0gP+MXAQRMHpLS1ZnJQfeDCqh3D0o10ZN5/ZPmYUeIKW0NFOu6+k3rx/hSAhN0atcnuAdws5DsRx3YvzC+dXJZcWLspUxOEpp37v3ZJDT2IuosPdjLsXdRLzdE33bvHDTYLetqy2SG7SfWoDc+su8QCOwd+cYqRpIcNeKzIbTJUYpubSpLH/4kjRh7I6dhh181/3INrqYKZvwcPayyD1CkfmwMShI925/xe1UB1MWxP2A3hgj2kmXd9D8r+7fUra8RXz4e/C5CkKMdrhktOId0g=="
}
```
