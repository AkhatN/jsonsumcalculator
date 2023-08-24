# jsonsumcalculator

## Описание

Программа для вычисления суммы чисел в файле.

Дан json файл в виде массива объектов как показано:
[
 {
    "a": 1,
   "b": 3
 },
 {
   "a": 5,
   "b": -9
 },
 {
   "a": -2,
   "b": 4
 }
 ...
]

Количество объектов = 1,000,000

Программа считывает файл и вычисляет сумму всех чисел. Для этого она:
* Разбивает данные на блоки (подмассивы), например по 100 обьектов
* Вычисляет сумму каждого блока, параллельно
* Считает сумму из чисел полученных от результатов блоков
* Выводит общий результат в консоль
* Количество горутин, для параллельной обработки блоков, можно получить при запуске программы через аргумент

## Как запустить
    # Usage: 
    go run main.go -file <путь к файлу> -numgo <Количество горутин> -numblocks <Количество блоков>

    # Пример: 
    go run main.go -file numbers.json -numgo 2 -numblocks 5

    # Пример без указания количества горутин и блоков (количество горутин и блоков будет 1 по умолчанию): 
    go run main.go -file numbers.json 

    # Компиляция проекта: 
    go build *.go

    # Пример запуска скомплилированного проекта: 
    ./main -file numbers.json -numgo 2 -numblocks 5

    # help: 
    go run main.go --help
    go run main.go -h

