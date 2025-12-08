package md

import (
	_ "embed"
)

//go:embed datatable.css
var datatableCSS string

// <link href="web/roboto.css" rel="stylesheet">
var headers string = `
<link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&amp;display=swap" rel="stylesheet">
<script type="importmap">
  {
    "imports": {
      "@material/web/": "https://esm.run/@material/web/"
    }
  }
</script>
<script type="module">
  import '@material/web/all.js';
  import {styles as typescaleStyles} from '@material/web/typography/md-typescale-styles.js';

  document.adoptedStyleSheets.push(typescaleStyles.styleSheet);
</script>
<style>
` + datatableCSS + `
</style>
<script>
// 列宽拖拽功能
document.addEventListener('DOMContentLoaded', function() {
    const tables = document.querySelectorAll('.mdc-data-table');
    
    tables.forEach(table => {
        const headers = table.querySelectorAll('.mdc-data-table__header-cell');
        let isResizing = false;
        let currentHeader = null;
        let startX = 0;
        let startWidth = 0;

        headers.forEach(header => {
            // 为每个表头单元格添加拖拽事件
            header.addEventListener('mousedown', function(e) {
                // 检查是否点击在拖拽区域（右侧4px范围内）
                const rect = header.getBoundingClientRect();
                if (e.clientX > rect.right - 8) {
                    isResizing = true;
                    currentHeader = header;
                    startX = e.clientX;
                    startWidth = header.offsetWidth;
                    
                    // 添加拖拽样式
                    document.body.style.cursor = 'col-resize';
                    document.body.style.userSelect = 'none';
                    e.preventDefault();
                }
            });
        });

        document.addEventListener('mousemove', function(e) {
            if (!isResizing || !currentHeader) return;

            const width = startWidth + (e.clientX - startX);
            if (width > 50) { // 最小宽度限制
                currentHeader.style.width = width + 'px';
                
                // 更新对应的数据单元格
                const colIndex = Array.from(headers).indexOf(currentHeader);
                const colNum = colIndex + 1;
                const cells = table.querySelectorAll('.mdc-data-table__cell:nth-child(' + colNum + ')');
                cells.forEach(cell => {
                    cell.style.width = width + 'px';
                });
            }
        });

        document.addEventListener('mouseup', function() {
            if (isResizing) {
                isResizing = false;
                currentHeader = null;
                document.body.style.cursor = '';
                document.body.style.userSelect = '';
            }
        });
    });
});
</script>`

func RawHeaders() []string {
	return []string{headers}
}
