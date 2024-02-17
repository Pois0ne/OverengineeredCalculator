# Быстрое начало

У данной программы **нет фронта**, только http-api. Использование подразумевается через curl.


Для запуска программы достаточно использовать одну эту команду:
```sh
go run main.go
```

_Так же обратите внимание на логи! Программа их активно пишет :)_

Теперь вы можете поиграться с готовыми примерам
# Примеры curl запросов:

1) Отправка выражения "2 + 2 * 2" на сервер
>⚠️ Из-за особенностей http нельзя передавать `+` в запросе. 
> 
> Надо использовать `%2B` или `p`.
> 
>`2 %2B 2` или `2 p 2` = `2 + 2`
``` sh
curl -X POST -d "expression=2 %2B 2 * 2" http://localhost:8080/expression
```
>Переданный в запросе пример будет взят в обработку. Это можно увидеть в логах
#

2) Получение списка выражений и их состояний:
``` sh
curl http://localhost:8080/expressions
```
> После использования данной команды нам возвращается json с данными о состоянии ранее введённых примеров. К примеру: `[{"id":"5713","content":"2 + 2 * 2","status":"completed","created":"2024-02-17T01:37:36.736146+03:00","updated":"2024-02-17T01:37:46.969893+03:00","result":6,"processed":true}]`
> 
> А теперь разберём понятнее:
> 
> `"id":"5713"` - ID генерируемый случайным образом
> 
> `"content":"2 + 2 * 2"` - Строка с примером которую мы передали
> 
> `"status":"completed"` - статус решения. Меняется в процессе решения. Может символизировать о том что: решение в процессе, решение невозможно, решено
> 
> `"created":"2024-02-17T01:37:36.736146+03:00"` - Дата создания
> 
> `updated":"2024-02-17T01:37:46.969893+03:00"` - Дата обновления (решения)
> 
> `"result":6` - результат решения
> 
> `"processed":true` - было в процессе решения?


#
3) Изменение времени на решение примера в секундах:
``` sh
curl -X POST -d "TimeToProceed=30" http://localhost:8080/timer
```
> Изменение времени решения наших "очень" ресурсоёмких задач
#

4) Добавить 5 вычислительных "демонов":
``` sh
curl -X POST -d "add=5" http://localhost:8080/computationalAgent
```
>С помощью этой команды можно добавлять "демонов" (горутины) которые решают примеры
#

5) Показать список вычислительных "демонов" и их состояние:
``` sh
curl http://localhost:8080/agentsList
```
>На выходе мы опять же получим json с такими данными: `[{"status":"working","num":1},{"status":"working","num":2},{"status":"idle","num":3},{"status":"idle","num":4},{"status":"working","num":5}]
`
> 
> А теперь вновь понятным языком:
> 
> Нам выводится массив "демонов". Каждый отдельный демон имеет свой номер и состояние. Номер думаю понятен, а состояния может быть два: `"working"` и `"idle"`
>(Состояние работы, когда происходят вычисления и состояние покоя, когда демон простаивает)
#

# Дополнительно
Для парсинга используется сторонняя библиотека
[govaluate](https://github.com/Knetic/govaluate).

Изначально код данной программы был заточен на работу с GUI, но в итоге я не успел его доделать. Поэтому все ответы приходят в формате JSON.

При написании данной программы руководствовался идеей максимальной простоты и надежности.
Старался максимально избегать так называемого Overengineering-га (Хотя название программы говорит об обратном XD)



