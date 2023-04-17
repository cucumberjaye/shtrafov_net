package parser

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cucumberjaye/shtrafov_net/internal/pb"
)

// пример ИНН для получения множественного вывода : 1656002652)

const queryURL = "https://www.rusprofile.ru/search?query=%v" //шаблон поисковой строки

var (
	ErrNotFound = errors.New("not found")
)

//Основная функция парсинга
func GetCompanyInfo(inn string) (*pb.ProfileResponse, error) {
	doc, err := getHTMLdocument(inn)
	if err != nil {
		return nil, err
	}

	success := isSuccessSearch(doc)

	if success {
		resp, err := parseCompanyInfo(doc)
		if err != nil {
			return nil, err
		}
		return resp, nil

	} else {
		description, err := getDescription(doc)
		if err != nil {
			return nil, err
		}

		resultNum := getResultNumber(description)
		if resultNum == "0" {
			return nil, ErrNotFound
		} else {
			log.Printf("Request : %s, results : %s", inn, resultNum)

			ogrn, err := findRequestedOGRN(doc, inn)
			if err != nil {
				return nil, err
			}
			docWithOGRN, err := getHTMLdocument(ogrn)

			if err != nil {
				return nil, err
			}
			response, err := parseCompanyInfo(docWithOGRN)

			if err != nil {
				return nil, err
			}

			return response, nil

		}
	}

}

//Получение description для анализа (при вводе ИНН в поисковую строку возможна множественная выдача)
func getDescription(doc *goquery.Document) (resultNum string, err error) {
	var metaDescription string
	doc.Find("meta").EachWithBreak(func(index int, item *goquery.Selection) bool {
		if item.AttrOr("name", "") == "description" {
			metaDescription = item.AttrOr("content", "")
			return false
		}
		return true
	})
	return metaDescription, nil

}

//количество результатов поисковой выдачи
func getResultNumber(desc string) (resultNum string) {
	slice := strings.Split(desc, " ")
	log.Println("количество результатов поисковой выдачи : " + slice[5])
	return slice[5]

}

//парсинг html-документа при выдаче единственного результата
func parseCompanyInfo(doc *goquery.Document) (*pb.ProfileResponse, error) {
	profile := &pb.ProfileResponse{}
	doc.Find("span").EachWithBreak(func(index int, item *goquery.Selection) bool {
		val, ok := item.Attr("id")
		if ok && val == "clip_inn" {
			profile.Inn = item.Text()
			return false
		}
		return true
	})

	if profile.Inn == "" {
		return nil, fmt.Errorf("company not found: %w", ErrNotFound)
	}

	doc.Find("span").EachWithBreak(func(index int, item *goquery.Selection) bool {
		val, ok := item.Attr("id")
		if ok && val == "clip_kpp" {
			profile.Kpp = item.Text()
			return false

		}
		return true
	})

	doc.Find("div").EachWithBreak(func(index int, item *goquery.Selection) bool {
		val, ok := item.Attr("itemprop")
		if ok && val == "legalName" {
			profile.CompanyName = item.Text()
			return false

		}
		return true
	})

	doc.Find("span").EachWithBreak(func(index int, item *goquery.Selection) bool {
		val, ok := item.Attr("class")
		if ok && val == "company-info__text" {
			profile.Supervisor = item.Text()
			return false

		}
		return true
	})

	return profile, nil
}

//Успешным считается поиск, при котором мы сразу получаем страницу компании, неуспешным - когда получаем множественный вывод либо
//отсутствие результата поиска
func isSuccessSearch(doc *goquery.Document) (successSearch bool) {
	doc.Find("title").EachWithBreak(func(index int, item *goquery.Selection) bool {
		if strings.Contains(item.Text(), "результаты поиска") {
			successSearch = false
			return false
		}
		successSearch = true
		return true
	})
	return successSearch
}

//получаем ОГРН для получения точного результата парсинга url с множественной выдачей (пример : https://www.rusprofile.ru/search?query=1656002652&search_inactive=0)
func findRequestedOGRN(doc *goquery.Document, requestedINN string) (ogrn string, err error) {

	//для получения ОГРН при множественном результате поисковой выдаче нужно получить следующее за ИНН поле <dd></dd>
	var keyForNextElement int
	m := make(map[int]string)
	doc.Find("dd").Each(func(i int, dd *goquery.Selection) {
		m[i] = dd.Text()
		if dd.Text() == requestedINN {
			keyForNextElement = i
		}
	})

	ogrn = m[keyForNextElement+1]

	if ogrn != "" {
		return ogrn, nil
	}
	return "", fmt.Errorf("inn not found: %w", ErrNotFound)
}

//Получение HTML-документа для парсинга
func getHTMLdocument(props string) (doc *goquery.Document, err error) {
	resp, err := http.Get(fmt.Sprintf(queryURL, props))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	return doc, nil
}
