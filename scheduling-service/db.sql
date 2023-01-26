CREATE SCHEMA IF NOT EXISTS dental_clinic;
USE dental_clinic;
CREATE TABLE dentists (
    id INT NOT NULL AUTO_INCREMENT,
    last_name VARCHAR(50) NOT NULL,
    name VARCHAR(25) NOT NULL,
    cro VARCHAR(10) NOT NULL UNIQUE,

    PRIMARY KEY (id)
)ENGINE = INNODB;

CREATE TABLE patients (
    id INT NOT NULL AUTO_INCREMENT,
    last_name VARCHAR(50) NOT NULL,
    name VARCHAR(25) NOT NULL,
    rg VARCHAR(10) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,

    PRIMARY KEY (id)
)ENGINE = INNODB;

CREATE TABLE appointments (
    id INT NOT NULL AUTO_INCREMENT,
    description VARCHAR(250) NOT NULL,
    date_and_time DATETIME NOT NULL,
    dentist_cro VARCHAR(10) NOT NULL,
    patient_rg VARCHAR(10) NOT NULL,

    PRIMARY KEY (id),

    CONSTRAINT fk_dentist
                          FOREIGN KEY (dentist_cro)
                          REFERENCES dentists(cro),
    CONSTRAINT fk_patient
                          FOREIGN KEY (patient_rg)
                          REFERENCES patients(rg)
)ENGINE = INNODB;