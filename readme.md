# REST service
Golang REST api service with PostgreSQL, Redis, gin-gonic and gorm.

There are features of this app:
1. JWT Authorization
2. Docker / Docker compose
3. Database Gorm + PostgreSQL
4. Cached db by Redis
5. [Postman documentation](<https://documenter.getpostman.com/view/11918079/TVmHFg9u> "Postman documentation")

## Installation
First of all download project

Run docker-compose with terminal
>docker-compose up --build

Open terminal of backend container for initial setup
> docker exec -it agyn-backend /bin/bash

Create standard roles
>go run cli/create-roles/main.go

Create admin user
>go run cli/create-admin/main.go

Create 10 test tasks
>go run cli/create-test-tasks/main.go

Open webpage
>localhost:8081

Done!

### Connection to database:

Open terminal of postgres container
>docker exec -it go_ismael_postgres_1 /bin/bash

Open postgreSQL terminal
>psql "dbname=agyn_test_rest host=localhost user=agyn password=agyn port=5432"
