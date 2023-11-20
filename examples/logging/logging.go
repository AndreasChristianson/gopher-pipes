package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/AndreasChristianson/gopher-pipes/reactive"
	"log"
	"math/rand"
	"time"
)

func main() {
	reactive.SetLogger(func(level reactive.Level, id interface{}, messageFormat string, args ...interface{}) {
		message := fmt.Sprintf(messageFormat, args...)
		log.Default().Printf("%s [%s]: %s", level, id, message)
	})
	returns := []string{"Hello", "world", "!"}
	pos := 0
	pipe := reactive.FromGeneratorWithDefaultBackoff(func() (*string, error) {
		if rand.Float32() > 0.8 {
			return nil, errors.New("simulated generator error")
		}
		if rand.Float32() > 0.5 {
			return nil, nil // simulated empty poll
		}
		if pos == len(returns) {
			return nil, reactive.GeneratorFinished{}
		}
		ret := returns[pos]
		pos++
		return &ret, nil
	})
	var response []string
	pipe.Observe(func(item string) error {
		response = append(response, item)
		return nil
	})
	pipe.Observe(func(item string) error {
		return errors.New("simulated sink error")
	})

	pipe.Observe(func(s string) error {
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		fmt.Println("Simulating long processing times", delay)

		time.Sleep(delay)
		return nil
	})
	reactive.Map[string, []byte](pipe, func(from string) ([]byte, error) {
		md5Hash := md5.Sum([]byte(from))
		asSlice := md5Hash[:]
		return asSlice, nil
	}).Observe(func(hash []byte) error {
		fmt.Printf("hashed: %x\n", hash)
		return nil
	})
	pipe.Start()
	pipe.AwaitCompletion()
	fmt.Println(response)
}
