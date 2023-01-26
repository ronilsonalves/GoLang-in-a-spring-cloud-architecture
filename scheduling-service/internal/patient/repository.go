package patient

import (
	"errors"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/pkg/store"
	"log"
)

var table = "patients"

type Repository interface {
	GetAll() (interface{}, error)
	GetByID(id int) (interface{}, error)
	Create(p domain.Patient) (interface{}, error)
	Update(id int, p domain.Patient) (interface{}, error)
	Delete(id int) error
}

type repository struct {
	store store.Store
}

func NewRepository(store store.Store) Repository {
	return &repository{store}
}

// GetAll - returns all patients at database
func (r *repository) GetAll() (interface{}, error) {
	return r.store.GetAll(table)
}

func (r *repository) GetByID(id int) (interface{}, error) {
	return r.store.GetByID(id, table)
}

func (r *repository) Create(p domain.Patient) (interface{}, error) {
	if !r.validateIdentificationNumber(p.RG) {
		return nil, errors.New("license number already exists at database")
	}
	return r.store.Save(p, table)
}

func (r *repository) Update(id int, p domain.Patient) (interface{}, error) {

	pInterface, err := r.GetAll()
	if err != nil {
		log.Fatalln("error while trying to fetch data from db while update a patient")
		return nil, err
	}
	patients, ok := pInterface.([]domain.Patient)
	if !ok {
		log.Fatalln("error while trying to fetch data from db while updating a patient")
		return nil, err
	}

	for _, patient := range patients {
		if patient.Id == id {
			if !r.validateIdentificationNumber(p.RG) && p.RG != p.RG {
				return nil, errors.New("there's a patient with same identity number")
			}
			return r.store.Update(id, p, table)
		}
	}
	return nil, errors.New("patient not found")
}

func (r *repository) Delete(id int) error {
	return r.store.Delete(id, table)
}

func (r *repository) validateIdentificationNumber(identityNumber string) bool {
	var patients []domain.Patient
	patientsInterface, err := r.GetAll()
	if err != nil {
		log.Fatalln("erro while trying to fetch data from db")
		return false
	}
	patients, ok := patientsInterface.([]domain.Patient)
	if !ok {
		log.Fatalln("error while trying to fetch data from db")
		return false
	}

	for _, patient := range patients {
		if patient.RG == identityNumber {
			return false
		}
	}

	return true
}
