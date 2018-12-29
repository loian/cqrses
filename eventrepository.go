package cqrses

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type EventRepository interface {
	Load(eventType EventType, aggregateId AggregateID) (AggregateRoot, error)
	Store(aggregateRoot AggregateRoot) error
}

type MysqlEventStore struct {
	connstr string
	eventBus EventBus

}

type ErrorConnection struct {
	Message string
}

func (e ErrorConnection) Error() string {
	return e.Message
}

type ErrorConcurrencyViolation struct {
	Message string
}

func (e ErrorConcurrencyViolation) Error() string {
	return e.Message
}



//connect to the database, it expect to fine a table
// CREATE TABLE events (
// 	 event_id CHAR(36) NOT NULL CHARACTER SET ascii,
// 	 aggregate_id CHAR(36) NOT NULL CHARACTER SET ascii,
//	 version BIGINT UNSIGNED NOT NULL,
//   event JSON,
//	 created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
//	 PRIMARY KEY (event_id),
//	 UNIQUE KEY (aggregate_id, version)
// )ENGINE=INNODB;


func (mystore *MysqlEventStore) connect() (*sql.DB, error){
	db, err := sql.Open("mysql", mystore.connstr)

	if err!= nil {
		return nil, ErrorConnection{fmt.Sprintf("unable to connect to mysql instance at %s", mystore.connstr)}
	}

	return db, nil
}


func (mystore *MysqlEventStore) storeChange(eventMessage EventMessage, expectedVersion int64) error {

	db, err := mystore.connect()
	defer db.Close()

	if err != nil {
		return err
	}

	query := `INSERT INTO events (?, ?, ?)
				SELECT aggregateID, version, event
				FROM dual
				WHERE NOT EXISTS (
					SELECT 1 FROM events
					WHERE 
						aggregate_id = ? AND
						version > ?
				)`

	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(string(eventMessage.AggregateID()), eventMessage.Event(), expectedVersion)

	if err != nil {
		return ErrorConnection{
			fmt.Sprintf(
				`version violation while saving change of aggregate "%s" with expected version "%s"`,
				string(eventMessage.AggregateID()),
				string(expectedVersion))}
	}

	return nil
}

func (mystore *MysqlEventStore) Store(root AggregateRoot) error{
 	for _, em := range root.Changes() {
 		err := mystore.storeChange(em, root.InitialVersion() + 1)

 		if err != nil {
 			return err
		}
 		root.IncrementVersion()
	}

 	return nil
}