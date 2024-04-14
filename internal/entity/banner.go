package entity

type Banner struct {
	ID              int
	JSONStructure   string
	FeatureID       int
	TagIDs          []int
	IsActive        bool
	UseLastRevision bool
}
