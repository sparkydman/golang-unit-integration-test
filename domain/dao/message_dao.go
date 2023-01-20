package domain

import (
	"database/sql"
	dto "efficient_api/domain/dto"
	"efficient_api/utils/error_formats"
	"efficient_api/utils/error_utils"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	MessageRepo MessageRepoInterface = &messageRepo{}
)

const (
	queryGetMessage    = "SELECT id, title, body, created_at FROM messages WHERE id=?;"
	queryInsertMessage = "INSERT INTO messages(title, body, created_at) VALUES(?, ?, ?);"
	queryUpdateMessage = "UPDATE messages SET title=?, body=? WHERE id=?;"
	queryDeleteMessage = "DELETE FROM messages WHERE id=?;"
	queryGetAllMessages = "SELECT id, title, body, created_at FROM messages;"
)


type MessageRepoInterface interface {
	Get(int64) (*dto.Message, error_utils.MessageErr)
	Create(*dto.Message) (*dto.Message, error_utils.MessageErr)
	Update(*dto.Message) (*dto.Message, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	GetAll() ([]dto.Message, error_utils.MessageErr)
	Initialize(string, string, string, string, string, string) *sql.DB
}
type messageRepo struct {
	db *sql.DB
}

func (mr *messageRepo) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *sql.DB  {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	mr.db, err = sql.Open(Dbdriver, DBURL)
	if err != nil {
		log.Fatal("This is the error connecting to the database:", err)
	}
	fmt.Printf("We are connected to the %s database\n", Dbdriver)

	return mr.db
}

func NewMessageRepository(db *sql.DB) MessageRepoInterface {
	return &messageRepo{db: db}
}

func (mr *messageRepo) Get(messageId int64) (*dto.Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryGetMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerErrorMessage(fmt.Sprintf("Error when trying to prepare message: %s", err.Error()))
	}
	defer stmt.Close()

	var msg dto.Message
	result := stmt.QueryRow(messageId)
	if getError := result.Scan(&msg.Id, &msg.Title, &msg.Body, &msg.CreatedAt); getError != nil {
		fmt.Println("this is the error man: ", getError)
		return nil,  error_formats.ParseError(getError)
	}
	return &msg, nil
}

func (mr *messageRepo) GetAll() ([]dto.Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryGetAllMessages)
	if err != nil {
		return nil, error_utils.NewInternalServerErrorMessage(fmt.Sprintf("Error when trying to prepare all messages: %s", err.Error()))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil,  error_formats.ParseError(err)
	}
	defer rows.Close()

	results := make([]dto.Message, 0)

	for rows.Next() {
		var msg dto.Message
		if getError := rows.Scan(&msg.Id, &msg.Title, &msg.Body, &msg.CreatedAt); getError != nil {
			return nil, error_utils.NewInternalServerErrorMessage(fmt.Sprintf("Error when trying to get message: %s", getError.Error()))
		}
		results = append(results, msg)
	}
	if len(results) == 0 {
		return nil, error_utils.NewNotFoundErrorMessage("no records found")
	}
	return results, nil
}

func (mr *messageRepo) Create(msg *dto.Message) (*dto.Message, error_utils.MessageErr) {
	fmt.Println("WE REACHED THE DOMAIN")
	stmt, err := mr.db.Prepare(queryInsertMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerErrorMessage(fmt.Sprintf("error when trying to prepare user to save: %s", err.Error()))
	}
	fmt.Println("WE DIDNT REACH HERE")

	defer stmt.Close()

	insertResult, createErr := stmt.Exec(msg.Title, msg.Body, msg.CreatedAt)
	if createErr != nil {
		return nil,  error_formats.ParseError(createErr)
	}
	msgId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, error_utils.NewInternalServerErrorMessage(fmt.Sprintf("error when trying to save message: %s", err.Error()))
	}
	msg.Id = msgId

	return msg, nil
}

func (mr *messageRepo) Update(msg *dto.Message) (*dto.Message, error_utils.MessageErr) {
	stmt, err := mr.db.Prepare(queryUpdateMessage)
	if err != nil {
		return nil, error_utils.NewInternalServerErrorMessage(fmt.Sprintf("error when trying to prepare user to update: %s", err.Error()))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(msg.Title, msg.Body, msg.Id)
	if updateErr != nil {
		return nil,  error_formats.ParseError(updateErr)
	}
	return msg, nil
}

func (mr *messageRepo) Delete(msgId int64) error_utils.MessageErr {
	stmt, err := mr.db.Prepare(queryDeleteMessage)
	if err != nil {
		return error_utils.NewInternalServerErrorMessage(fmt.Sprintf("error when trying to delete message: %s", err.Error()))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(msgId); err != nil {
		return error_utils.NewInternalServerErrorMessage(fmt.Sprintf("error when trying to delete message %s", err.Error()))
	}
	return nil
}
