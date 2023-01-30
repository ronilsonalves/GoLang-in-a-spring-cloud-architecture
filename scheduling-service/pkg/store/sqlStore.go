package store

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/config"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	"log"
	"time"
)

var (
	AP = "appointments"
	DE = "dentists"
	PE = "patients"
)

func init() {
	config.LoadConfig()
}

// NewSQLStore - Initialize Store interface
func NewSQLStore() Store {
	database, err := config.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	return &sqlStore{
		db: database,
	}
}

type sqlStore struct {
	db *sql.DB
}

// GetAll - Return all rows from selected table.
func (s *sqlStore) GetAll(tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxGetAllByTable(tableName, s)
	case DE:
		return auxGetAllByTable(tableName, s)
	case PE:
		return auxGetAllByTable(tableName, s)
	default:
		return nil, errors.New("an error occurred while trying to get data from db")
	}
}

// GetByID - Return a row from selected table by ID
func (s *sqlStore) GetByID(id int, tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxGetByIDByTable(tableName, id, s)
	case DE:
		return auxGetByIDByTable(tableName, id, s)
	case PE:
		return auxGetByIDByTable(tableName, id, s)
	default:
		return nil, errors.New("an error occurred while trying to get data from db")
	}
}

// Save - Insert into selected table a new row
func (s *sqlStore) Save(entity interface{}, tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxSave(tableName, s, entity)
	case DE:
		return auxSave(tableName, s, entity)
	case PE:
		return auxSave(tableName, s, entity)
	default:
		return nil, errors.New("failed to start inserting your data")
	}
}

// Update - update a row from selected table by ID
func (s *sqlStore) Update(entityID int, entity interface{}, tableName string) (interface{}, error) {
	switch tableName {
	case AP:
		return auxUpdate(tableName, s, entity, entityID)
	case DE:
		return auxUpdate(tableName, s, entity, entityID)
	case PE:
		return auxUpdate(tableName, s, entity, entityID)
	default:
		return nil, errors.New("failed to update entity data")
	}
}

// Delete - exclude a row from selected table by ID
func (s *sqlStore) Delete(entityID int, tableName string) error {
	switch tableName {
	case AP:
		return auxDelete(tableName, s, entityID)
	case DE:
		return auxDelete(tableName, s, entityID)
	case PE:
		return auxDelete(tableName, s, entityID)
	default:
		return errors.New("failed to delete")
	}
}

// auxGetAllByTable - Called function by GetAll, here the selected table is validated and all select queries are made.
func auxGetAllByTable(tableName string, s *sqlStore) (interface{}, error) {
	var entities []struct{}

	switch tableName {
	case AP:
		Query := "SELECT a.id, a.description, DATE_FORMAT(a.date_and_time,'%d/%m/%Y %H:%i') date_and_time,a.dentist_cro,a.patient_rg,d.id,d.last_name,d.name,d.cro,p.id,p.last_name,p.name,p.rg,DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') created_at FROM appointments a INNER JOIN dentists d on a.dentist_cro = d.cro INNER JOIN patients p on a.patient_rg = p.rg ORDER BY a.date_and_time"
		rows, err := s.db.Query(Query)
		if err != nil {
			return entities, err
		}

		//var appointment domain.Appointment
		//var appointments []domain.Appointment
		var appointment domain.AppointmentDTO
		var appointments []domain.AppointmentDTO

		for rows.Next() {
			if err := rows.Scan(
				&appointment.Id,
				&appointment.Description,
				&appointment.DateAndTime,
				&appointment.DentistCRO,
				&appointment.PatientRG, //; err != nil {
				//&appointment.Dentist.CRO,
				//&appointment.Patient.RG,
				&appointment.Dentist.Id,
				&appointment.Dentist.LastName,
				&appointment.Dentist.Name,
				&appointment.Dentist.CRO,
				&appointment.Patient.Id,
				&appointment.Patient.LastName,
				&appointment.Patient.Name,
				&appointment.Patient.RG,
				&appointment.Patient.CreatedAt); err != nil {
				return appointments, err
			}
			appointments = append(appointments, appointment)
		}
		return appointments, nil
	case DE:
		rows, err := s.db.Query("SELECT * FROM dentists")
		if err != nil {
			return entities, err
		}

		var dentist domain.Dentist
		var dentists []domain.Dentist

		for rows.Next() {
			if err := rows.Scan(
				&dentist.Id,
				&dentist.LastName,
				&dentist.Name,
				&dentist.CRO); err != nil {
				return dentists, err
			}
			dentists = append(dentists, dentist)
		}
		return dentists, nil
	case PE:
		rows, err := s.db.Query("SELECT p.id, p.last_name,p.name,p.rg, DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') FROM patients p")
		if err != nil {
			return entities, err
		}

		var patient domain.Patient
		var patients []domain.Patient
		for rows.Next() {
			if err := rows.Scan(
				&patient.Id,
				&patient.LastName,
				&patient.Name,
				&patient.RG,
				&patient.CreatedAt); err != nil {
				return patients, err
			}
			patients = append(patients, patient)
		}
		return patients, nil
	default:
		return nil, errors.New("failed to load data from db at sqlStore.go file")
	}
}

// auxGetByIDByTable - Called function by GetByID, here the selected table is validated and all select * from *table_name* where id = *entity_id* are made.
func auxGetByIDByTable(tableName string, entityID int, s *sqlStore) (interface{}, error) {
	var entity struct{}

	switch tableName {
	case AP:
		query := "SELECT a.id, a.description, DATE_FORMAT(a.date_and_time,'%d/%m/%Y %H:%i') date_and_time,a.dentist_cro,a.patient_rg,d.id,d.last_name,d.name,d.cro,p.id,p.last_name,p.name,p.rg,DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') created_at FROM appointments a INNER JOIN dentists d on a.dentist_cro = d.cro INNER JOIN patients p on a.patient_rg = p.rg WHERE a.id = ? ORDER BY a.date_and_time"
		rows, err := s.db.Query(query, entityID)
		if err != nil {
			return entity, err
		}
		defer rows.Close()

		var appointment domain.AppointmentDTO
		for rows.Next() {
			if err := rows.Scan(
				&appointment.Id,
				&appointment.Description,
				&appointment.DateAndTime,
				&appointment.DentistCRO,
				&appointment.PatientRG,
				&appointment.Dentist.Id,
				&appointment.Dentist.LastName,
				&appointment.Dentist.Name,
				&appointment.Dentist.CRO,
				&appointment.Patient.Id,
				&appointment.Patient.LastName,
				&appointment.Patient.Name,
				&appointment.Patient.RG,
				&appointment.Patient.CreatedAt); err != nil {
				return appointment, err
			}
			return appointment, nil
		}
		if rows.Next() {
			return appointment, nil
		}
		return nil, err
	case DE:
		rows, err := s.db.Query("SELECT * FROM dentists WHERE id = ?", entityID)
		if err != nil {
			return entity, err
		}
		defer rows.Close()

		var dentist domain.Dentist
		for rows.Next() {
			if err = rows.Scan(
				&dentist.Id,
				&dentist.LastName,
				&dentist.Name,
				&dentist.CRO); err != nil {
				return dentist, err
			}
			return dentist, nil
		}
		if rows.Next() {
			return dentist, nil
		}
		return nil, err
	case PE:
		rows, err := s.db.Query("SELECT p.id, p.last_name,p.name,p.rg, DATE_FORMAT(p.created_at,'%d/%m/%Y %H:%i') FROM patients p WHERE id = ?", entityID)
		if err != nil {
			return entity, err
		}
		defer rows.Close()

		var patient domain.Patient
		for rows.Next() {
			if err = rows.Scan(
				&patient.Id,
				&patient.LastName,
				&patient.Name,
				&patient.RG,
				&patient.CreatedAt); err != nil {
				return nil, err
			}
			return patient, nil
		}
		if rows.Next() {
			return patient, nil
		}
		return nil, err
	default:
		return nil, errors.New("failed to get by id from db")
	}
}

// auxSave - Called function by Save, here the inserts are made into selected table.
func auxSave(tableName string, s *sqlStore, entity interface{}) (interface{}, error) {
	switch tableName {
	case AP:
		log.Println("... inserting data into appointments table.")
		var appointment domain.Appointment
		appointment, ok := entity.(domain.Appointment)
		if ok {
			apDateAndTimeParsed, err := time.Parse("02/01/2006 15:04", appointment.DateAndTime)
			if err != nil {
				log.Println("failed to convert datetimee")
				return nil, errors.New("failed to convert datetimee")
			}
			//
			log.Println(apDateAndTimeParsed.String())
			result, err := s.db.Exec("INSERT INTO appointments(DESCRIPTION, DATE_AND_TIME, dentist_license, patient_identity) VALUES(?,?,?,?)",
				appointment.Description,
				apDateAndTimeParsed,
				appointment.DentistCRO,
				appointment.PatientRG) //err != nil {
			//&appointment.Dentist.CRO,
			//&appointment.Patient.RG,
			//&appointment.Dentist.Id,
			//&appointment.Dentist.LastName,
			//&appointment.Dentist.Name,
			//&appointment.Dentist.CRO,
			//&appointment.Patient.Id,
			//&appointment.Patient.LastName,
			//&appointment.Patient.Name,
			//&appointment.Patient.RG,
			//&appointment.Patient.CreatedAt)
			if err != nil {
				fmt.Println("inserting data failed :", err.Error())
				return nil, err
			}
			lastInsertedID, err := result.LastInsertId()
			if err != nil {
				fmt.Println("error trying to get id inserted:", err.Error())
				return nil, err
			}
			appointment.Id = int(lastInsertedID)
			log.Println("... INSERT operation was successfully")
			return s.GetByID(appointment.Id, AP)
			//return response, nil
		}
	case DE:
		var dentist domain.Dentist
		dentist, ok := entity.(domain.Dentist)
		if ok {
			result, err := s.db.Exec("INSERT INTO dentists(surname, name, license_number) VALUES (?,?,?)",
				dentist.LastName,
				dentist.Name,
				dentist.CRO)
			if err != nil {
				fmt.Println("inserting data failed :", err.Error())
				return nil, err
			}
			lastInsertedID, err := result.LastInsertId()
			if err != nil {
				fmt.Println("error trying to get id inserted:", err.Error())
				return nil, err
			}
			dentist.Id = int(lastInsertedID)
			fmt.Println("dentist inserted at db:", dentist)
			return dentist, nil
		}
	case PE:
		var patient domain.Patient
		patient, ok := entity.(domain.Patient)
		if ok {
			patCreatedAtParsed, err := time.Parse("02/01/2006 15:04:05", patient.CreatedAt)
			if err != nil {
				return nil, errors.New("failed to convert patient created_at field")
			}
			result, err := s.db.Exec("INSERT INTO patients(surname, name, identity_number, created_at) VALUES (?,?,?,?)",
				patient.LastName,
				patient.Name,
				patient.RG,
				patCreatedAtParsed)
			if err != nil {
				fmt.Println("inserting data failed :", err.Error())
				return nil, err
			}
			lastInsertedID, err := result.LastInsertId()
			if err != nil {
				fmt.Println("error trying to get id inserted:", err.Error())
				return nil, err
			}
			patient.Id = int(lastInsertedID)
			return patient, nil
		}
	default:
		return nil, errors.New("failed to insert data at database")
	}
	return nil, errors.New("failed to insert data at database")
}

// auxUpdate - Called function by Update, here the updates are made into selected table.
func auxUpdate(tableName string, s *sqlStore, entity interface{}, entityId int) (interface{}, error) {
	switch tableName {
	case AP:
		var appointment domain.Appointment
		appointment, ok := entity.(domain.Appointment)
		if ok {
			apDateAndTimeParsed, err := time.Parse("02/01/2006 15:04", appointment.DateAndTime)
			if err != nil {
				log.Println(err.Error(), "\nDate parsed: ", apDateAndTimeParsed)
				return nil, errors.New("failed to convert datetime")
			}
			_, err = s.db.Exec("UPDATE appointments SET description = ?, date_and_time = ?, dentist_license = ?, patient_identity = ? WHERE id = ?",
				appointment.Description,
				apDateAndTimeParsed,
				appointment.DentistCRO,
				appointment.PatientRG,
				entityId)
			if err != nil {
				return nil, err
			}
			return s.GetByID(entityId, AP)
		}
	case DE:
		var dentist domain.Dentist
		dentist, ok := entity.(domain.Dentist)
		if ok {
			_, err := s.db.Exec("UPDATE dentists SET surname = ?, name = ?, license_number = ? WHERE id = ?",
				dentist.LastName,
				dentist.Name,
				dentist.CRO,
				entityId)
			if err != nil {
				return nil, err
			}
			return dentist, nil
		}
	case PE:
		var patient domain.Patient
		patient, ok := entity.(domain.Patient)
		if ok {
			paCreatedAtParsed, err := time.Parse("02/01/2006 15:04", patient.CreatedAt)
			if err != nil {
				return nil, errors.New("failed to convert patient created_at field: " + patient.CreatedAt)
			}

			_, err = s.db.Exec("UPDATE patients SET surname = ?, name = ?, identity_number = ?, created_at = ? WHERE id = ?",
				patient.LastName,
				patient.Name,
				patient.RG,
				paCreatedAtParsed,
				entityId)
			if err != nil {
				return nil, err
			}
			return patient, nil
		}
	default:
		return nil, errors.New("failed to update data into database")
	}
	return nil, errors.New("failed to update data into database")
}

// auxDelete - Called function by Delete, here the deletes are made at selected table.
func auxDelete(tableName string, s *sqlStore, entityID int) error {
	switch tableName {
	case AP:
		result, err := s.db.Exec("DELETE FROM appointments WHERE id =?", entityID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if count == 0 {
			return errors.New("entity not found at database")
		}
		return nil
	case DE:
		result, err := s.db.Exec("DELETE FROM dentists WHERE id =?", entityID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if count == 0 {
			return errors.New("entity not found at database")
		}
		return nil
	case PE:
		result, err := s.db.Exec("DELETE FROM patients WHERE id =?", entityID)
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if count == 0 {
			return errors.New("entity not found at database")
		}
		return nil
	default:
		return errors.New("failed to delete row")
	}
}
