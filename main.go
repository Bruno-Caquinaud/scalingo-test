package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Scalingo/go-utils/logger"
)

type OrganizationRepositoryAnswer struct {
	FullName     string `json:"full_name"`
	Name         string `json:"name"`
	Owner        Owner  `json:"owner"`
	LanguagesUrl string `json:"languages_url"`
}

type Owner struct {
	Login string `json:"login"`
}

type ApiAnswer struct {
	Repositories []DetailsRepository `json:"repositories"`
}

type DetailsRepository struct{
	FullName string `json:"fullname"`
	Owner string `json:"owner`
	Repository string `json:"repository"`
	Languages map[string]CodeSize `json:"languages"`
}

type CodeSize struct{
	Bytes int `json:"bytes"`
}

func parseOrganizationRepositoriesAnswer(jsonAnswer *[]byte, reply *[]OrganizationRepositoryAnswer) {

	json.Unmarshal(*jsonAnswer, reply)

	fmt.Printf("Chasse et Peche")

	fmt.Printf(" fullname: %s, Owner:%s, LanguagesUrl:%s", (*reply)[0].FullName, (*reply)[0].Owner.Login, (*reply)[0].LanguagesUrl)
}

func callOrganizationEndpoint(org string, answer *[]byte) {

	ListOrganizationRepoUrl := "https://api.github.com/orgs/" + org + "/repos"
	ListOrganizationRepoUrl1, err := url.Parse(ListOrganizationRepoUrl)

	if err != nil {

		fmt.Printf("Invalid Url to parse: %s\n", err)
		os.Exit(1)
	}

	rawquery := ListOrganizationRepoUrl1.Query()
	rawquery.Add("type", "public")
	rawquery.Add("per_page", "3")

	ListOrganizationRepoUrl1.RawQuery = rawquery.Encode()
	ListOrganizationRepoUrl2 := ListOrganizationRepoUrl1.String()

	req, err := http.NewRequest(http.MethodGet, ListOrganizationRepoUrl2, nil)
	if err != nil {

		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(0)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {

		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: response body: %s\n", resBody)
	*answer = resBody

}

func callLanguageEndpoint(urlp string, answer *[]byte) {

	//ListLanguageRepoUrl := "https://api.github.com/repos/" + fullname + "/languages"

	ListLanguageRepoUrl1, err := url.Parse(urlp)

	if err != nil {

		fmt.Printf("Invalid Url to parse: %s\n", err)
		os.Exit(1)
	}

	ListLanguageRepoUrl2 := ListLanguageRepoUrl1.String()
	fmt.Printf(" ListLanguageRepoUrl: %s\n", ListLanguageRepoUrl2)

	req, err := http.NewRequest(http.MethodGet, ListLanguageRepoUrl2, nil)

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf(" client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	*answer = resBody

	fmt.Printf("client: response body: %s\n", resBody)
}

func decodeLanguagesAnswer(answer *[]byte, mapper *map[string]int) {
	if err := json.Unmarshal(*answer, mapper); err != nil {
		panic(err)
	}
}

func processAnswer(listOrgRepo *[]OrganizationRepositoryAnswer, listLanguagesRepo *[]map[string]int, languageSelected string, apiAnswer * ApiAnswer) {

	if len(*listOrgRepo) != len(*listLanguagesRepo) {
		fmt.Printf("len reply and mapper are not the same size")
		os.Exit(1)
	}

	for i, currentRepo := range *listOrgRepo {
		var currentRepoDetails DetailsRepository
		
		currentRepoDetails.FullName = currentRepo.FullName
		currentRepoDetails.Owner = currentRepo.Owner.Login
		currentRepoDetails.Repository = currentRepo.Name
		currentRepoDetails.Languages = make(map[string]CodeSize)
		currentRepoDetails.Languages[languageSelected] = CodeSize{((*listLanguagesRepo)[i])[languageSelected]}
		
		(*apiAnswer).Repositories = append((*apiAnswer).Repositories, currentRepoDetails)
	}
}

func main() {
	var org string = "adobe"
	var answer []byte
	var objectOrgRepo []OrganizationRepositoryAnswer
	var LanguagesanswerMapper []map[string]int
	var apiAnswer ApiAnswer
	var selectedLanguage string = "C++"
	var apiAnswerJson []byte
	
	callOrganizationEndpoint(org, &answer)
	parseOrganizationRepositoriesAnswer(&answer, &objectOrgRepo)
	
	for _, orgRepo := range objectOrgRepo{
		var objectAnswer map[string]int
		callLanguageEndpoint(orgRepo.LanguagesUrl, &answer)
		decodeLanguagesAnswer(&answer, &objectAnswer)
		
		LanguagesanswerMapper = append(LanguagesanswerMapper, objectAnswer)
	}
	processAnswer(&objectOrgRepo, &LanguagesanswerMapper, selectedLanguage, &apiAnswer)
	apiAnswerJson, err := json.Marshal(apiAnswer)
	
	if err != nil {
		panic(err)
	}
	fmt.Printf("api answer : %s\n", string(apiAnswerJson))
	return
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}

