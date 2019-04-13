# Команда `close spend channel`
## Описание
Команда закрытия закрывает исходящий канал и переводит оставшиеся средства на исходящем адресе кошелька на адрес кошелька трекера.

Для того, чтобы адрес и канал были закрыты необходимо отсутствие открытых трансфер-каналов.
## Вход
Значение поля `cmd`: 64
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
	"message": "{\"address\":\"0x0cGNc_v0DLFQEMr0vG269bgZavAQbBrYdw-B9S8A==\",\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\\nnQIDAQAB\\n-----END PUBLIC KEY-----\\n\"}",
	"pk": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhCoNLv7e2FOQ488z4Qks\neCKIBYeWcz69IXBlEC0dy81Q+JTyp+zzYQomQ/e0jtzgkTk3lWey0DdB0wJXVe6c\nQSa1wvVQuv72Na+Xkknn19m7Kk6UkkDZ5PvkCAllaL8d02wMAO8tJKACTPeMVYP1\nOCQOdn5CoozrTEcVGhQ7ZBfmgIKJIpoiNYmebRY/63dwkeMUqgIjmiJDsKGI4aIF\nPF8cphRuIubHC8o/44/k+lREMje499aV9+q6BbXNnw0CcyjRdpDfA0SrlCSIJVYB\nWbzA951r0sJdutnTvwDjm64UT0nNEnYohcdS25JNQJhvIal6cOYBXdNZZ1hR1K/m\nnQIDAQAB\n-----END PUBLIC KEY-----\n",
	"sign": "T9LlRNgOwQ1VORwc6Q/HRmeV+1F7EGbq982+uKZ1Wpy50CU0Vanw/5Oiza7I1VJCAiNuxXVpn4wi0lBLj5iIvy3Ia27XqF55uszU1ohudCck7AR1LyZTBoT/IsHilXdNqeu+8122EovCdAoCokiO3WHKucWfvK1aj+Onfy9L9DYIJvJg6qEvoqgp8MHhETgcW1H7Zqx2QYSMLGtEMPJIMhK4nXLMZrs5pO8jiHbwzKD8jbQLp1Z5QiUJARrjR9/BAz7vX26/1I13lAH6mzWQ/+7y3vhlwSfzQnRr060mj+QEJd42dlu6Y2bU68ECNcmFcLqagDnkHVqqYBQ3G6tyVQ=="
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
      // Число байт, за которые списывались токены.
      // Type: number (uint64)
      "quantumVolume":0,
      // Фактическое число полученных байт
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
  "message":"{\"pk\":\"-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAn996Li8/jWBEdXRQyjXw\\nOAUdXTm/nLe3qKWBAitLHSeDKoh7vez2TOrFi3l96FChOIQI77YdNq3mYNjM1LNv\\nqXbHdP6IxTOIME5F1mxZqhVrsAjdDbSOJY6/gLE99GMsv6VuyLdf7PXbpMZrOLDV\\nj3p2nzhH0uaXQGWCuXYbxQkRp/qlj3bqxhU7EQpp94dRv06LpLUf6KAJzb0onBL8\\nGTkDpIwjg4b34L5989yY6eJqQUoysczSQPgHwvjZH5lJvTSE9LTENS6ZhhuBzVHh\\nw280UVBZLKt5vlwd6VnN+fR9n0G6XbeWU1gZ+moJLySNpHTgJdFZFnn8NlNFqRJC\\nSwIDAQAB\\n-----END PUBLIC KEY-----\\n\",\"state\":{\"amount\":\"10000\",\"quantumVolume\":0,\"volume\":0}}",
  "sign":"L4AIQgiQhL/+ilLdH9xsQT55mG3M0irKUICG0vq6+CkiXhGdX0HyzorT3T/3zlX9rvKdG8GCc313T4/c57w9d0cxbZpz3dv7MjoJn2U2ELkMd6dw5O64Uh1upQhY1tx69gu6p8V6FXjGLrnOzUnzZKklkkOAjcyiQVHQUxkfJ7YhPA5Oa964fF5U+0D5RuZElbXJDtGtpDke2p8a9OzlR8K0sgM0pYdJNKTj5ebB+DywTtU0A9yL6IRg8Jh4uXrk+5UEQud2Vm0/dkYk+XzbHcotcCheqOJIjaBIscG5anQrwxd84h4Oq/c3gw+gs90ab3laDfMzKlBDuKmeKz6whw=="
}
```
