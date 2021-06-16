package main

import "html/template"

const (
	pageTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <title>Server {{.Version}}</title>
		{{if .Refresh}}<meta http-equiv="refresh" content="5;url=/">{{end}}
    </head>
    <body>
		{{if .Refresh}}
		<h1>Update is in progress, this page will reload in 5 seconds...</h1>
		{{else}}
        <h1>This server is version {{.Version}}</h1>
        <a href="/check">Check for new version</a>
        <br>
        {{if .NewVersion}}New version is available: {{.NewVersion.LastVer}} | <a href="/install">Upgrade</a>{{end}}
		{{end}}
    </body>
</html>
`
)

var PageTemplate = mustNewTemplate(template.New("page").Parse(pageTemplate))
