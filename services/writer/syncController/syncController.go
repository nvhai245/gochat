func (s *server) Write(stream pb.Sync_WriteServer) error {
	var allMessages = []*pb.WriteRequest{}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		allMessages = append(allMessages, message)
	}
	for _, message := range allMessages {
		quoted := pq.QuoteIdentifier(message.Table)
		sqlStatement := fmt.Sprintf(`
	INSERT INTO %s (count, author, message, deleted)
	VALUES ($1, $2, $3, $4)`, quoted)
		rows, err := db.Query(sqlStatement, message.Count, message.Author, message.Message, false)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close();
	}
	return stream.SendAndClose(&pb.WriteResponse{
		Success: true,
	})
}
func (s *server) GetDifference(ctx context.Context, localCount *pb.GetDifferenceRequest) (*pb.GetDifferenceResponse, error) {
	quoted := pq.QuoteIdentifier(localCount.Table)
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s ORDER BY count DESC LIMIT 1`, quoted)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var (
		count   int64
		author  string
		message string
		deleted bool
	)

	for rows.Next() {
		if err := rows.Scan(&count, &author, &message, &deleted); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	difference := count
	return &pb.GetDifferenceResponse{Difference: difference}, nil
}
func (s *server) Read(messagesRange *pb.ReadRequest, stream pb.Sync_ReadServer) error {
	for i := messagesRange.First; i <= messagesRange.Last; i++ {
		quoted := pq.QuoteIdentifier(messagesRange.Table)
		sqlStatement := fmt.Sprintf(`SELECT count, author, message, deleted FROM %s WHERE count=$1`, quoted)
		rows, err := db.Query(sqlStatement, i)
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()
		message := &pb.ReadResponse{}
		for rows.Next() {
			if err := rows.Scan(&message.Count, &message.Author, &message.Message, &message.Deleted); err != nil {
				log.Fatal(err)
			}
		}
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) CheckExist(ctx context.Context, tables *pb.CheckRequest) (*pb.CheckResponse, error) {
	quoted1 := pq.QuoteIdentifier(tables.Table1)
	quoted2 := pq.QuoteIdentifier(tables.Table2)
	var (
		exists1 bool
		exists2 bool
	)
	sqlStatement1 := fmt.Sprintf(`SELECT 1 FROM %s LIMIT 1`, quoted1)
	sqlStatement2 := fmt.Sprintf(`SELECT 1 FROM %s LIMIT 1`, quoted2)
	rows1, err1 := db.Query(sqlStatement1)
	if err1 != nil {
		log.Println(tables.Table1, " does not exist")
	} else {
		exists1 = true
		log.Println(err1)
		defer rows1.Close();
	}
	rows2, err2 := db.Query(sqlStatement2)
	if err2 != nil {
		log.Println(tables.Table2, " does not exist")
	} else {
		exists2 = true
		log.Println(err2)
		defer rows2.Close();
	}
	log.Println(exists1)
	log.Println(exists2)
	if exists1 {
		return &pb.CheckResponse{Table: tables.Table1}, nil
	} else if exists2 {
		return &pb.CheckResponse{Table: tables.Table2}, nil
	} else {
		quoted := pq.QuoteIdentifier(tables.Table1)
		sqlStatement := fmt.Sprintf(`CREATE TABLE %s (count INT NOT NULL UNIQUE, author TEXT NOT NULL, message TEXT NOT NULL, deleted BOOL)`, quoted)
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Println(err)
		}
		sqlStatement = fmt.Sprintf(`
	INSERT INTO %s (count, author, message, deleted)
	VALUES ($1, $2, $3, $4)`, quoted)
		rows4, err := db.Query(sqlStatement, 0, "admin", "Say hi", false)
		if err != nil {
			log.Println(err)
		}
		defer rows4.Close();
		return &pb.CheckResponse{Table: tables.Table1}, nil
	}
}

func (s *server) Delete(ctx context.Context, message *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	quoted := pq.QuoteIdentifier(message.Table)
	sqlStatement := fmt.Sprintf(`
	UPDATE %s
    SET deleted = $1
    WHERE
    count = $2`, quoted)
	rows, err := db.Query(sqlStatement, true, message.Count)
	if err != nil {
		log.Println(err)
		return &pb.DeleteResponse{Success: false}, nil
	}
	defer rows.Close()
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *server) Restore(ctx context.Context, message *pb.RestoreRequest) (*pb.RestoreResponse, error) {
	quoted := pq.QuoteIdentifier(message.Table)
	sqlStatement := fmt.Sprintf(`
	UPDATE %s
    SET deleted = $1
    WHERE
    count = $2`, quoted)
	rows, err := db.Query(sqlStatement, false, message.Count)
	if err != nil {
		log.Println(err)
		return &pb.RestoreResponse{Success: false}, nil
	}
	defer rows.Close()
	return &pb.RestoreResponse{Success: true}, nil
}