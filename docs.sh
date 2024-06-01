cd src/ACMESkyNCC/middleware
go mod tidy
swag init --parseDependency --instanceName ncc
swag fmt
cd ../../ACMESkyService
go mod tidy
swag init --parseDependency --instanceName acme
swag fmt
cd ../BankService/server
go mod tidy
swag init --parseDependency --instanceName bank
swag fmt
cd ../../FlightCompanyService
go mod tidy
swag init --parseDependency --instanceName fc
swag fmt
cd ../GeoDistance
go mod tidy
swag init --parseDependency --instanceName geo
swag fmt