package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dataPath := "./Input/"
	//outputPath := "./Output/"

	files, err := ioutil.ReadDir("./Input")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		targetFile := dataPath + f.Name()
		path := flag.String("path", targetFile, "Path of the file")
		flag.Parse()
		fileBytes, fileNPath := CSVRead(path)
		SaveFile(fileBytes, fileNPath)
	}

	fmt.Println(strings.Repeat("=", 10), "끝!", strings.Repeat("=", 10))
}

func CSVRead(path *string) ([]byte, string) {
	csvFile, err := os.Open(*path)
	if err != nil {
		log.Fatal("파일이 없습니다. || 파일 경로가 다릅니다.")
	}
	defer csvFile.Close()

	fmt.Printf("파일명 : %s\n", csvFile.Name())
	b := make([]byte, 3)
	_, err = csvFile.Read(b)
	if err != nil {
		log.Fatalln("Read failed", err)
	}

	if !(b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xEF) {
		csvFile.Seek(0, os.SEEK_SET)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	content, _ := reader.ReadAll()

	if len(content) < 1 {
		log.Fatal("파일이 잘못됐거나 라인의 길이가 다릅니다")
	}

	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		fmt.Printf(" 내용물 : %s", headE)
		headersArr = append(headersArr, headE)
	}

	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString("{")
	buffer.WriteString(`"` + "value" + `":`)
	buffer.WriteString(`{`)
	buffer.WriteString(`"` + "table" + `":`)
	buffer.WriteString("[")

	for i, d := range content {
		buffer.WriteString("{")

		//fmt.Printf("몇번째 줄 : %d", i)
		for j, y := range d {
			buffer.WriteString(`"` + headersArr[j] + `":`)

			_, fErr := strconv.ParseFloat(y, 32)
			_, bErr := strconv.ParseBool(y)
			if fErr == nil {
				buffer.WriteString(y)
			} else if bErr == nil {
				buffer.WriteString(strings.ToLower(y))
			} else {
				buffer.WriteString((`"` + y + `"`))
			}
			//end of property
			if j < len(d)-1 {
				buffer.WriteString(",")
			}
		}

		//fmt.Printf(" 현재 : %s", &buffer)
		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}

		//fmt.Println("")
	}

	buffer.WriteString("]")
	buffer.WriteString("}")
	buffer.WriteString("}")
	buffer.WriteString(`]`)
	rawMessage := json.RawMessage(buffer.String())
	x, _ := json.MarshalIndent(rawMessage, "", "  ")
	newFileName := filepath.Base(*path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"

	println(newFileName)
	r := filepath.Dir("./Output/")
	return x, filepath.Join(r, newFileName)
}

func SaveFile(myFile []byte, path string) {
	if err := ioutil.WriteFile(path, myFile, os.FileMode(0644)); err != nil {
		panic(err)
	}
}
