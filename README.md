1. export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/slackchat?sslmode=disable'
2. migrate -database ${POSTGRESQL_URL} -path db/migrations up