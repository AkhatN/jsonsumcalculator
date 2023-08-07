package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

type Object struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {

	// Флаги для получения пути к json файлу и количество горутин из аргументов командной строки
	jsonFilePath := flag.String("file", "", "Путь к файлу")
	numOfGoroutines := flag.Int("numgo", 1, "Количество горутин")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go -file <путь к файлу> -numgo <Количество горутин>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *jsonFilePath == "" {
		log.Println("Нужно название файла\nUsage: go run main.go -file <Путь к файлу> -numgo <Количество горутин>")
		return
	}

	if *numOfGoroutines < 1 {
		log.Println("Количество горутин не должно быть меньше 1\nUsage: go run main.go -file <Путь к файлу> -numgo <Количество горутин>")
		return
	}

	data, err := os.ReadFile(*jsonFilePath)
	if err != nil {
		log.Println(err)
		return
	}

	var objs []Object

	if err := json.Unmarshal(data, &objs); err != nil {
		log.Println(err)
		return
	}

	// Нахождение размера блоков
	isBlockSizeFixed := false
	blockSize := len(objs) / *numOfGoroutines

	if blockSize == 0 {
		blockSize = 1
		isBlockSizeFixed = true
	}

	leftElements := 0
	if !isBlockSizeFixed {

		numOfElements := blockSize * *numOfGoroutines

		if numOfElements == len(objs) {
			isBlockSizeFixed = true
		}

		if !isBlockSizeFixed {
			leftElements = len(objs) - numOfElements
		}
	}

	wg := &sync.WaitGroup{}
	result := make(chan int)

	// Распределение данных по блокам в горутины
	for start := 0; start < len(objs); start += blockSize {

		end := start + blockSize

		if leftElements > 0 {
			end++
		}

		if end > len(objs) {
			end = len(objs)
		}

		wg.Add(1)

		go calculateBlockSum(objs[start:end], result, wg)

		if leftElements > 0 {
			start++
			leftElements--
		}
	}

	// Ждет завершения всех горутин и закрытие канала
	go func() {
		wg.Wait()
		close(result)
	}()

	// Суммирование общего результата
	totalResult := 0
	for r := range result {
		totalResult += r
	}

	fmt.Println("Обший результат:", totalResult)
}

// calculateBlockSum находит сумму каждого отдельного блока
func calculateBlockSum(block []Object, result chan int, wg *sync.WaitGroup) {

	defer wg.Done()

	sum := 0
	for _, obj := range block {
		sum += obj.A + obj.B
	}

	result <- sum
}
