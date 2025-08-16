package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(GetDaily(-1, 9))
}

func GetDaily(addDays, addHours int) string {
	utcDate := time.Now().UTC().AddDate(0, 0, addDays).Add(time.Hour * time.Duration(addHours))
	dailyKey := utcDate.Format("20060102")
	return dailyKey
}
