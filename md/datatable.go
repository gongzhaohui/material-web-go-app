package md

import (
	"fmt"
	"sort"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// DataTableConfig 配置DataTable的行为和外观
type DataTableConfig struct {
	Selectable  bool
	Sortable    bool
	Pagination  bool
	RowsPerPage int
	Elevation   int
	Outlined    bool
	Editable    bool // 是否启用单元格编辑
}

// DataTableColumn 定义表格列
type DataTableColumn struct {
	Header   string
	Field    string
	Numeric  bool
	Sortable bool
	Width    string
	Align    string
}

// DataTable 组件状态
type DataTable struct {
	app.Compo
	config        DataTableConfig
	columns       []DataTableColumn
	data          []map[string]interface{}
	sortBy        string
	sortAscending bool
	selectedRows  map[int]bool
	currentPage   int
	editingCell   struct { // 当前编辑的单元格
		rowIndex int
		colField string
	}
}

// NewDataTable 创建新的DataTable实例
func NewDataTable(config DataTableConfig) *DataTable {
	return &DataTable{
		config:       config,
		selectedRows: make(map[int]bool),
		currentPage:  1,
	}
}

// SetColumns 设置表格列
func (dt *DataTable) SetColumns(columns []DataTableColumn) *DataTable {
	dt.columns = columns
	return dt
}

// SetData 设置表格数据
func (dt *DataTable) SetData(data []map[string]interface{}) *DataTable {
	dt.data = data
	return dt
}

// Render 渲染DataTable组件
func (dt *DataTable) Render() app.UI {
	if len(dt.data) == 0 || len(dt.columns) == 0 {
		return app.Div().Text("No data available")
	}

	// 排序数据
	sortedData := dt.sortData()

	// 分页数据
	paginatedData := dt.paginateData(sortedData)

	return app.Div().
		Class("mdc-data-table").
		Style("overflow-x", "auto").
		Body(
			// 表格容器
			app.Div().
				Class("mdc-data-table__table-container").
				Body(
					// 表格
					app.Table().
						Class("mdc-data-table__table").
						Aria("label", "Data Table").
						Body(
							// 表头
							app.THead().Body(dt.renderHeader()),
							// 表格主体
							app.TBody().
								Class("mdc-data-table__content").
								Body(dt.renderRows(paginatedData)...),
						),
				),
			// 分页控件
			dt.renderPagination(len(sortedData)),
		)
}

// renderHeader 渲染表头
func (dt *DataTable) renderHeader() app.UI {
	return app.Tr().
		Class("mdc-data-table__header-row").
		Body(dt.renderHeaderCells()...)
}

// renderHeaderCells 渲染表头单元格
func (dt *DataTable) renderHeaderCells() []app.UI {
	var cells []app.UI

	// 选择列
	if dt.config.Selectable {
		cells = append(cells,
			app.Th().
				Class("mdc-data-table__header-cell", "mdc-data-table__header-cell--checkbox").
				Role("columnheader").
				Scope("col").
				Body(
					app.Div().
						Class("mdc-data-table__header-cell-wrapper").
						Body(
							app.Input().
								Type("checkbox").
								Class("mdc-checkbox__native-control").
								Attr("aria-label", "Select all rows").
								OnChange(dt.onSelectAll),
						),
				),
		)
	}

	// 数据列
	for _, col := range dt.columns {
		cells = append(cells, dt.renderHeaderCell(col))
	}

	return cells
}

// renderHeaderCell 渲染单个表头单元格
func (dt *DataTable) renderHeaderCell(col DataTableColumn) app.UI {
	cell := app.Th().
		Class("mdc-data-table__header-cell").
		Role("columnheader").
		Scope("col").
		Style("width", col.Width)

	if col.Numeric {
		cell = cell.Class("mdc-data-table__header-cell--numeric")
	}

	content := app.Div().
		Class("mdc-data-table__header-cell-wrapper")

	if col.Sortable && dt.config.Sortable {
		content = content.Body(
			app.Button().
				Class("mdc-data-table__sort-indicator").
				OnClick(func(ctx app.Context, e app.Event) {
					dt.onSort(col.Field)
				}).
				Body(
					app.Span().Text(col.Header),
					dt.renderSortIndicator(col.Field),
				),
		)
	} else {
		content = content.Body(app.Span().Text(col.Header))
	}

	return cell.Body(content)
}

// renderSortIndicator 渲染排序指示器
func (dt *DataTable) renderSortIndicator(field string) app.UI {
	if dt.sortBy != field {
		return app.Span().Class("mdc-data-table__sort-indicator")
	}

	if dt.sortAscending {
		return app.Span().Class("mdc-data-table__sort-indicator", "mdc-data-table__sort-indicator--ascending")
	}
	return app.Span().Class("mdc-data-table__sort-indicator", "mdc-data-table__sort-indicator--descending")
}

// renderRows 渲染表格行
func (dt *DataTable) renderRows(data []map[string]interface{}) []app.UI {
	var rows []app.UI
	startIndex := (dt.currentPage - 1) * dt.config.RowsPerPage

	for i, rowData := range data {
		globalIndex := startIndex + i
		rows = append(rows, dt.renderRow(globalIndex, rowData))
	}
	return rows
}

// renderRow 渲染单行
func (dt *DataTable) renderRow(index int, rowData map[string]interface{}) app.UI {
	row := app.Tr().
		Class("mdc-data-table__row").
		Body(dt.renderRowCells(index, rowData)...)

	if dt.selectedRows[index] {
		row = row.Class("mdc-data-table__row--selected")
	}

	return row
}

// renderRowCells 渲染行单元格
func (dt *DataTable) renderRowCells(index int, rowData map[string]interface{}) []app.UI {
	var cells []app.UI

	// 选择单元格
	if dt.config.Selectable {
		cells = append(cells,
			app.Td().
				Class("mdc-data-table__cell", "mdc-data-table__cell--checkbox").
				Body(
					app.Input().
						Type("checkbox").
						Class("mdc-checkbox__native-control").
						Checked(dt.selectedRows[index]).
						OnChange(func(ctx app.Context, e app.Event) {
							dt.onRowSelect(index)
						}),
				),
		)
	}

	// 数据单元格
	for _, col := range dt.columns {
		value := rowData[col.Field]
		cells = append(cells, dt.renderCell(index, value, col))
	}

	return cells
}

// renderCell 渲染单个单元格
func (dt *DataTable) renderCell(rowIndex int, value interface{}, col DataTableColumn) app.UI {
	cellClass := "mdc-data-table__cell"
	if dt.config.Editable {
		cellClass += " mdc-data-table__editable-cell"
	}

	cell := app.Td().
		Class(cellClass).
		OnClick(func(ctx app.Context, e app.Event) {
			if dt.config.Editable {
				dt.editingCell.rowIndex = rowIndex
				dt.editingCell.colField = col.Field
				dt.Update()
			}
		}).
		OnDblClick(func(ctx app.Context, e app.Event) {
			if dt.config.Editable {
				dt.editingCell.rowIndex = rowIndex
				dt.editingCell.colField = col.Field
				dt.Update()
			}
		})

	if col.Numeric {
		cell = cell.Class("mdc-data-table__cell--numeric")
	}

	// 如果是正在编辑的单元格
	if dt.config.Editable && dt.editingCell.rowIndex == rowIndex && dt.editingCell.colField == col.Field {
		return cell.Body(
			app.Input().
				Type("text").
				Value(formatValue(value)).
				Class("mdc-data-table__edit-input").
				OnBlur(func(ctx app.Context, e app.Event) {
					// 保存编辑内容
					newValue := e.Get("target").Get("value").String()
					dt.data[rowIndex][col.Field] = newValue
					dt.editingCell.rowIndex = -1 // 结束编辑
					dt.Update()
				}).
				OnKeyDown(func(ctx app.Context, e app.Event) {
					if e.Get("key").String() == "Enter" {
						// 回车保存
						newValue := e.Get("target").Get("value").String()
						dt.data[rowIndex][col.Field] = newValue
						dt.editingCell.rowIndex = -1
						dt.Update()
					} else if e.Get("key").String() == "Escape" {
						// ESC取消
						dt.editingCell.rowIndex = -1
						dt.Update()
					}
				}),
		)
	}

	return cell.Text(formatValue(value))
}

// renderPagination 渲染分页控件
func (dt *DataTable) renderPagination(totalRows int) app.UI {
	if !dt.config.Pagination {
		return app.Div()
	}

	totalPages := (totalRows + dt.config.RowsPerPage - 1) / dt.config.RowsPerPage
	if totalPages <= 1 {
		return app.Div()
	}

	return app.Div().
		Class("mdc-data-table__pagination").
		Body(
			app.Div().
				Class("mdc-data-table__pagination-navigation").
				Body(
					dt.renderPaginationButton("first", "First", dt.currentPage == 1),
					dt.renderPaginationButton("prev", "Previous", dt.currentPage == 1),
					app.Span().
						Class("mdc-data-table__pagination-info").
						Textf("Page %d of %d", dt.currentPage, totalPages),
					dt.renderPaginationButton("next", "Next", dt.currentPage == totalPages),
					dt.renderPaginationButton("last", "Last", dt.currentPage == totalPages),
				),
		)
}

// renderPaginationButton 渲染分页按钮
func (dt *DataTable) renderPaginationButton(action, label string, disabled bool) app.UI {
	button := app.Button().
		Class("mdc-button", "mdc-data-table__pagination-button").
		Attr("aria-label", label).
		OnClick(func(ctx app.Context, e app.Event) {
			dt.onPaginate(action)
		})

	if disabled {
		button = button.Disabled(true)
	}

	return button.Text(label)
}

// sortData 排序数据
func (dt *DataTable) sortData() []map[string]interface{} {
	if dt.sortBy == "" {
		return dt.data
	}

	sorted := make([]map[string]interface{}, len(dt.data))
	copy(sorted, dt.data)

	sort.Slice(sorted, func(i, j int) bool {
		val1 := sorted[i][dt.sortBy]
		val2 := sorted[j][dt.sortBy]

		if dt.sortAscending {
			return compareValues(val1, val2)
		}
		return compareValues(val2, val1)
	})

	return sorted
}

// paginateData 分页数据
func (dt *DataTable) paginateData(data []map[string]interface{}) []map[string]interface{} {
	if !dt.config.Pagination {
		return data
	}

	start := (dt.currentPage - 1) * dt.config.RowsPerPage
	end := start + dt.config.RowsPerPage
	if end > len(data) {
		end = len(data)
	}

	return data[start:end]
}

// 事件处理函数
func (dt *DataTable) onSort(field string) {
	if dt.sortBy == field {
		dt.sortAscending = !dt.sortAscending
	} else {
		dt.sortBy = field
		dt.sortAscending = true
	}
}

func (dt *DataTable) onRowSelect(index int) {
	if dt.selectedRows[index] {
		delete(dt.selectedRows, index)
	} else {
		dt.selectedRows[index] = true
	}
}

func (dt *DataTable) onSelectAll(ctx app.Context, e app.Event) {
	selectAll := e.Get("target").Get("checked").Bool()
	dt.selectedRows = make(map[int]bool)

	if selectAll {
		for i := range dt.data {
			dt.selectedRows[i] = true
		}
	}
}

func (dt *DataTable) onPaginate(action string) {
	totalPages := (len(dt.data) + dt.config.RowsPerPage - 1) / dt.config.RowsPerPage

	switch action {
	case "first":
		dt.currentPage = 1
	case "prev":
		if dt.currentPage > 1 {
			dt.currentPage--
		}
	case "next":
		if dt.currentPage < totalPages {
			dt.currentPage++
		}
	case "last":
		dt.currentPage = totalPages
	}
}

// Update triggers a UI update for the component; this is a no-op placeholder
// so that calls to dt.Update() compile — replace with actual update logic if
// your go-app version provides a different mechanism to request a re-render.
// Update 触发组件重新渲染
func (dt *DataTable) Update() {
	dt.Render()
}

// 辅助函数
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%.2f", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func compareValues(a, b interface{}) bool {
	switch aVal := a.(type) {
	case string:
		bVal := b.(string)
		return aVal < bVal
	case int:
		bVal := b.(int)
		return aVal < bVal
	case float64:
		bVal := b.(float64)
		return aVal < bVal
	case float32:
		bVal := b.(float32)
		return aVal < bVal
	default:
		return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
	}
}
