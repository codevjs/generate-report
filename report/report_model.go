package report

// ReportRow holds the data for each row returned by the SQL query.
// Fields correspond to the aliases in the SQL query.
type ReportRow struct {
	No               int    // Corresponds to '@a:=@a+1 AS No'
	JenisProject     string // Maps to 'Jenis Project'
	NamaProject      string // Maps to 'Nama Project'
	ProjectObjective string // Maps to 'Project Objective'
	BisnisProses     string // Maps to 'Bisnis Proses'
	Leader           string // Maps to 'Leader'
	Member           string // Maps to 'Member'
	MulaiProject     string // Maps to 'Mulai Project'
	SelesaiProject   string // Maps to 'Selesai Project'
	UploadProject    string // Maps to 'Upload Project'
	SubHolding       string // Maps to 'Sub Holding'
	SBU              string // Maps to 'SBU'
	StatusApprove    string // Maps to 'Status Approve'
	StepProject      string // Maps to 'Step Project'
	NQIPotential     string // Maps to 'NQI Potential'
	NQIReal          string // Maps to 'NQI Real'
	StatusNQI        string // Maps to 'Status NQI'
}
