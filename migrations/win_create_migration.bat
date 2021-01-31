@echo off
echo Sha zaeboshu
migrations\win_migrate -source file://migrations/migration_list -database postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable create -ext sql -dir migrations/migration_list %*
echo huyak