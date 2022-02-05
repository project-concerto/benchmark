package utils

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var APITokens map[int32]string
var OrderNumbers map[int32]string

func LoadKeyValuePairs(m map[int32]string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 2 {
			continue
		}

		id, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		m[int32(id)] = fields[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func LoadIntStringPairs(filename string) map[int]string {
	m := make(map[int]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 2 {
			continue
		}

		id, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		m[id] = fields[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return m
}

func init() {
	APITokens = make(map[int32]string)
	OrderNumbers = make(map[int32]string)
	LoadKeyValuePairs(APITokens, "./utils/spree_api_token")
	LoadKeyValuePairs(OrderNumbers, "./utils/spree_order_number")
}
