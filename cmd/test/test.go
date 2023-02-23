package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	//u1 := uuid.NewV4()
	//fmt.Printf("UUIDv4: %s\n", u1)
	//
	//u1 = uuid.NewV4()
	//fmt.Printf("UUIDv4: %s\n", u1)
	//
	//u1 = uuid.NewV4()
	//fmt.Printf("UUIDv4: %s\n", u1)

	cardNum := genCardNum("526855", 16)
	fmt.Println(cardNum)
	fmt.Println(checkCarNum("4624 0096 9192 1299"))
}

func checkCarNum(cardNum string) bool {
	cardNum = strings.ReplaceAll(cardNum, " ", "")

	sum, err := getCardNumSum(cardNum)
	if err != nil {
		return false
	}
	return sum%10 == 0
}
func getCardNumSum(cardNum string) (int64, error) {
	sum := int64(0)
	length := len(cardNum)
	index := length - 1
	for {
		t, err := strconv.ParseInt(string(cardNum[index]), 10, 64)
		if err != nil {
			return 0, err
		}
		if index%2 == 0 {
			t = t * 2
			if t >= 10 {
				t = t%10 + t/10
			}
		}
		sum += t
		if index <= 0 {
			break
		}
		index--
	}
	return sum, nil
}
func genCardNum(startWith string, totalNum int) string {
	result := startWith
	length := len(result)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		result += fmt.Sprintf("%d", rand.Intn(10))
		if length == totalNum-1 {
			break
		}
		length++
	}
	sum, _ := getCardNumSum(result + "0")
	t := 10 - sum%10
	if t == 10 {
		t = 0
	}
	result += fmt.Sprintf("%d", t)
	return result
}
