package main

import (
	"encoding/json"
	. "github.com/dave/jennifer/jen"
	"io/ioutil"
	"log"
	"strings"
)

type APIparameters struct {
	PathParam []pathParameters
	Host string
	BasePath string
	Security string
}

type pathParameters struct {
	Path string
	Method []MethodParameters
}

type MethodParameters struct {
	MethodParam []OneMethodParam
	MethodType string
	MethodConsumes []interface{}
}
type OneMethodParam struct {
	In interface{}
	TypeParam interface{}
	Name interface{}
}

func main () {
	//Загружаем swagger.json файл документации
	file, _ := ioutil.ReadFile("swagger.json")
	var api APIparameters
	var pathParam pathParameters
	var methodParameters MethodParameters
	var oneMethodParam OneMethodParam
	m := make(map[string]interface{})

	_ = json.Unmarshal(file, &m)

	//Endpoint => Method => Parameters
	for k := range m["paths"].(map[string]interface{}) {
		pathParam.Path = k
		for j := range m["paths"].(map[string]interface{})[k].(map[string]interface{}) {
			methodParameters.MethodType = strings.ToUpper(j)
			for _, l := range m["paths"].(map[string]interface{})[k].(map[string]interface{})[j].(map[string]interface{})["parameters"].([]interface{}) {
				oneMethodParam.TypeParam = l.(map[string]interface{})["type"]
				oneMethodParam.In = l.(map[string]interface{})["in"]
				oneMethodParam.Name = l.(map[string]interface{})["name"]
				methodParameters.MethodParam = append(methodParameters.MethodParam, oneMethodParam)
			}
			for _, l := range m["paths"].(map[string]interface{})[k].(map[string]interface{})[j].(map[string]interface{})["consumes"].([]interface{}) {
				methodParameters.MethodConsumes = append(methodParameters.MethodConsumes, l)
			}
			pathParam.Method = append(pathParam.Method, methodParameters)
			for _, w := range m["paths"].(map[string]interface{})[k].(map[string]interface{})[j].(map[string]interface{})["consumes"].([]interface{}) {
				methodParameters.MethodConsumes = append(methodParameters.MethodConsumes, w)
			}
		}

		api.PathParam = append(api.PathParam, pathParam)
	}
	//Проверка Авторизации
	if m["securityDefinitions"] != nil {
		for k := range m["securityDefinitions"].(map[string]interface{}) {
			api.Security = k
		}
	}

	//Ссылка "localhost"
	api.Host = m["host"].(string)

	//Базовый endpoint
	api.BasePath = m["basePath"].(string)
	if api.BasePath == "/" {
		api.BasePath = ""
	}



	//Генерация API
	f := NewFile("main")

	//func main
	f.Func().Id("main").Params().BlockFunc(func(group *Group) {
		group.Id("router").Op(":=").Qual("github.com/gin-gonic/gin", "Default").Call()
		for  i := 0; i < len(api.PathParam); i++ {
			group.Id("router").Dot(api.PathParam[i].Method[i].MethodType).Call(Lit(api.BasePath + api.PathParam[i].Path), Id("handlerAPI"))
		}
		group.Id("router").Dot("Run").Call(Lit(":80"))
	})

	//func handlerAPI
	f.Func().Id("handlerAPI").Params(Id("c").Add(Op("*")).Qual("github.com/gin-gonic/gin", "Context")).BlockFunc(func(group *Group) {
		group.Id("body").Op(":=").Id("c").Dot("Request").Dot("Body")
		group.Id("header").Op(":=").Id("c").Dot("Request").Dot("Header")
		group.Id("method").Op(":=").Id("c").Dot("Request").Dot("Method")
		group.Id("endpoint").Op(":=").Id("c").Dot("Request").Dot("RequestURI")

		group.Id("timeout").Op(":=").Qual("time", "Duration").Call(Id("10").Op("*").Qual("time", "Second"))
		group.Id("client").Op(":=").Qual("net/http", "Client").Values(Dict {
			Id("Timeout"): Id("timeout"),
		})
		group.Defer().Id("body").Dot("Close").Call()
		group.List(Id("request"), Id("err")).Op(":=").Qual("net/http", "NewRequest").Call(Id("method"),Add(Lit("http://"+ api.Host+api.BasePath)).Add(Op("+")).Id("endpoint"), Id("body"))
		group.If(
			Id("err").Op("!=").Id("nil").Block(
				Id("c").Dot("JSON").Call(List(Id("404"), Qual("github.com/gin-gonic/gin", "H").Values(Dict{
					Lit("error"): Id("err"),
				}))),
				Return(),
			),
		)
		group.Id("request").Dot("Header").Op("=").Id("header")
		group.List(Id("response"), Id("err")).Op(":=").Id("client").Dot("Do").Call(Id("request"))
		group.If(
			Id("err").Op("!=").Id("nil").Block(
				Id("c").Dot("JSON").Call(List(Id("404"), Qual("github.com/gin-gonic/gin", "H").Values(Dict{
					Lit("error"): Id("err"),
				}))),
				Return(),
			),
		)
		group.Defer().Id("response").Dot("Body").Dot("Close").Call()

		group.List(Id("bodyResp"), Id("err")).Op(":=").Qual("io/ioutil", "ReadAll").Call(Id("response").Dot("Body"))
		group.If(
			Id("err").Op("!=").Id("nil").Block(
				Id("c").Dot("JSON").Call(List(Id("500"), Qual("github.com/gin-gonic/gin", "H").Values(Dict{
					Lit("error"): Id("err"),
				}))),
				Return(),
			),
		)
		group.Id("c").Dot("JSON").Call(List(Id("200"), Qual("github.com/gin-gonic/gin", "H").Values(Dict{
			Lit("result"): String().Call(Id("bodyResp")),
		})))
	})


	//Сохраняем API под названием generatedAPI
	err := f.Save("generatedAPI.go")
	if err != nil {
		log.Fatalln(err)
	}
}