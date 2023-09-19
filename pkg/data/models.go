package data

import (
	"encoding/json"
	"github.com/google/uuid"
)

type PatientDataHeader struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	PatientId int
	DayNumber int
	Items     []PatientDataEntry `gorm:"foreignKey:PatientDataHeaderId"`
}

type PatientDataEntry struct {
	Id                  uuid.UUID `gorm:"primaryKey"`
	PatientDataHeaderId uuid.UUID `gorm:"index"`
	Offset              int
	Length              int
	Mean                int
}

type PatientDto struct {
	PatientId int
	DayNumber int
	Items     []PatientEntryDto
}

type PatientEntryDto struct {
	Offset int
	Length int
	Mean   int
}

func NewPatient(patientId int, dayNumber int) *PatientDto {
	return &PatientDto{
		PatientId: patientId,
		DayNumber: dayNumber,
		Items:     make([]PatientEntryDto, 0),
	}
}

func NewPatientFromBytes(bytes []byte) (*PatientDto, error) {
	var patient PatientDto
	err := json.Unmarshal(bytes, &patient)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (p *PatientDto) AddItem(offset int, length int, mean int) {
	p.Items = append(p.Items, PatientEntryDto{
		Offset: offset,
		Length: length,
		Mean:   mean,
	})
}
