package dao

func (d *DaoInstance) NewTranscodeOutput(filename string, filesize int, source *SourceAsset, profile string, resolution Resolution) (*TranscodeOutput, error) {
	result, err := d.db.Exec(
		"INSERT INTO transcode_asset (source, filename, filesize, profile_name, resolution_width, resolution_height) VALUES (?,?,?,?,?,?);",
		source.Id, filename, filesize, profile, resolution.Width, resolution.Height)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &TranscodeOutput{
		Id:          id,
		Filename:    filename,
		Source:      source.Id,
		Filesize:    filesize,
		ProfileName: &profile,
		Resolution:  resolution,
	}, nil
}

func (d *DaoInstance) GetTranscodeOutputs(where string, args ...any) (outputs []TranscodeOutput, err error) {
	query := "SELECT id, source, filename, filesize, profile_name, resolution_width, resolution_height from transcode_asset " + where
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return outputs, err
	}

	defer rows.Close()

	for rows.Next() {
		to := TranscodeOutput{DaoInstance: d}
		err = rows.Scan(&to.Id, &to.Source, &to.Filename, &to.Filesize, &to.ProfileName, &to.Resolution.Width, &to.Resolution.Height)
		if err != nil {
			return outputs, err
		}
		outputs = append(outputs, to)
	}

	return outputs, err
}
