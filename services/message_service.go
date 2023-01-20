package services

import (
	dao "efficient_api/domain/dao"
	dto "efficient_api/domain/dto"
	"efficient_api/utils/error_utils"
	"time"
)

var MessagesService messageServiceInterface = &messageService{}

type messageService struct{}

type messageServiceInterface interface {
	GetMessage(int64) (*dto.Message, error_utils.MessageErr)
	CreateMessage(*dto.Message)(*dto.Message, error_utils.MessageErr)
	UpdateMessage(*dto.Message)(*dto.Message, error_utils.MessageErr)
	DeleteMessage(int64) error_utils.MessageErr
	GetAllMessages() ([]dto.Message, error_utils.MessageErr)
}

func (s *messageService) GetMessage(msgId int64) (*dto.Message, error_utils.MessageErr){
	res, err := dao.MessageRepo.Get(msgId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *messageService) CreateMessage(message *dto.Message)(*dto.Message, error_utils.MessageErr){
	err := message.Validate()
	if err != nil {
		return nil, err
	}
	message.CreatedAt = time.Now()
	res, err := dao.MessageRepo.Create(message)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *messageService) UpdateMessage(message *dto.Message)(*dto.Message, error_utils.MessageErr){
	err := message.Validate()
	if err != nil {
		return nil, err
	}
	currentMessage, err := dao.MessageRepo.Get(message.Id)
	if err != nil {
		return nil, err
	}
	currentMessage.Title = message.Title
	currentMessage.Body = message.Body
	res, err := dao.MessageRepo.Update(currentMessage)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *messageService) DeleteMessage(msgId int64) error_utils.MessageErr{
	msg, err := dao.MessageRepo.Get(msgId)
	if err != nil {
		return err
	}
	if err = dao.MessageRepo.Delete(msg.Id); err != nil{
		return err
	}
	
	return nil
}

func (s *messageService) GetAllMessages() ([]dto.Message,error_utils.MessageErr){
	msgs, err := dao.MessageRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
