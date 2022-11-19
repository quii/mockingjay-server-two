package matching

import "github.com/google/uuid"

type ReportCRUD struct {
	matchReports map[uuid.UUID]Report
}

func NewReportCRUD() *ReportCRUD {
	return &ReportCRUD{matchReports: make(map[uuid.UUID]Report)}
}

func (r ReportCRUD) GetAll() ([]Report, error) {
	var reports Reports
	for _, report := range r.matchReports {
		reports = append(reports, report)
	}
	reports.Sort()
	return reports, nil
}

func (r ReportCRUD) GetByID(id uuid.UUID) (Report, bool, error) {
	report, exists := r.matchReports[id]
	return report, exists, nil
}

func (r ReportCRUD) Create(t Report) error {
	r.matchReports[t.ID] = t
	return nil
}

func (r ReportCRUD) Delete(id uuid.UUID) error {
	delete(r.matchReports, id)
	return nil
}
