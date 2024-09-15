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
    (1, 1, 2, "2024-12-01 00:00:00", "2024-12-01 01:30:00", 100, 50.50),
    (2, 2, 1, "2024-12-01 02:00:00", "2024-12-01 03:30:00", 100, 60.50),
    (3, 1, 2, "2024-12-01 04:00:00", "2024-12-01 05:30:00", 100, 70.50),
    (4, 2, 1, "2024-12-01 06:00:00", "2024-12-01 07:30:00", 100, 80.50),
    (5, 1, 2, "2024-12-01 08:00:00", "2024-12-01 09:30:00", 100, 90.50),
    (6, 2, 1, "2024-12-01 10:00:00", "2024-12-01 11:30:00", 100, 100.50),
    (7, 1, 2, "2024-12-01 12:00:00", "2024-12-01 13:30:00", 100, 110.50),
    (8, 2, 1, "2024-12-01 15:00:00", "2024-12-01 16:30:00", 100, 120.50),
    (9, 1, 2, "2024-12-01 17:00:00", "2024-12-01 18:30:00", 100, 130.50),
    (10, 2, 1, "2024-12-01 23:00:00", "2024-12-02 00:30:00", 100, 140.50),
    (11, 1, 2, "2024-12-02 01:00:00", "2024-12-02 02:30:00", 100, 150.50),
    (12, 2, 1, "2024-12-02 06:00:00", "2024-12-02 07:30:00", 100, 160.50),
    (13, 1, 2, "2024-12-02 08:00:00", "2024-12-02 09:30:00", 100, 170.50),
    (14, 2, 1, "2024-12-02 10:00:00", "2024-12-02 11:30:00", 100, 180.50),
    (15, 1, 2, "2024-12-02 12:00:00", "2024-12-02 13:30:00", 100, 190.50),
    (16, 2, 1, "2024-12-02 15:00:00", "2024-12-02 16:30:00", 100, 200.50),
    (17, 1, 2, "2024-12-02 17:00:00", "2024-12-02 18:30:00", 100, 210.50),

    (18, 1, 2, "2024-12-03 00:00:00", "2024-12-03 01:30:00", 100, 550.50),
    (19, 2, 1, "2024-12-03 02:00:00", "2024-12-03 03:30:00", 100, 560.50),
    (20, 1, 2, "2024-12-03 04:00:00", "2024-12-03 05:30:00", 100, 570.50),
    (21, 2, 1, "2024-12-03 06:00:00", "2024-12-03 07:30:00", 100, 580.50),
    (22, 1, 2, "2024-12-03 08:00:00", "2024-12-03 09:30:00", 100, 590.50),
    (23, 2, 1, "2024-12-03 10:00:00", "2024-12-03 11:30:00", 100, 500.50),
    (24, 1, 2, "2024-12-03 12:00:00", "2024-12-03 13:30:00", 100, 510.50),
    (25, 2, 1, "2024-12-03 15:00:00", "2024-12-03 16:30:00", 100, 520.50),
    (26, 1, 2, "2024-12-03 17:00:00", "2024-12-03 18:30:00", 100, 530.50),
    (27, 2, 1, "2024-12-03 23:00:00", "2024-12-04 00:30:00", 100, 540.50),
    (28, 1, 2, "2024-12-04 01:00:00", "2024-12-04 02:30:00", 100, 550.50),
    (29, 2, 1, "2024-12-04 06:00:00", "2024-12-04 07:30:00", 100, 560.50),
    (30, 1, 2, "2024-12-04 08:00:00", "2024-12-04 09:30:00", 100, 570.50),
    (31, 2, 1, "2024-12-04 10:00:00", "2024-12-04 11:30:00", 100, 580.50),
    (32, 1, 2, "2024-12-04 12:00:00", "2024-12-04 13:30:00", 100, 590.50),
    (33, 2, 1, "2024-12-04 15:00:00", "2024-12-04 16:30:00", 100, 500.50),
    (34, 1, 2, "2024-12-04 17:00:00", "2024-12-04 18:30:00", 100, 510.50);
    

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


CREATE VIEW FlightCurrentSeats AS
SELECT F.FlightID, COALESCE(F.AvailableSeats - B.CurrentlyReservedSeats, F.AvailableSeats) AS LeftSeats
FROM Flights F
    LEFT JOIN 
    (
        SELECT FlightID, SUM(SeatsCount) AS CurrentlyReservedSeats
        FROM FlightBookings
        GROUP BY FlightID
    ) B
    ON F.FlightID = B.FlightID
;

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
