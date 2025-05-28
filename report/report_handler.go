package report

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"report-service/report" // For ReportRow struct

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	// _ "github.com/go-sql-driver/mysql" // Driver registered in main.go or database.go
)

// ExportReportHandler creates a gin.HandlerFunc for exporting reports.
func ExportReportHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 3.a. Define the main SQL query string
		// SQL query adapted from the prompt for this specific subtask.
		// Removed DATE_FORMAT from date fields as per current prompt's SQL.
		sqlQuery := `
SELECT
    @a:=@a+1 AS No,
    g.category_name AS 'Jenis Project',
    a.project_name AS 'Nama Project',
    a.project_objective AS 'Project Objective',
    d.name AS 'Bisnis Proses',
    b.fullname AS 'Leader',
    i.Member,
    a.project_start AS 'Mulai Project',
    a.project_finish AS 'Selesai Project',
    a.project_created_date AS 'Upload Project',
    f.subholding_name AS 'Sub Holding',
    c.corporate_name AS 'SBU',
    CASE a.project_status
        WHEN "New" THEN "Waiting Approve"
        WHEN "Ongoing" THEN "Approve"
        WHEN "Reject" THEN "Reject"
        WHEN "Done" THEN "Done"
        WHEN "Drop" THEN "Delete"
        ELSE ""
    END AS 'Status Approve',
    h.step_project AS 'Step Project',
    a.project_nqi_potential AS 'NQI Potential',
    a.project_nqi_real AS 'NQI Real',
    CASE a.project_approve_nqi
        WHEN "None" THEN "Belum Diajukan"
        WHEN "Waiting" THEN "Menunggu Persetujuan"
        WHEN "Approved" THEN "Disetujui"
        WHEN "Reject" THEN "Ditolak"
        ELSE ""
    END AS 'Status NQI'
FROM kaizen_list_project a
LEFT JOIN users b ON b.username = a.project_pic_nik
LEFT JOIN corporate c ON b.corporate_id = c.corporate_id
LEFT JOIN kaizen_ref_bisnis_proses d ON a.project_bispro_id = d.id
LEFT JOIN departments e ON b.department_id = e.department_id
LEFT JOIN subholding f ON b.subholding_id = f.subholding_id
LEFT JOIN kaizen_ref_category g ON a.project_category = g.category_id
LEFT JOIN (
    SELECT id_project, MAX(step_project) AS step_project
    FROM kaizen_list_file
    GROUP BY id_project
) h ON a.project_id = h.id_project
LEFT JOIN (
    SELECT klm.member_project_id, GROUP_CONCAT(u.fullname SEPARATOR ', ') AS Member
    FROM kaizen_list_member klm -- Corrected alias from 'a' to 'klm'
    JOIN users u ON klm.member_id = u.username -- Corrected alias from 'b' to 'u'
    GROUP BY klm.member_project_id
) i ON a.project_id = i.member_project_id
WHERE
    b.corporate_id = 'dc6529e8-7c9b-48a4-9ceb-a6df1dae7c78' AND
    DATE_FORMAT(a.project_created_date, '%Y-%m-%d') >= '2025-05-04' AND
    DATE_FORMAT(a.project_created_date, '%Y-%m-%d') <= '2025-05-06'
ORDER BY a.project_created_date ASC;
`
		// 3.b. Execute SET @a := 0;
		_, err := db.ExecContext(c.Request.Context(), "SET @a := 0")
		if err != nil {
			log.Printf("Error setting session variable @a: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize report query"})
			return
		}

		// 3.c. Execute the main SELECT query
		rows, err := db.QueryContext(c.Request.Context(), sqlQuery)
		if err != nil {
			log.Printf("Error executing main report query: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve report data"})
			return
		}
		defer rows.Close()

		// 3.d. Create a slice results := []report.ReportRow{}
		results := []report.ReportRow{}

		// 3.e. Iterate through the query results
		for rows.Next() {
			var r report.ReportRow
			var member, stepProject, nqiPotential, nqiReal sql.NullString
			// For date fields, if not formatted in SQL, they might be time.Time or []byte.
			// Scanning directly into string might work or might need sql.NullTime and then formatting.
			// Assuming for now that the driver handles conversion of date/time types to string appropriately
			// when the scan destination is a string field. If not, this will need sql.NullTime.
			var mulaiProject, selesaiProject, uploadProject sql.NullString


			err := rows.Scan(
				&r.No,
				&r.JenisProject,
				&r.NamaProject,
				&r.ProjectObjective,
				&r.BisnisProses,
				&r.Leader,
				&member,
				&mulaiProject,    // Changed to scan into sql.NullString
				&selesaiProject, // Changed to scan into sql.NullString
				&uploadProject,  // Changed to scan into sql.NullString
				&r.SubHolding,
				&r.SBU,
				&r.StatusApprove,
				&stepProject,
				&nqiPotential,
				&nqiReal,
				&r.StatusNQI,
			)
			if err != nil { // 3.e.iv. Handle errors from rows.Scan()
				log.Printf("Error scanning row data: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing report data"})
				return
			}

			r.Member = ""
			if member.Valid { r.Member = member.String }
			
			r.MulaiProject = ""
			if mulaiProject.Valid { r.MulaiProject = mulaiProject.String }
			
			r.SelesaiProject = ""
			if selesaiProject.Valid { r.SelesaiProject = selesaiProject.String }
			
			r.UploadProject = ""
			if uploadProject.Valid { r.UploadProject = uploadProject.String }
			
			r.StepProject = ""
			if stepProject.Valid { r.StepProject = stepProject.String }
			
			r.NQIPotential = ""
			if nqiPotential.Valid { r.NQIPotential = nqiPotential.String }
			
			r.NQIReal = ""
			if nqiReal.Valid { r.NQIReal = nqiReal.String }
			
			results = append(results, r) // 3.e.iii. Append populated ReportRow
		}

		// 3.f. Check for errors after the loop
		if err = rows.Err(); err != nil {
			log.Printf("Error after iterating rows: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finalizing report data processing"})
			return
		}

		// 3.g. Define the path to the template
		templatePath := "templates/template.xlsx"

		// 3.h. Open the Excel template
		f, err := excelize.OpenFile(templatePath)
		if err != nil {
			log.Printf("Error opening template file '%s': %v", templatePath, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open report template"})
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("Error closing template file '%s': %v", templatePath, err)
			}
		}()

		// 3.i. Get the first sheet name
		sheetName := f.GetSheetName(0) // excelize v2 uses 0-based index
		if sheetName == "" {
			log.Printf("Sheet at index 0 not found in template '%s'. Trying active sheet.", templatePath)
			activeSheetIdx := f.GetActiveSheetIndex()
			sheetName = f.GetSheetName(activeSheetIdx)
			if sheetName == "" {
				log.Printf("Active sheet in '%s' also not found or has no name. Defaulting to 'Sheet1'.", templatePath)
				sheetName = "Sheet1"
			}
		}
		
		// 3.j. Define the starting row for data entry
		startRow := 4

		// 3.k. Iterate through the results and populate the sheet
		for index, rowData := range results {
			excelRowIndex := startRow + index

			values := []interface{}{
				rowData.No, rowData.JenisProject, rowData.NamaProject, rowData.ProjectObjective,
				rowData.BisnisProses, rowData.Leader, rowData.Member, rowData.MulaiProject,
				rowData.SelesaiProject, rowData.UploadProject, rowData.SubHolding, rowData.SBU,
				rowData.StatusApprove, rowData.StepProject, rowData.NQIPotential,
				rowData.NQIReal, rowData.StatusNQI,
			}
			
			 colNames := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q"}
			 for colIdx, val := range values {
				 cellName := fmt.Sprintf("%s%d", colNames[colIdx], excelRowIndex)
				 if err := f.SetCellValue(sheetName, cellName, val); err != nil {
					 log.Printf("Error setting cell value for %s on sheet '%s': %v", cellName, sheetName, err)
					 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to populate report data into Excel sheet"})
					 return
				 }
			 }
		}

		// 3.l. Set HTTP headers for Excel download
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=rekap_report.xlsx")

		// 3.m. Write the Excel file to the response
		if err := f.Write(c.Writer); err != nil {
			log.Printf("Error writing Excel file to HTTP response: %v", err)
		}
	}
}
