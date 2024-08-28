package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func writeDipId(dipId string) {
	dipIdByte := []byte(dipId + "\n")

	runningIds, err2 := os.OpenFile("./ids/collections-complete-log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		log.Println("Error writing dip ids to log: ", err2)
	}
	defer runningIds.Close()

	_, err := runningIds.Write(dipIdByte)
	if err != nil {
		log.Println("Error writing to dip id log: ", err)
	}
}

func main() {
	baseURL := "https://exploreuk.uky.edu"
	mimetypeList := make(map[string]int)

	errorLogFile, err := os.OpenFile("errLog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(errorLogFile)
	log.Println("Start log")

	file, err := os.Open("ids/run-dips.txt")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		dipId, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("Error reading file:", err)
			continue
		}
		dipId = strings.TrimSuffix(dipId, "\n")
		fmt.Println("Collection: ", dipId)

		// pipe each id into base/dips/id/data/mets.xml
		response, err := http.Get(baseURL + "/dips/" + dipId + "/data/mets.xml")
		if err != nil {
			log.Println("Error with METS get request: ", err)
			continue
		}
		defer response.Body.Close()

		decoder := xml.NewDecoder(response.Body)
		var currentMimetype string

		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error with xml encoding in %s: %v\n", dipId, err)
				break
			}

			// Find the items that you need to extract by finding acceptable mimetypes
			switch t := token.(type) {
			case xml.StartElement:
				if t.Name.Local == "file" {
					for _, attr := range t.Attr {
						if attr.Name.Local == "MIMETYPE" {
							currentMimetype = attr.Value
							_, exists := mimetypeList[currentMimetype]
							if exists {
								mimetypeList[currentMimetype] += 1
							} else {
								mimetypeList[currentMimetype] = 1
							}
							fmt.Printf("Collection: %v\tRunning MimeTypes: %v\n", dipId, mimetypeList)
							break
						}
					}
				}
				// TODO: Make a list of the mimetypes you're interested in
				if t.Name.Local == "FLocat" && currentMimetype == "text/plain" || currentMimetype == "application/xml" {
					for _, attr := range t.Attr {
						if attr.Name.Local == "href" {
							href := attr.Value
							// fmt.Printf("href is: %v", href)
							response, err := http.Get(baseURL + "/dips/" + dipId + "/data/" + href)
							if err != nil {
								log.Println("Error with resource request: ", err)
							}
							defer response.Body.Close()
							fileName := dipId + "_" + filepath.Base(href)
							file, err := os.Create("./resources/" + fileName)
							if err != nil {
								log.Println("Error creating file: ", err)
							}
							defer file.Close()

							_, err = io.Copy(file, response.Body)
							if err != nil {
								log.Println("Error copying file: ", err)
							}
							// fmt.Printf("File %v successfully downloaded\n", fileName)
							// time.Sleep(time.Millisecond * 25)
						}
					}
				}
			}
		}
		writeDipId(dipId)
	}
}
