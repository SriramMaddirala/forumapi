package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var machineBits string
var sequenceNumber int64 = 0

func generateSnowflake(time time.Time) int64 {
	timeBits := fmt.Sprintf("%b", time.Unix())
	timeBits = strings.Repeat("0", 41-utf8.RuneCountInString(timeBits)) + timeBits
	sequenceBits := fmt.Sprintf("%b", sequenceNumber)
	sequenceBits = strings.Repeat("0", 12-utf8.RuneCountInString(sequenceBits)) + sequenceBits
	sequenceNumber = sequenceNumber + 1
	test := "0" + timeBits + sequenceBits + machineBits
	fmt.Println(test)
	snowflakeId, err := strconv.ParseInt(test, 2, 64)
	if err != nil {
		panic(err)
	}
	return snowflakeId
}

//can do something similar with sequence num, store it in .env after you're done and read it from .env

// this works for up to 1023 machine id
func init() {
	machineId, machineIdError := os.LookupEnv("MACHINE_ID")
	if !machineIdError {
		panic("Couldn't get Machine Id")
	}
	macIdInt, _ := strconv.Atoi(machineId)
	machineBits = fmt.Sprintf("%b", macIdInt)
	machineBits = strings.Repeat("0", 10-utf8.RuneCountInString(machineBits)) + machineBits
}
