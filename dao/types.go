package dao

import "database/sql"

type Resolution struct {
	Width  int
	Height int
}

type SourceAsset struct {
	Id             int64
	Filename       string
	Codec          string
	Filesize       int
	DurationTime   float32
	DurationFrames int
	Fps            float32
	Resolution     Resolution

	daoInstance *DaoInstance
}

type DaoInstance struct {
	db *sql.DB
}

func NewDaoInstance(db *sql.DB) *DaoInstance {
	return &DaoInstance{db: db}
}
