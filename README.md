##install

##create migration
`migrate -source file://./backend/storage/migrations -database postgres://localhost:5432/meteostaion create -dir backend/storage/migrations -ext sql initialize`