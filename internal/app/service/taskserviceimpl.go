package service

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"grpc-rusprofile-task/configs"
	pb "grpc-rusprofile-task/proto"
	"log"
	"net/http"
)

type TaskServiceImpl struct {
	cli *http.Client
	cfg *configs.Config
}

func NewTaskService(client *http.Client, cfg *configs.Config) TaskService {
	return &TaskServiceImpl{cli: client, cfg: cfg}
}


func (t TaskServiceImpl) CompanyByInn(ctx context.Context, request *pb.CompanyByINNRequest) (*pb.CompanyByINNResponse, error) {
	response := &pb.CompanyByINNResponse{}
	res, err := t.cli.Get(t.cfg.HTTP.Url+"/search?query="+request.Inn)
	if err != nil {
		return nil, &pb.ErrorNotFound{Msg: err.Error()}
	}

	defer res.Body.Close()

	node, err := html.Parse(res.Body)
	if err != nil {
		return nil, &pb.Error{Msg: err.Error()}
	}

	doc := goquery.NewDocumentFromNode(node)

	doc.Find("#main").Each(func(i int, selection *goquery.Selection) {
		class, _ := selection.Attr("class")
		if class == "company-main renewed"{
			foundCompany(doc, response)
		}else {
			sel := doc.Find(".company-item::first-child")
			sel.Each(func(i int, selection *goquery.Selection) {
				selection.Find("a").Each(func(i int, selection *goquery.Selection) {
					href, _ := selection.Attr("href")
					docById, err := goquery.NewDocument("https://www.rusprofile.ru" + href)
					if err != nil{
						log.Fatal(err)
					}
					foundCompany(docById, response)
				})
			})
		}
	})
	return response, nil
}

func foundCompany(doc *goquery.Document, response *pb.CompanyByINNResponse) {
	doc.Find(".company-name::first-child").Each(func(i int, selection *goquery.Selection) {
		response.Name = selection.Text()
	})

	sel := doc.Find("span")
	sel.Each(func(i int, selection *goquery.Selection) {
		id, _ := selection.Attr("id")
		if id == "clip_inn" {
			response.Inn = selection.Text()
		}
		if id == "clip_kpp" {
			response.Kpp = selection.Text()
		}
	})

	doc.Find(".company-row").Each(func(i int, selection *goquery.Selection) {
		selection.Find("span").Each(func(i int, selection *goquery.Selection) {
			class, _ := selection.Attr("class")
			if class == "company-info__title" {
				if selection.Text() == "Руководитель" {
					asd := selection.Next()
					qwe := asd.Next()
					response.FullName = qwe.Text()
				}
			}
		})
	})
}
