/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

var (
	routerTplName     = "router.go"
	middlewareTplName = "middleware.go"
	handlerTplName    = "handler.go"
	modelTplName      = "model.go"
	registerTplName   = "register.go"
	clientTplName     = "client.go"

	insertPointNew        = "//INSERT_POINT: DO NOT DELETE THIS LINE!"
	insertPointPatternNew = `//INSERT_POINT\: DO NOT DELETE THIS LINE\!`
	insertPointPatternOld = `func register\(r \*(hertz|server)\.Hertz\) \{\n`
)

func IsDefaultTpl(name string) bool {
	return name == routerTplName || name == handlerTplName || name == modelTplName || name == registerTplName || name == middlewareTplName || name == clientTplName
}

var defaultPkgConfig = TemplateConfig{
	Layouts: []Template{
		{
			Path:   defaultHandlerDir + sp + handlerTplName,
			Delims: [2]string{"{{", "}}"},
			Body: `// Code generated by hertz generator.

package {{.PackageName}}

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

{{- range $k, $v := .Imports}}
	{{$k}} "{{$v.Package}}"
{{- end}}
)

{{range $_, $MethodInfo := .Methods}}
{{$MethodInfo.Comment}}
func {{$MethodInfo.Name}}(ctx context.Context, c *app.RequestContext) { 
	var err error
	{{if ne $MethodInfo.RequestTypeName "" -}}
	var req {{$MethodInfo.RequestTypeName}}
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	{{end}}
	resp := new({{$MethodInfo.ReturnTypeName}})

	c.{{.Serializer}}(200, resp)
}
{{end}}
			`,
		},
		{
			Path:   defaultRouterDir + sp + routerTplName,
			Delims: [2]string{"{{", "}}"},
			Body: `// Code generated by hertz generator. DO NOT EDIT.

package {{$.PackageName}}

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	{{range $k, $v := .HandlerPackages}}{{$k}} "{{$v}}"{{end}}
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
 */

{{define "g"}}
{{- if eq .Path "/"}}r
{{- else}}{{.GroupName}}{{end}}
{{- end}}

{{define "G"}}
{{- if ne .Handler ""}}
	{{- .GroupName}}.{{.HttpMethod}}("{{.Path}}", append({{.MiddleWare}}Mw(), {{.Handler}})...)
{{- end}}
{{- if ne (len .Children) 0}}
{{.MiddleWare}} := {{template "g" .}}.Group("{{.Path}}", {{.MiddleWare}}Mw()...)
{{- end}}
{{- range $_, $router := .Children}}
{{- if ne .Handler ""}}
	{{template "G" $router}}
{{- else}}
	{	{{template "G" $router}}
	}
{{- end}}
{{- end}}
{{- end}}

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {
{{template "G" .Router}}
}

		`,
		},
		{
			Path: defaultRouterDir + sp + registerTplName,
			Body: `// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	{{$.PkgAlias}} "{{$.Pkg}}"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz){
	` + insertPointNew + `
	{{$.PkgAlias}}.Register(r)
}
`,
		},
		// Model tpl is imported by model generator. Here only decides model directory.
		{
			Path: defaultModelDir + sp + modelTplName,
			Body: ``,
		},
		{
			Path:   defaultRouterDir + sp + middlewareTplName,
			Delims: [2]string{"{{", "}}"},
			Body: `// Code generated by hertz generator.

package {{$.PackageName}}

import (
	"github.com/cloudwego/hertz/pkg/app"
)

{{define "M"}}
func {{.MiddleWare}}Mw() []app.HandlerFunc {
	// your code...
	return nil
}
{{range $_, $router := $.Children}}{{template "M" $router}}{{end}}
{{- end}}

{{template "M" .Router}}

		`,
		},
		{
			Path:   defaultClientDir + sp + clientTplName,
			Delims: [2]string{"{{", "}}"},
			Body: `// Code generated by hertz generator.

package {{$.PackageName}}

import (
    "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/config"
)

type {{.ServiceName}}Client struct {
	client * client.Client
}

func New{{.ServiceName}}Client(opt ...config.ClientOption) (*{{.ServiceName}}Client, error) {
	c, err := client.NewClient(opt...)
	if err != nil {
		return nil, err
	}

	return &{{.ServiceName}}Client{
		client: c,
	}, nil
}
		`,
		},
	},
}
