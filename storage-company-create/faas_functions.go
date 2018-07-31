package function

import (
	"fmt"
	"github.com/hecatoncheir/Storage"
	"net/http"
	"bytes"
	"encoding/json"
	"log"
	"os"
)

var FAASLogger = log.New(os.Stdout, "FAASFunctions: ", log.Lshortfile)

type FAASFunctions struct {
	FAASGateway string
}

func (functions FAASFunctions) CompaniesReadByName(companyName, language, DatabaseGateway string) []storage.Company {
	functionPath := fmt.Sprintf("%v/%v/%v", functions.FAASGateway, "functino", "storage-company-read-by-name")

	body := struct {
		Language        string
		CompanyName     string
		DatabaseGateway string
	}{
		Language:language,
		CompanyName:companyName,
		DatabaseGateway:DatabaseGateway	}

	encodedBody,err:=json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
	}

	http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	return []storage.Company{}
}
