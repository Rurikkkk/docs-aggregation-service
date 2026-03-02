package aggregate

import "time"

type Document struct {
	Doc DocumentFields `bson:"doc" json:"doc"`
}

type DocumentFields struct {
	DateTime             time.Time `bson:"dateTime" json:"dateTime"`
	FiscalDocumentNumber int64     `bson:"fiscalDocumentNumber" json:"fiscalDocumentNumber"`
	FiscalDriveNumber    string    `bson:"fiscalDriveNumber" json:"fiscalDriveNumber"`
	Items                []Item    `bson:"items" json:"items"`
}

type Item struct {
	Name string `bson:"name" json:"name"`
}
