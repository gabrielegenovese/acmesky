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
    (2,"Aeroporto del Salento", "Brindisi"),
    (3,"Aeroporto di Roma Fiumicino", "Roma"),
    (4,"Aeroporto di Milano Malpensa", "Milano"),
    (5, "Aeroporto di Amsterdam Schiphol", "Amsterdam"),
    (6, "Aeroporto di Los Angeles", "Los Angeles"),
    (7, "Aeroporto di Venezia Marco Polo", "Venezia"),
    (8, "Aeroporto di Napoli Capodichino", "Napoli"),
    (9, "Aeroporto di Bergamo Orio al Serio", "Bergamo"),
    (10, "Aeroporto di Torino Caselle", "Torino"),
    (11, "Aeroporto di Singapore Changi", "Singapore"),
    (12, "Aeroporto di Palermo Falcone e Borsellino", "Palermo"),
    (13, "Aeroporto di Catania Fontanarossa", "Catania"),
    (14, "Aeroporto di Firenze Peretola", "Firenze"),
    (15, "Aeroporto di Parigi Charles de Gaulle", "Parigi"),
    (16, "Aeroporto di Londra Heathrow", "Londra"),
    (17, "Aeroporto di New York John F. Kennedy", "New York"),
    (18, "Aeroporto di Francoforte", "Francoforte"),
    (19, "Aeroporto di Madrid Barajas", "Madrid"),
    (20, "Aeroporto di Dubai", "Dubai"),
    (21, "Aeroporto di Tokyo Narita", "Tokyo"),
    (22, "Aeroporto di Sydney Kingsford Smith", "Sydney"),
    (23, "Aeroporto di San Francisco", "San Francisco"),
    (24, "Aeroporto di Istanbul", "Istanbul"),
    (25, "Aeroporto di Zurigo", "Zurigo"),
    (26, "Aeroporto di Mosca Sheremetyevo", "Mosca");

CREATE TABLE TravelPreferences (
    TravelPreferenceID int NOT NULL AUTO_INCREMENT,
    ProntogramID VARCHAR(256) NOT NULL,
    Budget DECIMAL(8, 2) NOT NULL,
    TravelDateStart DATE NOT NULL,
    TravelDateEnd DATE NOT NULL,
    AirportOriginID int NOT NULL,
    AirportDestinationID int NOT NULL,
    SeatsCount int NOT NULL DEFAULT(1),
    PRIMARY KEY (TravelPreferenceID),
    FOREIGN KEY (AirportOriginID) REFERENCES Airports(AirportID),
    FOREIGN KEY (AirportDestinationID) REFERENCES Airports(AirportID)
);

CREATE TABLE Customers (
    CustomerID int NOT NULL AUTO_INCREMENT,
    Name VARCHAR(64) NOT NULL,
    Surname VARCHAR(64) NOT NULL,
    PRIMARY KEY (CustomerID)
);

CREATE TABLE FlightCompanies (
    CompanyID int NOT NULL AUTO_INCREMENT,
    CompanyName VARCHAR(64) NOT NULL,
    PRIMARY KEY (CompanyID)
);

INSERT INTO FlightCompanies
    (CompanyID, CompanyName)
VALUES
    (1, "FlightCompany");

CREATE TABLE Flights (
    CompanyFlightID int NOT NULL,
    CompanyID int NOT NULL,
    AirportOriginID int NOT NULL,
    AirportDestinationID int NOT NULL,
    DepartDatetime DATETIME NOT NULL,
    ArrivalDatetime DATETIME NOT NULL,
    PassengerFlightPrice DECIMAL(8, 2) NOT NULL,
    AvailableSeats int NOT NULL DEFAULT(0),
    PRIMARY KEY (CompanyFlightID, CompanyID),
    FOREIGN KEY (CompanyID) REFERENCES FlightCompanies(CompanyID)
);

CREATE TABLE ReservedOffers (
    OfferCode int NOT NULL AUTO_INCREMENT,
    TravelPreferenceID int NOT NULL,
    StartReservationDatetime DATETIME DEFAULT NOW(),
    EndReservationDatetime DATETIME AS (DATE_ADD(StartReservationDatetime, INTERVAL 24 HOUR)),
    TotalOfferPrice DECIMAL(8, 2) NOT NULL,
    PRIMARY KEY (OfferCode),
    FOREIGN KEY (TravelPreferenceID) REFERENCES TravelPreferences(TravelPreferenceID)
);

CREATE TABLE OffersBundles(
    OfferCode int NOT NULL,
    CompanyFlightID int NOT NULL,
    CompanyID int NOT NULL,
    PRIMARY KEY (OfferCode, CompanyFlightID, CompanyID),
    FOREIGN KEY (OfferCode) REFERENCES ReservedOffers(OfferCode),
    FOREIGN KEY (CompanyFlightID, CompanyID) REFERENCES Flights(CompanyFlightID, CompanyID)
);

CREATE TABLE SoldOffers (
    OfferCode int NOT NULL,
    TravelPreferenceID int NOT NULL,
    EndReservationDatetime DATETIME DEFAULT NULL,
    CustomerFlightPrice DECIMAL(8, 2) NOT NULL,
    PRIMARY KEY (OfferCode),
    FOREIGN KEY (TravelPreferenceID) REFERENCES TravelPreferences(TravelPreferenceID)
);
