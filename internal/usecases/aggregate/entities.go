package aggregate

import "time"

type Document struct {
	Doc DocumentFields `bson:"doc"`
}

type DocumentFields struct {
	DateTime             time.Time `bson:"dateTime"`
	FiscalDocumentNumber int64     `bson:"fiscalDocumentNumber"`
	FiscalDriveNumber    string    `bson:"fiscalDriveNumber"`
	Items                []Item    `bson:"items"`
}

type Item struct {
	Name string `bson:"name"`
}
