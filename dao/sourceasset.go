package dao

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
)

func (d *DaoInstance) NewSourceAsset(filename string) (*SourceAsset, error) {
	result, err := d.db.Exec("INSERT INTO source_asset (filename) values (?);", filename)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &SourceAsset{
		Id:          id,
		Filename:    filename,
		daoInstance: d,
	}, nil
}

func (d *DaoInstance) GetSourceAssets(where string, args ...any) ([]*SourceAsset, error) {
	var assets []*SourceAsset

	query := "SELECT id,filename,codec,filesize,duration_time,duration_frames,fps,resolution_width,resolution_height from source_asset " + where
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return assets, fmt.Errorf("unable to get assets: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		sa := &SourceAsset{daoInstance: d}
		err := rows.Scan(&sa.Id, &sa.Filename, &sa.Codec, &sa.Filesize, &sa.DurationTime, &sa.DurationFrames, &sa.Fps, &sa.Resolution.Width, &sa.Resolution.Height)
		if err != nil {
			log.Errorf("unable to scan source_asset row: %s", err)
		}
		assets = append(assets, sa)
	}

	return assets, nil
}

func (d *DaoInstance) GetSourceAssetByFilename(filename string) (*SourceAsset, error) {
	sa := &SourceAsset{daoInstance: d}
	query := "SELECT id,filename,codec,filesize,duration_time,duration_frames,fps,resolution_width,resolution_height from source_asset where filename = :filename"
	row := d.db.QueryRow(query, sql.Named("filename", filename))
	err := row.Scan(&sa.Id, &sa.Filename, &sa.Codec, &sa.Filesize, &sa.DurationTime, &sa.DurationFrames, &sa.Fps, &sa.Resolution.Width, &sa.Resolution.Height)
	if err != nil {
		return nil, err
	}
	return sa, nil
}

func (s *SourceAsset) Save() error {
	if s.daoInstance == nil {
		return fmt.Errorf("no daoInstance for SourceAsset")
	}
	query := "UPDATE source_asset set codec = :codec, filesize = :filesize, duration_time = :duration_time, duration_frames = :duration_frames, fps = :fps, resolution_width = :width, resolution_height = :height where id = :id;"
	_, err := s.daoInstance.db.Exec(query,
		sql.Named("codec", s.Codec),
		sql.Named("filesize", s.Filesize),
		sql.Named("duration_time", s.DurationTime),
		sql.Named("duration_frames", s.DurationFrames),
		sql.Named("fps", s.Fps),
		sql.Named("width", s.Resolution.Width),
		sql.Named("height", s.Resolution.Height),
		sql.Named("id", s.Id),
	)
	if err != nil {
		return fmt.Errorf("unable to update source asset %s: %s", s.Filename, err)
	}
	return nil
}
