# Исследование методов сериализации данных

> Приложение писалось в качестве домашнего задания из курса "Сервис-ориентированные архитектуры" ПМИ ВШЭ.

## Постановка задачи

**Цель:** на языке Go реализовать приложение для тестирования эффективности работы с различными форматами сериализации данных. В процессе тестирования форматов сериализации необходимо учитывать следующие характеристики:
1. Размер сериализованной структуры данных;
2. Время сериализации/десериализации.

В сериализуемой структуре желательно представить несколько различных видов данных, включая:
- строковые данные,
- массивы данных,
- словари,
- целочисленные данные,
- данные с плавающей запятой.

Приложение должно обеспечивать:
- выполнение операций сериализации/десериализаци набора данных, во все форматы:
  - Нативный Gob
  - XML
  - JSON
  - Google Protocol Buffers
  - Apache Avro
  - YAML
  - MessagePack
- представление в наглядном виде результатов выполнения сериализации/десериализации.

Необходимо предоставить отчет в формате таблицы Excel о форматах сериализации данных.


## Запуск приложения

### Локально

> Необходимо наличие `go` версии `>1.15`

Запуск приложения:

```
mkdir files
go run main.go
```

По завершению в консоль выведутся данные, которые показывают, сколько наносекунд/байт в сумме для 1000 итераций было потрачено для каждого формата.
```
	Volume	SerTime	DeserTime
nat:	{13679705 90723503 140447777}
xml:	{33291253 369779757 1676470454}
json:	{20025640 204605494 498539402}
proto:	{15650379 146048448 229550197}
avro:	{12927683 38657463 115304117}
yaml:	{20513768 1884369776 1507802880}
msg:	{14894918 116914643 219894076}
```

Отчет будет находится в файле `./report/report.xlsx`

### Docker

> Необходимо наличие докера

```
docker pull zhukowladimir/serialization_test
docker run -d --name your-container-name zhukowladimir/serialization_test
docker cp your-container-name:/serialization-test/report/report.xlsx /your/local/path/your_report_name.xlsx
```
С помощью этого вы получите отчет, который будет храниться здесь: `/your/local/path/your_report_name.xlsx`.

Все рабочие файлы внутри докера хранятся в `/serialization-test`.

## Описание структуры

```go
type Map map[string]int32

type Qwe struct {
	Rty string
	Ott int64
}

type Test struct {
	ID                           int32    
	Name                         string    
	ServiceIDs                   []int32   
	Tests                        []Qwe     
	Flts                         []float64 
	Dict                         Map       
	VeryLongNameForSmallVariable float32 
}
```
