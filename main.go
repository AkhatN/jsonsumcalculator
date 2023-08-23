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

	// Флаги для получения пути к json файлу, количество горутин и количество блоков из аргументов командной строки
	jsonFilePath := flag.String("file", "", "Путь к файлу")
	numOfGoroutines := flag.Int("numgo", 1, "Количество горутин")
	numOfBlocks := flag.Int("numblocks", 1, "Количество блоков")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: go run main.go -file <путь к файлу> -numgo <Количество горутин> -numblocks <Количество блоков>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *jsonFilePath == "" {
		log.Println("Нужно название файла\nUsage: go run main.go -file <Путь к файлу> -numgo <Количество горутин> -numblocks <Количество блоков>")
		return
	}

	if *numOfGoroutines < 1 {
		log.Println("Количество горутин не должно быть меньше 1\nUsage: go run main.go -file <Путь к файлу> -numgo <Количество горутин> -numblocks <Количество блоков>")
		return
	}

	if *numOfBlocks < 1 {
		log.Println("Количество горутин не должно быть меньше 1\nUsage: go run main.go -file <Путь к файлу> -numgo <Количество горутин> -numblocks <Количество блоков>")
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

	jobs := make(chan []Object, *numOfBlocks)
	results := make(chan int, *numOfBlocks)

	var wg sync.WaitGroup

	for w := 0; w < *numOfGoroutines; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results)
	}

	// Нахождение размера блоков
	isBlockSizeFixed := false
	blockSize := len(objs) / *numOfBlocks

	if blockSize == 0 {
		blockSize = 1
		isBlockSizeFixed = true
	}

	leftElements := 0
	if !isBlockSizeFixed {

		numOfElements := blockSize * *numOfBlocks

		if numOfElements == len(objs) {
			isBlockSizeFixed = true
		}

		if !isBlockSizeFixed {
			leftElements = len(objs) - numOfElements
		}
	}

	count := 0
	go func() {
		// Распределение данных по блокам в горутины
		for start := 0; start < len(objs); start += blockSize {

			end := start + blockSize

			if leftElements > 0 {
				end++
			}

			if end > len(objs) {
				end = len(objs)
			}

			jobs <- objs[start:end]

			count++
			if leftElements > 0 {
				start++
				leftElements--
			}
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// Суммирование общего результата
	totalResult := 0
	for r := range results {
		totalResult += r
	}

	fmt.Println("Обший результат:", totalResult)
}

func worker(wg *sync.WaitGroup, jobs <-chan []Object, result chan<- int) {
	defer wg.Done()

	for {

		job, ok := <-jobs
		if !ok {
			return
		}

		sum := 0
		for _, obj := range job {
			sum += obj.A + obj.B
		}

		result <- sum
	}
}
