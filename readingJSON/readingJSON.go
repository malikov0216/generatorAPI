package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"io/ioutil"
	"strings"
)

type APIparameters struct {
	PathParam []pathParameters
	Host string
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
var ifMethodType MethodParameters
func main () {
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
	//Авторизация
	for k := range m["securityDefinitions"].(map[string]interface{}) {
		api.Security = k
	}
	//Порт
	api.Host = m["host"].(string)[9:]


	// MAKING API.go

	importLib := map[string]string{
		"github.com/gin-gonic/gin": "gin",
		"io/ioutil": "ioutil",
		"log": "log",
		"net/http": "net/http",
	}

	f := NewFile("main")
	f.ImportNames(importLib)

	f.Func().Id("main").Params().BlockFunc(func(group *Group) {
		group.Id("router").Op(":=").Qual("github.com/gin-gonic/gin", "Default").Call()
		for  i := 0; i < len(api.PathParam); i++ {
			group.Id("router").Dot(api.PathParam[i].Method[i].MethodType).Call(Lit(api.PathParam[i].Path), Id("handlerAPI"))
		}
		group.Id("router").Dot("Run").Call(Lit(":80"))
	})
	fmt.Println(pathParam.Method[0].MethodType)
	f.Func().Id("handlerAPI").Params(Id("c").Add(Op("*")).Qual("github.com/gin-gonic/gin", "Context")).BlockFunc(func(group *Group) {
		group.Id("body").Op(":=").Id("c").Dot("Request").Dot("Body")
		group.Id("header").Op(":=").Id("c").Dot("Request").Dot("Header")
		group.Id("method").Op(":=").Id("c").Dot("Request").Dot("Method")
		group.Id("endpoint").Op(":=").Id("c").Dot("Request").Dot("RequestURI")

		group.Id("timeout").Op(":=").Id("time").Dot("Duration").Call(Id("10").Op("*").Id("time").Dot("Second"))
		group.Id("client").Op(":=").Id("http").Dot("Client").Values(Dict {
			Id("Timeout"): Id("timeout"),
		})
		group.Defer().Id("body").Dot("Close").Call()
		group.List(Id("request"), Id("err")).Op(":=").Id("http").Dot("NewRequest").Call(Id("method"),Add(Lit("http://localhost"+api.Host)).Add(Op("+")).Id("endpoint"), Id("body"))
		group.If(
			Id("err").Op("!=").Id("nil").Block(
				Id("log").Dot("Fatal").Call(Id("err")),
			),
		)
		group.Id("request").Dot("Header").Op("=").Id("header")
		group.List(Id("response"), Id("err")).Op(":=").Id("client").Dot("Do").Call(Id("request"))
		group.If(
			Id("err").Op("!=").Id("nil").Block(
				Id("log").Dot("Fatal").Call(Id("err")),
			),
		)
		group.Defer().Id("response").Dot("Body").Dot("Close").Call()

		group.List(Id("bodyResp"), Id("err")).Op(":=").Id("ioutil").Dot("ReadAll").Call(Id("response").Dot("Body"))
		group.If(
			Id("err").Op("!=").Id("nil").Block(
				Id("log").Dot("Fatal").Call(Id("err")),
			),
		)
		group.Id("log").Dot("Println").Call(String().Call(Id("bodyResp")))
	})



	err := f.Save("main.go")
	if err != nil {
		fmt.Println(err)
		Continue()
	}
}