package main

import (
	"github.com/astaxie/beego"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
	"net/http"
)

// MyController is a test controller.
type MyController struct {
	beego.Controller
}

// SchemaSelector is the interface that your schema selector should implements.
// type SchemaSelector interface {
// 	Schema(req *http.Request) gojsonschema.JSONLoader
// }

// SchemaSelector is the interface that your schema selector should implements.
type SchemaSelector func(req *http.Request) gojsonschema.JSONLoader

// DefaultSchemaSelector is the default schema loader for http.Request.
func DefaultSchemaSelector(req *http.Request) gojsonschema.JSONLoader {
	return gojsonschema.NewReferenceLoader("file:///home/baibin/workspace/go/src/error-response-400/schema/test1.json")
}

// Prepare override the default Prepare() and do JSON request schema validation.
func (c *MyController) Prepare() {
	log.Printf("Prepare()\n")
	ValidateRequest(&c.Controller, DefaultSchemaSelector)
}

// Post override the default Post().
func (c *MyController) Post() {
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.ServeJSON()
}

// ValidateRequest validate each incoming request.
func ValidateRequest(c *beego.Controller, selector SchemaSelector) {
	log.Printf("ValidateRequest()\n")
	schema := selector(c.Ctx.Request)
	b, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		return
	}
	result, err := gojsonschema.Validate(schema, gojsonschema.NewBytesLoader(b))
	if err != nil {
		log.Printf("JSON format error, err = %v\n", err)
		return
	}
	if result.Valid() {
		log.Printf("Request valid.\n")
		return
	}
	for _, err := range result.Errors() {
		// Err implements the ResultError interface
		log.Printf("- %s\n", err)
	}
	c.Ctx.Output.SetStatus(http.StatusBadRequest)
	c.ServeJSON()
}

func main() {
	beego.AddNamespace(
		beego.NewNamespace(
			"validator",
			beego.NSRouter("/test1", &MyController{}, "post:Post"),
		),
	)
	beego.BConfig.Listen.HTTPPort = 3000
	// beego.BConfig.CopyRequestBody = true
	beego.Run()
}
