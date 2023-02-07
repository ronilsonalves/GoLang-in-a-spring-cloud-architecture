package appointment

import (
	"errors"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/amqp"
)

type Service interface {
	GetAll() ([]domain.AppointmentDTO, error)
	GetByID(id int) (domain.AppointmentDTO, error)
	GetAllByIdentityNumber(identityNumber string) ([]domain.AppointmentDTO, error)
	GetAllByLicenseNumber(licenseNumber string) ([]domain.AppointmentDTO, error)
	Create(a domain.Appointment) (domain.AppointmentDTO, error)
	Update(id int, a domain.Appointment) (domain.AppointmentDTO, error)
	Delete(id int) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetAll() ([]domain.AppointmentDTO, error) {
	list, err := s.r.GetAll()
	if err != nil {
		return nil, err
	}
	appointments, ok := list.([]domain.AppointmentDTO)
	if !ok {
		return nil, errors.New("an error occurred while trying to fetch data from database")
	}
	return appointments, nil
}

func (s *service) GetByID(id int) (domain.AppointmentDTO, error) {
	aInterface, err := s.r.GetByID(id)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}
	appointment, ok := aInterface.(domain.AppointmentDTO)
	if !ok {
		return appointment, errors.New("not found an appointment with id provided")
	}
	return appointment, nil
}

func (s *service) GetAllByIdentityNumber(identityNumber string) ([]domain.AppointmentDTO, error) {
	list, err := s.r.GetAllByIdentityNumber(identityNumber)
	if err != nil {
		return nil, err
	}
	appointments, ok := list.([]domain.AppointmentDTO)
	if !ok {
		return nil, errors.New("an error occurred while trying to fetch data from database")
	}
	return appointments, nil
}

func (s *service) GetAllByLicenseNumber(licenseNumber string) ([]domain.AppointmentDTO, error) {
	list, err := s.r.GetAllByLicenseNumber(licenseNumber)
	if err != nil {
		return nil, err
	}
	appointments, ok := list.([]domain.AppointmentDTO)
	if !ok {
		return nil, errors.New("an error occurred while trying to fetch data from database")
	}
	return appointments, nil
}

func (s *service) Create(a domain.Appointment) (domain.AppointmentDTO, error) {
	aSavedInterface, err := s.r.Create(a)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}
	apSaved, ok := aSavedInterface.(domain.AppointmentDTO)
	if ok {
		amqp.PublishMessage(apSaved)
		return apSaved, nil
	}

	return domain.AppointmentDTO{}, errors.New("failed to save a new appointment")
}

func (s *service) Update(id int, a domain.Appointment) (domain.AppointmentDTO, error) {
	aUpdate, err := s.GetByID(id)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}

	if a.Description == "" {
		a.Description = aUpdate.Description
	}
	if a.DateAndTime == "" {
		a.DateAndTime = aUpdate.DateAndTime
	}
	if a.DentistCRO == "" {
		a.DentistCRO = aUpdate.DentistCRO
	}
	if a.PatientRG == "" {
		a.PatientRG = aUpdate.PatientRG
	}
	a.Id = aUpdate.Id

	updated, err := s.r.Update(id, a)
	if err != nil {
		return domain.AppointmentDTO{}, err
	}
	response, ok := updated.(domain.AppointmentDTO)
	if !ok {
		return domain.AppointmentDTO{}, errors.New("failed to update appointment")
	}

	amqp.PublishMessage(response)
	return response, nil
}

func (s *service) Delete(id int) error {
	return s.r.Delete(id)
}

// sendMsg - send a msg to RabbitMQ queue when an appointment is made or updated
//func sendMsg(a domain.AppointmentDTO) {
//	mq, err := amqp.ConnectRabbitMQ(os.Getenv("RABBIT_MQ_URL_CONN"), "appointment-service")
//	log.Println(err)
//	body, _ := json.Marshal(a)
//	err = mq.Publish(&rabbitmq.MQConfigPublish{
//		RoutingKey: mq.Queue().Name,
//		Message: amqpi.Publishing{
//			ContentType: "application/json",
//			Body:        body,
//		},
//	})
//	defer mq.Close()
//}
