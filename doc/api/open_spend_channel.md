# Команда `open spend channel`
## Описание
Команда возвращает *новый* доступный для использования в качестве исходящего адреса в `transfer`-каналах адрес кошелька.

## <a name='open_spend_channel-input'></a>Вход
Значение поля `cmd`: 4
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
  "message": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----\n",
  "pk": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----\n",
  "sign": "HYNV2vQl8n1M9vueTxmZI6ofb7iuTEMax2as201MdodXV5MoARflsNTduGy3cTVNpjdAZrvqOlp5HZausg0m6f60wN+Y9SlyG/Fkl4quxemTCV13qu0LTsTVnHBSxcwmfwAbdXtU5n83p2YhqDyEGHGjHhCExRT/YjjoCF9yKyG/6UqlVqzDHkwy46cDl8sDrQIyQuDXybp2tkS8V+03M70HXALf9rPJXFkruGcx5N+ZrA11GDuZ8Zv1iLDF44fdb0+fCMqOR2/aMB2K0StLIbK0rLNFfi5D/i/QrmHyHjivwqTYCZkSKX47wXs+X8e3tBs2rZR1S8NX7JqxjTsGfg=="
}
```
## <a name='open_spend_channel-output'></a>Выход
```json
{
  // Server response
  // Type: string
  // Value: JSON.stringify(message-object)
  "message": {
    // Spend address
    // Type: string
    "address": "",
    // Текущий состояние адреса
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
  "message":"{\"address\":\"0x0cDQK1qBrUwD9iidPjueduzwpwIWXJJY3jJ0OVpw==\",\"current\":{\"amount\":\"10000\",\"quantumVolume\":0,\"volume\":0},\"limit\":{\"amount\":\"10000\",\"quantumVolume\":0},\"price\":{\"amount\":\"500\",\"quantumPower\":25},\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtKP/V1c3VvcxoCk7LZbP\\nclf/y6x+qwtDjt1YbQulP6YrcDI6j0UCYhpYThn8SjohER5TmrgffroUJ0dfQqHo\\nMrXcrBKjxfuFBdypFX4wKcDBRgvy1lDE+jow7GCg3A0CWDyCqLBI0nt1TbRkb9OM\\nc3Ie0d/YArQHVmGMKebqdeMMwEASUBhCknAmHIJOQApwm4L2qiKPAnUUHEKMUQMK\\nuGWVkF01R+l9qalq2nRUpRfeXKj4S6vSo7VPWFDXERkXFGQc1Zvu6XFH6zzOuZgs\\nuVYAaP65Q35w06/j4PIP37l0ObaROJUeG56MOjTFsxidryO2R0Iw8+XCloio5Lu2\\noQIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"timelock\":1532519081139,\"lifetime\":1799996}",
  "sign":"LC/4WAGCi4ZACAHRPmyMStH6ijt/nDqu3LV4YzonRyqq5bPo6M9TEvX93SXyeLHTUJjFzb76E1yfx6Wf/u7ngZb/0kj1NFVgWXN4RJBbATzm/bJOELGDVhh+MUJwnlJoEtVqUk8DxlqSHbWg/U//AFa8Tfapbtj58PCTfq7595aeeId7yjE6+jOACblaRUkAPNfCym7r2MwqD101LWXNHe06iLeufZ3LcaXvrBYIJvgzcqObIfGpz03Ih4HXNtT/yFp8yaBSBN4QHmU/fTkGDzSSEUPPMdD0Huv7Dm3MCIkJQ9qUoNuNdiBmw/nFPOIqv/YO2KQSNoy3Cdbx4rvE0w=="
}
```
