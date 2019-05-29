package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"io/ioutil"
	"log"
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
	}

	f := NewFile("main")
	f.ImportNames(importLib)

	f.Func().Id("main").Params().BlockFunc(func(group *Group) {
		group.Id("router").Op(":=").Qual("github.com/gin-gonic/gin", "Default").Call()
		for  i := 0; i < len(api.PathParam); i++ {
			group.Id("router").Dot(api.PathParam[i].Method[i].MethodType).Call(Lit(api.PathParam[i].Path), Id("handlerAPI"))
		}
		group.Id("router").Dot("Run").Call(Lit(":1234"))
	})

	f.Func().Id("handlerAPI").Params(Id("c").Add(Op("*")).Qual("github.com/gin-gonic/gin", "Context")).BlockFunc(func(group *Group) {
		for k := range pathParam.Method {
			if pathParam.Method[k].MethodType == "post" {
				for i := range pathParam.Method[k].MethodConsumes {
					if pathParam.Method[k].MethodConsumes[i] == "multipart/form-data" {

					}
				}
			} else if pathParam.Method[k].MethodType == "GET" {
				fmt.Println("GET", pathParam.Method[k].MethodConsumes)
			}
		}
	})



	err := f.Save("main.go")
	if err != nil {
		log.Fatal("Eror on saving", err)
	}
}