echo Sha zaeboshu
./migrations/migrate -source file://migrations/migration_list -database postgres://superuser:superuser@localhost:5432/mydb?sslmode=disable $*
echo huyak