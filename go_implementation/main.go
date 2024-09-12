package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const ArrayLength = 50_000_000
const ThreadCount = 7

type Array [ArrayLength]int
type ChunkResult struct {
	result       int
	workerNumber int
}

func fillArray() Array {
	var a Array
	for i := 0; i < ArrayLength; i++ {
		a[i] = i
	}
	return a
}

func sumArray(a []int) int {
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

func worker(a []int, workerNumber int, ch chan<- ChunkResult, wg *sync.WaitGroup) {
	defer func() {
		fmt.Printf("потік № %d завершив роботу\n", workerNumber)
		wg.Done()
	}()
	if workerNumber%2 == 1 {
		/*
			додаємо затримку до потоків з непарним номером
			для кращої демонстрації паралелізму
		*/
		time.Sleep(time.Second)
	}
	ch <- ChunkResult{
		result:       sumArray(a),
		workerNumber: workerNumber,
	}
}

func waitWorkers(ch chan ChunkResult, wg *sync.WaitGroup) {
	/*
		чекаємо, поки завершаться всі потоки, щоб коректно закрити канал
	*/
	wg.Wait()
	fmt.Println("всі потоки завершили роботу, закриваємо канал")
	close(ch)
}

func main() {
	fmt.Printf("довжина масиву: %d\nк-ть потоків: %d\n", ArrayLength, ThreadCount)

	a := fillArray()
	syncCalculatedSum := sumArray(a[:])
	parallelCalculatedSum := 0

	chunkSize := int(math.Ceil(float64(ArrayLength) / float64(ThreadCount)))

	fmt.Printf("синхронно порахована сума масиву: %d\n", syncCalculatedSum)
	/*
		створюємо буферизований канал, щоб вичитувати проміжні результати,
		та не блокувати потоки на запис
	*/
	ch := make(chan ChunkResult, ThreadCount)
	/*
		створюємо WaitGroup, щоб дочекатись, поки всі потоки відпрацюють,
		щоб коректно закрити канал
	*/
	wg := new(sync.WaitGroup)

	for i := 0; i < ThreadCount; i++ {
		wg.Add(1)
		chunkStart := i * chunkSize
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > ArrayLength {
			chunkEnd = ArrayLength
		}
		fmt.Printf("запускаємо потік № %d\n", i)
		go worker(a[chunkStart:chunkEnd], i, ch, wg)
	}

	go waitWorkers(ch, wg)

	for {
		res, ok := <-ch
		if !ok {
			fmt.Printf("канал закрито\n")
			break
		}
		parallelCalculatedSum += res.result
	}

	fmt.Printf("паралельно порахована сума масиву: %d\n", parallelCalculatedSum)
}
