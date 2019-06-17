package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestA(t *testing.T){
	bhMap := map[string]string{}
	f,err := os.Open("bh.txt")
	if err != nil {
		log.Println(err)
	}
	scannerF := bufio.NewScanner(f)
	for scannerF.Scan(){
		line := scannerF.Text()
		fsList := strings.Fields(strings.TrimSpace(line))
		if len(fsList) >= 2{
			bhMap[fsList[0]] = fsList[1]
		}
	}
	if err := scannerF.Err(); err != nil {
		log.Println("读取文件 err ",err)
	}
	fmt.Println(bhMap)

}
