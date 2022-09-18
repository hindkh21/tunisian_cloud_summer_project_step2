package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	//"io/ioutil"
	"log"
	"os"

	//"strconv"
	//"strings"
	"unicode"

	"github.com/gocolly/colly"
)

type Secteur_Et_URL struct {
	Secteur_Full_URL string `json:"secteur_url"`
	Secteur_Nom      string `json:"secteur_nom"`
}

type Domaine_Et_URL struct {
	domaine_Full_URL string `json:"domaine_url"`
	domaine_Nom      string `json:"domaine_nom"`
}

type Entreprise_Result struct {
	Type       string
	Company    string
	Coordinate int
	Link       string
}

type Domaine_Et_URL_Entreprise struct {
	Domaine_Full_URL string `json:"domaine_url"`
	Domaine_Nom      string `json:"domaine_nom"`
}
type Pagination struct {
	Page_URL     string `json:"page_url"`
	Pages_Nombre string `json:"pages_nombre"`
}
type Entreprise struct {
	Source  string `json:"source"`
	Nom     string `json:"nom"`
	Adresse string `json:"adresse"`
	Tel     string `json:"tel"`
	Fax     string `json:"fax"`
	Secteur string `json:"secteur"`
	Domaine string `json:"domaine"`
	Site    string `json:"site"`
	Email   string `json:"email"`
}

func convertJSONToCSV(test, destination string) error {
	// 2. Read the JSON file into the struct array
	testFile, err := os.Open(test)
	if err != nil {
		return err
	}

	// remember to close the file at the end of the function
	defer testFile.Close()

	var Entreprises []Entreprise
	if err := json.NewDecoder(testFile).Decode(&Entreprises); err != nil {
		return err
	}

	// 3. Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	writer.Comma = ';'
	defer writer.Flush()

	header := []string{"Source", "Nom", "Adresse", "Tel", "Fax", "Secteur", "Domaine", "Site", "Email"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range Entreprises {
		var csvRow []string
		csvRow = append(csvRow, r.Source, r.Nom, r.Adresse, r.Tel, r.Fax, r.Secteur, r.Domaine, r.Site, r.Email)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

var entreprise_xpath string

func scraping_domain_transferts_de_fonds() string {
	entreprise_xpath = "/html/body/div[6]/div[2]/div[2]/div[2]/div[1]/div[2]/div/div[*]" //#TODO checkitpath
	return entreprise_xpath
}
func scraping_domain_Epargne() string {
	entreprise_xpath = "/html/body/div[6]/div[2]/div[2]/div[2]/div[1]/div[2]/div/div[*]"
	return entreprise_xpath
}

func scraping_domain_default() string {
	entreprise_xpath = "/html/body/div[6]/div[2]/div[2]/div[3]/div[1]/div[2]/div/div[*]"
	return entreprise_xpath
}

func main() {
	Liste_entreprises := make([]Entreprise, 0)

	var DData1 Secteur_Et_URL
	list_of_Secteur_Et_URL := make([]Secteur_Et_URL, 0)
	collector_secteur := colly.NewCollector()

	collector_secteur.OnHTML("a.stretched-link.text-center", func(element *colly.HTMLElement) {
		DData1.Secteur_Nom = element.Text
		list_of_Secteur_Et_URL = append(list_of_Secteur_Et_URL, DData1)
	})
	i := 0
	collector_secteur.OnHTML("h3", func(element *colly.HTMLElement) {
		DData1.Secteur_Full_URL = "https://www.goafricaonline.com" + element.ChildAttr("a", "href")
		list_of_Secteur_Et_URL[i].Secteur_Full_URL = "https://www.goafricaonline.com" + element.ChildAttr("a", "href")
		i++
	})

	collector_secteur.OnError(func(r *colly.Response, err error) {
		log.Fatal("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector_secteur.Visit("https://www.goafricaonline.com/annuaire")
	for n := 0; n < len(list_of_Secteur_Et_URL); n++ {

		fmt.Printf("%d) le nom du secteur est : %s", n, list_of_Secteur_Et_URL[n].Secteur_Nom)
		fmt.Println()
		fmt.Printf("%d) l'URL du secteur est : %s", n, list_of_Secteur_Et_URL[n].Secteur_Full_URL)
		fmt.Println()
		fmt.Println("=====================================================")

		//Debut +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
		//Debut +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

		var Base_URL = "https://www.goafricaonline.com"
		var Scraping_URL_Source = list_of_Secteur_Et_URL[n].Secteur_Full_URL // "https://www.goafricaonline.com/annuaire/secteur-energie" //"https://www.goafricaonline.com/ml/annuaire/entreprises-finances"

		//##############################################################################################################3
		fmt.Println()
		fmt.Println("===============================================================")
		fmt.Printf("Starting Scraping URL: %s", Scraping_URL_Source)
		fmt.Println()
		fmt.Println("===============================================================")

		liste_des_noms_de_domaines := make([]Domaine_Et_URL_Entreprise, 0)

		collector_domaine := colly.NewCollector()
		collector_domaine.OnHTML("a", func(element *colly.HTMLElement) {
		})

		collector_domaine.OnXML("/html/body/div[6]/div[2]/div/div[5]/div/div[*]", func(element *colly.XMLElement) {

			Domaine_Nom := element.ChildText("/a/h3")
			Domaine_URL := element.ChildAttr("a", "href")
			Full_URL := Base_URL + Domaine_URL

			var DData Domaine_Et_URL_Entreprise
			DData.Domaine_Full_URL = Full_URL
			DData.Domaine_Nom = Domaine_Nom
			liste_des_noms_de_domaines = append(liste_des_noms_de_domaines, DData)

			fmt.Printf("Domaine_Nom: %s %s ", strings.TrimSpace(Domaine_Nom), strings.TrimSpace(Full_URL))
			fmt.Println()

		})

		collector_domaine.OnRequest(func(request *colly.Request) {
			//fmt.Println("Visiting", request.URL.String())
		})

		collector_domaine.Visit(Scraping_URL_Source)

		fmt.Println("le nombre des domaines est", len(liste_des_noms_de_domaines))

		//##############################################################################################################3

		for k := 0; k < len(liste_des_noms_de_domaines); k++ {

			var Full_URL = liste_des_noms_de_domaines[k].Domaine_Full_URL
			var Domain_To_Scrap = liste_des_noms_de_domaines[k].Domaine_Nom

			fmt.Println()
			fmt.Println("=====================================================")
			fmt.Printf(" Sub_Domain_To_Scrap = %s ", Domain_To_Scrap)
			fmt.Println()
			fmt.Printf(" Sub_Domain_Full_URL= %s", Full_URL)
			fmt.Println()
			fmt.Println("=====================================================")

			Nb_Pages := 1
			var DData1 Pagination
			DData1.Page_URL = Full_URL
			DData1.Pages_Nombre = strconv.Itoa(Nb_Pages)

			liste_pagination := make([]Pagination, 0)
			liste_pagination = append(liste_pagination, DData1)
			fmt.Printf("Page_Number: '%s' Page_Full_URL:%s ", strings.TrimSpace(liste_pagination[Nb_Pages-1].Pages_Nombre), strings.TrimSpace(liste_pagination[Nb_Pages-1].Page_URL))
			fmt.Println()

			collector_pagination := colly.NewCollector()
			collector_pagination.OnXML("/html/body/div[6]/div[2]/div[2]/div[3]/div[1]/div[3]/ul/li[*]", func(element *colly.XMLElement) {
				Page_Number := ""
				Page_URL := ""

				if isInt(strings.TrimSpace(element.ChildText("a"))) == true && strings.TrimSpace(element.ChildText("a")) != "" {
					Nb_Pages++
					Page_Number = element.ChildText("a")
					Page_URL = element.ChildAttr("a", "href")
					Page_Full_URL := Base_URL + Page_URL

					DData1.Page_URL = Page_Full_URL
					DData1.Pages_Nombre = Page_Number
					liste_pagination = append(liste_pagination, DData1)

					fmt.Printf("Page_Number: '%s' Page_Full_URL:%s ", strings.TrimSpace(liste_pagination[Nb_Pages-1].Pages_Nombre), strings.TrimSpace(liste_pagination[Nb_Pages-1].Page_URL))
					fmt.Println()
				}
			})
			collector_pagination.OnRequest(func(request *colly.Request) {
				//fmt.Println()
				//fmt.Println("Visiting", request.URL.String())
			})
			collector_pagination.Visit(Full_URL)
			fmt.Printf("Page Number : %d ", Nb_Pages)
			fmt.Println()
			fmt.Println("=====================================================")

			//##############################################################################################################3

			for j := 0; j < Nb_Pages; j++ {
				i := 0
				fmt.Println()
				fmt.Println("=====================================================")
				fmt.Printf(" Sub_Domain_To_Scrap = %s ", Domain_To_Scrap)
				fmt.Println()
				fmt.Printf(" Sub_Domain_URL_Page_%d_To_Scrap = %s", j, liste_pagination[j].Page_URL)
				fmt.Println()
				fmt.Printf("Page numero : %s ", liste_pagination[j].Pages_Nombre)
				fmt.Println()
				fmt.Println("=====================================================")

				collector_entreprise := colly.NewCollector()

				entreprise_xpath = scraping_domain_default()

				if Domain_To_Scrap == "Transferts de fonds" {
					entreprise_xpath = scraping_domain_transferts_de_fonds()
				} else {

					if Domain_To_Scrap == "Epargne" {
						entreprise_xpath = scraping_domain_Epargne()
					}
				}
				collector_entreprise.OnHTML("div.mb-3", func(element *colly.HTMLElement) {
					fmt.Println(element.Text)
				})

				collector_entreprise.OnXML(entreprise_xpath, func(element *colly.XMLElement) {

					Entreprise_Nom := element.ChildText("/article/div[1]/div/div[1]/div[1]/h2/a")
					Entreprise_Adresse_1 := element.ChildText("/article/div[1]/div/div[2]/div[1]/address/text()[1]")
					Entreprise_Adresse_2 := element.ChildText("/article/div[1]/div/div[2]/div[1]/address/text()[2]")
					Entreprise_Adresse := Entreprise_Adresse_1 + " " + Entreprise_Adresse_2
					Entreprise_Tel1 := element.ChildText("/article/div[1]/div/div[2]/div[2]/a/span[1]")
					Entreprise_Tel2 := element.ChildText("/article/div[1]/div[2]/div[2]/div[2]/a/span[2]")
					Entreprise_Tel := Entreprise_Tel1 + Entreprise_Tel2
					Entreprise_json_sites := element.ChildAttr("/article/div[2]/div", "data-collect-event-on-click")
					var Entreprise_string_sites Entreprise_Result
					json.Unmarshal([]byte(Entreprise_json_sites), &Entreprise_string_sites)

					if strings.TrimSpace(Entreprise_Nom) != "" {
						i++
						var Entreprise_Data Entreprise

						Entreprise_Data.Secteur = list_of_Secteur_Et_URL[n].Secteur_Nom //"*"
						Entreprise_Data.Adresse = Entreprise_Adresse
						Entreprise_Data.Domaine = Domain_To_Scrap
						Entreprise_Data.Email = ""
						Entreprise_Data.Fax = ""
						Entreprise_Data.Nom = Entreprise_Nom
						Entreprise_Data.Site = Entreprise_string_sites.Link
						Entreprise_Data.Source = Scraping_URL_Source
						Entreprise_Data.Tel = Entreprise_Tel

						Liste_entreprises = append(Liste_entreprises, Entreprise_Data)
						fmt.Println()
						fmt.Println("____________________________________________________________")

						fmt.Printf("Entreprise numero : %d  Secteur: %s", i, Entreprise_Data.Secteur)
						fmt.Println()
						fmt.Printf("Entreprise numero : %d  Domaine: %s", i, Entreprise_Data.Domaine)
						fmt.Println()
						fmt.Printf("Entreprise numero : %d  Nom: %s", i, Entreprise_Data.Nom)
						fmt.Println()
						fmt.Printf("Entreprise numero : %d  Adresse: %s", i, Entreprise_Data.Adresse)
						fmt.Println()
						fmt.Printf("Entreprise numero : %d  Tel: %s", i, Entreprise_Data.Tel)
						fmt.Println()
						fmt.Printf("Entreprise numero : %d  Web Site: %s", i, Liste_entreprises[i-1].Site)

					}
				})

				collector_entreprise.OnRequest(func(request *colly.Request) {
					//fmt.Println("Visiting", request.URL.String())
				})

				collector_entreprise.Visit(liste_pagination[j].Page_URL)

			} // for j boucle
		} // for k boucle
	} // for n boucle
	//Fin +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	//Fin +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	file, _ := json.MarshalIndent(Liste_entreprises, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)

	if err := convertJSONToCSV("test.json", "data.csv"); err != nil {
		log.Fatal(err)
	}

}

//##############################################################################################################3
