package main

import (
	"log"
	"net/http"

	"github.com/BrownNPC/material-web-go-app/md"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Example struct {
	app.Compo
}

func (h *Example) Render() app.UI {
	return app.Div().
		//center using grid layout
		Body(

			md.LinearProgress().
				Attr("indeterminate", true),
			md.Icon().Text("settings"),
			app.H1().Text("Hello World").Class("md-typescale-display-large"),
			md.FilledTextField().Attr("value", "Hello World"),
			md.FilledIconButton().
				Body(md.Icon().
					Text("settings"),
				),
			md.Slider(),
			md.ChipSet().
				Body(
					md.AssistChip().Text("assist chip"),
					md.FilterChip().Text("filter chip"),
				),
			// DataTable Example
			app.Div().Style("margin-top", "2rem").Body(
				app.H2().Text("Data Table Example").Class("md-typescale-headline-medium"),
				md.NewDataTable(md.DataTableConfig{
					Selectable:  true,
					Sortable:    true,
					Pagination:  true,
					RowsPerPage: 5,
					Elevation:   1,
					Editable:    true, // 启用编辑功能
				}).
					SetColumns([]md.DataTableColumn{
						{Header: "Name", Field: "name", Sortable: true, Width: "200px"},
						{Header: "Age", Field: "age", Numeric: true, Sortable: true, Width: "100px"},
						{Header: "Email", Field: "email", Width: "300px"},
					}).
					SetData([]map[string]interface{}{
						{"name": "John Doe", "age": 28, "email": "john@example.com"},
						{"name": "Jane Smith", "age": 32, "email": "jane@example.com"},
						{"name": "Bob Johnson", "age": 45, "email": "bob@example.com"},
						{"name": "Alice Brown", "age": 24, "email": "alice@example.com"},
						{"name": "Charlie Wilson", "age": 36, "email": "charlie@example.com"},
						{"name": "Diana Miller", "age": 29, "email": "diana@example.com"},
					}),
			),
		)
}
func main() {
	handler := &app.Handler{
		Name:       "Material Design App",
		RawHeaders: append(md.RawHeaders(), md.LinkIconFontOutlined()),
	}
	app.Route("/", func() app.Composer { return &Example{} })
	app.RunWhenOnBrowser()
	log.Fatal(http.ListenAndServe(":8000", handler))
}
