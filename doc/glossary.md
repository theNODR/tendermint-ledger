# Глоссарий

## Адрес
Адрес кошелька ноды.
Для обозначения обычно используется A<sub>обозначение ноды</sub>, например, A<sub>T</sub>.
Вычисляется как `A = sha256.sha224(PK)`

## Входящий адрес
Входящий адрес раздающей ноды на ноде трекера кошелька.
Для обозначения обычно используется <sup>I<sub>обозначение раздающей ноды</sub></sup>A<sub>обозначение ноды трекера</sub>, например, <sup>I<sub>B</sub></sup>A<sub>T</sub>.
Вычисляется как ```A = sha256.sha224(...)```.
Находится во владении у ноды трекера кошелька.
Может иметь явные и неявные ограничения, например, на использование в транзакциях, на время жизни, и т.п.

## Закрытие контракта
Подтверждение выполнения контракта или его части сначала принимающей, а затем раздающей нодой.

## Исходящий адрес
Исходящий адрес принимающей ноды на ноде трекера кошелька.
Для обозначения обычно используется <sup>S<sub>обозначение раздающей ноды</sub></sup>A<sub>обозначение ноды трекера</sub>, например, <sup>S<sub>C</sub></sup>A<sub>T</sub>.
Вычисляется как ```A = sha256.sha224(...)```.
Находится во владении у ноды трекера кошелька.
Может иметь явные и неявные ограничения, например, на использование в транзакциях, на время жизни и т.п.

## Контракт
Условия подтверждающие готовность:
* со стороны раздающей ноды предоставить данные видеопотока;
* со стороны принимающей ноды предоставить фиксированное число токенов за фиксированный размер переданных данных.
В рамках контракта может быть выполнено меньшее число работы (передачи трафика) за меньшее вознаграждение (число токенов), пропорционально заявленному в изначальном контракте.
При этом максимально возможное вознаграждение (число монет), полученное раздающей нодой не может быть больше указанного в контракте. 

## Нода

## Нода трекера кошелька
Для обозначения обычно используется T.

## Открытие контракта
Процесс подписи контракта сначала принимающей, а затем раздающей нодой. Если обе ноды подписали контракт, он считается открытым.

## Пиринговая сеть

## Приватный ключ ноды
Для обозначения обычно используется SK<sub>обозначение ноды</sub>, например, SK<sub>B</sub>. 

## Принимающая нода
Для обозначения обычно используется C.

## Публичный ключ ноды
Для обозначения обычно используется PK<sub>обозначение ноды</PK>, например, PK<sub>T</sub>

## Раздающая нода
Для обозначения обычно используется B.

## Сегмент
Атомарная порция данных видеопотока