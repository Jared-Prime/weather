include .env

connect:
	ssh -i $(PEMFILE) ubuntu@$(PGHOST)

dbconnect:
	PGPASSWORD=$(PGPASSWORD) psql -U postgres -h $(PGHOST) -d weather

generate:
	protoc -I . conditions.proto --go_out=plugins=grpc:conditions
