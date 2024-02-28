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
    FOREIGN KEY (AirportOriginID, AirportDestinationID) REFERENCES Airports(AirportID)
);

CREATE TABLE FlightsForReservedOffers{
    ReservedOfferCode int NOT NULL,
    CompanyFlightID int NOT NULL,
    CompanyID int NOT NULL,
    PRIMARY KEY (ReservedOfferCode, CompanyFlightID, CompanyID),
    FOREIGN KEY (ReservedOfferCode) REFERENCES ReservedOffers(ReservedOfferCode),
    FOREIGN KEY (CompanyFlightID, CompanyID) REFERENCES Flights(CompanyFlightID, CompanyID)
};

CREATE TABLE FlightBookings {
    BookingID int NOT NULL AUTO_INCREMENT,
    CustomerName VARCHAR(64) NOT NULL,
    CustomerSurname VARCHAR(64) NOT NULL,
    FlightID int NOT NULL,
    SeatsCount int NOT NULL DEFAULT(1),
    ReservationFlightPrice DECIMAL(8, 2) NOT NULL,
    ReservationDatetime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    BoughtDatetime DATETIME NULL,
    PRIMARY KEY (BookingID),
    FOREIGN KEY (FlightID) REFERENCES Flights(FlightID)
}
