package function

import "github.com/hecatoncheir/Storage"

type FAASFunctions struct {
	FAASGateway string
}

func (functions FAASFunctions) CompaniesReadByName(companyName string, language string) ([]storage.Company, error) {
	return []storage.Company{}, nil
}
