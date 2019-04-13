# Краткое содержание
1. Стоимость задается строкой, которое представляет собой число с 4 знаками после запятой.
    1. На бакенде go используется готовая реализация типа decimal.
    2. В хранилище пишется строка.
    3. На js и c++ подозреваю тоже есть встроенные/готовые реализации типа decimal.
2. При задании цены траффика размер еденицы объема траффика (квант) задается как двоичный логарифм от числа байт. Это число должно быть целым.
3. При подсчете вознаграждения учитывается целое число квантов траффика (с округлением в большую сторону).

# Полное содежание
## Способ записи телепортов

Телепорт является единственной рассчётной еденицей системы.

Телепорт состоит из атомов.
Каждый телепорт состоит из **10000** атомов.
Атом являнется неделимой еденицей.
> То есть нельзя каким-либо образом передать между кошельками нецелое число атомов, например, `0.5`.

При взаимодействии с API леджера количество телепортов задается строкой, которая должна удовлетворять регулярному выражению:
```javascript
/^((\.[0-9]{1,4})|([0-9]+(\.[0-9]{0,4})))?$/g
```
В случае, если переданная строка не удовлетворяет указанному регулярному выражению, строка не задает число телепортов и является ошибочной.
Отсутствие точки в переданной строке означает, что задано целое число телепортов.
В случае, если переданная строка содержит точку:
* Блок цифр слева от точки означает число телепортов:
    * Пустой блок цифр означает, что число телепортов в записи 0.
> В записи телепортов следует обратить внимание на:
> * Префиксные нули *не являются* значимыми и *могут* быть опущены при записи, то есть запись `00100.1` **эквивалентна** записи `100.1` и означает 100 телепортов и 1000 атомов. 
> * Суффиксные нули *являются* значимыми и *не могут* быть опущены, то есть запись `100.1` **не эквивалентна** записи `1.1`.
* Блок цифр справа от точки означает число атомов.
    * Пустой блок цифр означает, что число атомов.
    * Максимальное число атомов в записи - 9999. Запись 10001 атома соответствует записи 1 телепорт и 1 атом.
> В записи атомов следует обратить внимание на:
> * Префиксные нули *являются* значимыми и *не могут* быть опущены при записи, то есть запись `123.1` **не эквивалентна** записи `123.001`.
> * Суффиксные нули  *не являются значимыми* и *могут* быть опущены при записи, то есть запись `123.100` **эквивалентна** записи `123.1` и означает 123 телепорта и 1000 атомов.

## Способ записи траффика
### Мощность кванта траффика для оплаты

Размер траффика для оплаты должен является степенью двойки байт.
То есть для рассчетов могут использоваться величины в `1`, `2`, `128` байт.

Назовем такой размер траффика квантом (порцией, ...).
Квант траффика характеризуется мощностью (от английского power). Мощность кванта записывается как степень двойки. То есть:
* для 1Б = 2<sup>0</sup> записывается как квант мощности `P = 0`,
* для 128Б = 2<sup>7</sup> записывается как квант мощности `P = 7`,
* для 1КБ = 2<sup>10</sup> записывается как квант мощности `P = 10`,
* для 1МБ = 2<sup>20</sup> записывается как квант мощности `P = 20`,
* для 1ГБ = 2<sup>30</sup> записывается как квант мощности `P = 30`.  

Можность кванта - это целое беззнаковое число.

### Для фактического учета
Целое беззнаковое число, представляющее число переданных байт.

## Определение цены траффика

Цена траффика опредеяется парой значений:
1. `P` - квантом.
2. `A` - числом телепортов, которым вознаграждается передача кванта.

Для обозначения цены траффика используется символ `C`.

Например, цена может определяться как 1 телепорт 100 атомов за квант мощности 3. Иначе, 1 телепорт 100 атомов за каждые 8 байт траффика.
Формальная запись:
1. `P = 3` (квант мощности 3, 8 = 2<sup>3</sup>).
2. `A = "1.01"` (1 телепорт 100 атомов).

Две цены траффика С<sub>1</sub> и C<sub>2</sub> можно сравнивать толко в том случае, если они обе обладают одинаковой мощностью квантов.
В этом случае сравнение осуществляется путем сравнения вознаграждения за передачу кванта.

Значение мощности кванта для цены *можно* повышать.
В этом случае при увеличении на каждую еденицу мощности кванта вознагражение за передачу кванта увеличивается в 2 раза.

Значение мощности кванта для цены *нельзя* понижать.

## Как считается траффик

Перед передачей траффика открывается трансфер-канал. В нём фиксируется цена траффика C.
В трансфер-канале можность кванта является неделима. Т.е. траффик считается в целых квантах. 

> Никак нельзя передать 1.5 кванта траффика. В данном случае необходимо считать 2 кванта траффика.

> Заметим, что в этом случае заданная цена C<sub>1</sub> (A = "1.01", P = 3) в общем случае **не экивалентно** заданной цене C<sub>2</sub> (A = "2.02", P = 4)

Число переданных квантов получается из фактического размера траффика в байтах и из мощности кванта.
Определим число переданных квантое как наименьшее целое число, которое при умножении на число байт в кванте превысит число переданных байт. 

Положим, что цена траффика в трансфер согласована в 1 телепорт 100 атомов за 8 байт:
1. `P = 3` (т.к. 8 = 2<sup>3</sup>).
2. `A = "1.01"`.

Размер траффика 36 байта. Число переданых квантов - `Q = 5`: 
```
4 * 8 < 36 <= 5 * 8
```
.
Стоимость 6 квантов составляет 1.01 * 5 = 5.05 (5 телепортов и 500 атомов).

Итого:
1. Перечисленное вознаграждение `S = "5.05"`.
2. Переданное число едениц траффика `Q = 5`.
3. Переданно байт `V = 44`.
4. Стоимость кванта `A = "1.01"`.
5. Мощность кванта `P = 3`.

> Если `C (A = "2.02", P = 4)`, то для размера траффика `V = 36` байт `Q = 3`, `S="6.06"` 