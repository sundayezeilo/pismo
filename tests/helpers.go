package tests

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GetEnvPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("Error reading current working directory: ", err)
	}
	return path + "/.env"
}

func GenerateRandomNumberString(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	result := ""
	for i := 0; i < length; i++ {
		num := r.Intn(10)
		result += strconv.Itoa(num)
	}
	return result
}
