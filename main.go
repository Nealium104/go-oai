package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// OAIPMH was generated 2024-08-20 15:53:27 by https://xml-to-go.github.io/ in Ukraine.
	type OAIPMH struct {
		XMLName        xml.Name `xml:"OAI-PMH"`
		Text           string   `xml:",chardata"`
		Xmlns          string   `xml:"xmlns,attr"`
		Xsi            string   `xml:"xsi,attr"`
		SchemaLocation string   `xml:"schemaLocation,attr"`
		ResponseDate   string   `xml:"responseDate"`
		Request        struct {
			Text            string `xml:",chardata"`
			Verb            string `xml:"verb,attr"`
			ResumptionToken string `xml:"resumptionToken,attr"`
		} `xml:"request"`
		ListIdentifiers struct {
			Text   string `xml:",chardata"`
			Header []struct {
				Text       string `xml:",chardata"`
				Identifier string `xml:"identifier"`
				Datestamp  string `xml:"datestamp"`
				SetSpec    string `xml:"setSpec"`
			} `xml:"header"`
			ResumptionToken string `xml:"resumptionToken"`
		} `xml:"ListIdentifiers"`
	}

	baseURL := "https://exploreuk.uky.edu/"
	f, err := os.Create("EUK_IDs")
	if err != nil {
		log.Fatal("Error with creating file: ", err)
	}
	defer f.Close()

	for {
		request, err := http.Get(fmt.Sprintf("%vcatalog/oai?verb=ListIdentifiers&resumptionToken=oai_dc.s(default).f(2018-09-25T12:48:35Z).u(2024-08-20T17:22:46Z):2", baseURL))

		if err != nil {
			log.Fatalf("Error with request: %v", err)
			return
		}
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		var response OAIPMH
		error := xml.Unmarshal(body, &response)
		if error != nil {
			log.Fatal("Error with unmarshal: ", error)
		}

		for _, header := range response.ListIdentifiers.Header {
			str := header.Identifier
			cleaned := strings.Split(str, "/")
			identifier := cleaned[1]
			_, err := f.WriteString(identifier + "\n")
			if err != nil {
				log.Fatal("Error writing to file: ", err)
			}
		}

		fmt.Println("Request header: ", request)
		fmt.Println("Request Body: ", string(body))
		fmt.Printf("Type: %T", body)

		time.Sleep(2 * time.Second)
	}
}

// TODO
// open the ids file
// pipe each id into base/dips/id/data/mets.xml
// Find the items that you need to extract by finding acceptable mimetypese
// Grab the xlink:href value and pipe it into base/dips/id/data/href
