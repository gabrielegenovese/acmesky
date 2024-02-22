CREATE DATABASE IF NOT EXISTS acmesky_db;

USE acmesky_db;

CREATE TABLE Airports (
    AirportID int NOT NULL AUTO_INCREMENT,
    Name VARCHAR(64) NOT NULL,
    City VARCHAR(32) NOT NULL,
    PRIMARY KEY (AirportID)
);

INSERT INTO Airports
    (AirportID,Name,City)
VALUES
    (1,"Aeroporto di Bologna-Guglielmo Marconi","Bologna"),
    (2,"Aeroporto del Salento", "Brindisi");


CREATE TABLE TravelPreferences (
    TravelPreferenceID int NOT NULL AUTO_INCREMENT,
    ProntogramID VARCHAR(256) NOT NULL,
    Budget DECIMAL(8, 2) NOT NULL,
    TravelDateStart DATE NOT NULL,
    TravelDateEnd DATE NOT NULL.
    AirportOriginID int NOT NULL,
    AirportDestinationID int NOT NULL,
    PRIMARY KEY (TravelPreferenceID),
    FOREIGN KEY (AirportOriginID) REFERENCES Airports(AirportID),
    FOREIGN KEY (AirportDestinationID) REFERENCES Airports(AirportID)
);

