CREATE DATABASE IF NOT EXISTS flightcompany_db;

USE flightcompany_db;

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

CREATE TABLE Flights (
    FlightID int NOT NULL AUTO_INCREMENT,
    AirportOriginID int NOT NULL,
    AirportDestinationID int NOT NULL,
    DepartDatetime DATETIME NOT NULL,
    ArrivalDatetime DATETIME NOT NULL,
    FlightPrice DECIMAL(8, 2) NOT NULL,
    AvailableSeats int NOT NULL DEFAULT(0),
    PRIMARY KEY (FlightID),
    FOREIGN KEY (AirportOriginID) REFERENCES Airports(AirportID),
    FOREIGN KEY (AirportDestinationID) REFERENCES Airports(AirportID)
);

INSERT INTO Flights 
    (FlightID, AirportOriginID, AirportDestinationID, DepartDatetime, ArrivalDatetime, AvailableSeats, FlightPrice )
VALUES
    (1, 1, 2, "2024-01-01 00:00:00", "2024-01-01 01:30:00", 100, 60.50),
    (2, 2, 1, "2024-01-01 02:00:00", "2024-01-01 03:30:00", 100, 70.50),
    (3, 1, 2, "2024-01-01 04:00:00", "2024-01-01 05:30:00", 100, 80.50),
    (4, 2, 1, "2024-01-01 06:00:00", "2024-01-01 07:30:00", 100, 90.50),
    (5, 1, 2, "2024-01-01 08:00:00", "2024-01-01 09:30:00", 100, 100.50),
    (6, 2, 1, "2024-01-01 10:00:00", "2024-01-01 11:30:00", 100, 110.50),
    (7, 1, 2, "2024-01-01 12:00:00", "2024-01-01 13:30:00", 100, 120.50),
    (8, 2, 1, "2024-01-01 15:00:00", "2024-01-01 16:30:00", 100, 250.50),
    (9, 1, 2, "2024-01-01 17:00:00", "2024-01-01 18:30:00", 100, 350.50),
    (10, 2, 1, "2024-01-01 23:00:00", "2024-01-02 00:30:00", 100, 150.50),
    (11, 1, 2, "2024-01-02 01:00:00", "2024-01-02 02:30:00", 100, 230.50);

CREATE TABLE FlightBookings (
    BookingID int NOT NULL AUTO_INCREMENT,
    CustomerName VARCHAR(64) NOT NULL,
    CustomerSurname VARCHAR(64) NOT NULL,
    FlightID int NOT NULL,
    SeatsCount int NOT NULL DEFAULT(1),
    ReservationFlightPrice DECIMAL(8, 2) NOT NULL,
    ReservationDatetime DATETIME NOT NULL DEFAULT(CURRENT_TIMESTAMP),
    BoughtDatetime DATETIME NULL DEFAULT(NULL),
    PRIMARY KEY (BookingID),
    FOREIGN KEY (FlightID) REFERENCES Flights(FlightID)
);

/*
CREATE TABLE FlightsForReservedOffers (
    ReservedOfferCode int NOT NULL AUTO_INCREMENT,
    CompanyFlightID int NOT NULL,
    CompanyID int NOT NULL,
    PRIMARY KEY (ReservedOfferCode, CompanyFlightID, CompanyID),
    FOREIGN KEY (ReservedOfferCode) REFERENCES ReservedOffers(ReservedOfferCode),
    FOREIGN KEY (CompanyFlightID, CompanyID) REFERENCES Flights(CompanyFlightID, CompanyID)
);
*/