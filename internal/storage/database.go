package storage

func (s *Storage) GetShortUrlFromDB(id uint32, short *string) error {
	query := "SELECT short_url FROM short_url WHERE id = $1"
	return s.db.Get(short, query, id)
}

func (s *Storage) InsertShortUrlIntoDB(id uint32, url string, short string) error {
	query1 := "INSERT INTO short_url (id, short_url) VALUES ($1, $2)"
	query2 := "INSERT INTO long_url (id, long_url) VALUES ($1, $2)"

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(query1, id, short)
	_, err = tx.Exec(query2, id, url)

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (s *Storage) GetLongUrlFromDB(short string, long *string) error {
	query := "SELECT long_url FROM long_url JOIN short_url ON long_url.id = short_url.id WHERE short_url.short_url = $1;"
	return s.db.Get(long, query, short)
}

func (s *Storage) GetLongUrlByID(id uint32, long *string) error {
	query := "SELECT long_url FROM long_url WHERE id = $1;"
	return s.db.Get(long, query, id)
}
