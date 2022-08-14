# Movie list management app
  ### Tools:
  - go 1.19
  - postgres

 ### How to use this on Windows
 Run container with postgres:
```cmd
docker run -d --name crud-app-v2 -e POSTGRES_PASSWORD=password -v %CD%/pgdata/:/var/lib/postgresql/data -p 5432:5432 postgres
```
Launch container terminal and run psql:
```cmd
docker exec -it crud-movie-app bash
psql -U postgres
```
Сreate a table to store movies:
```sql
create table movies (
id serial not null unique,
title varchar(255) not null,
release varchar(255) not null,
streaming_service varchar(255) not null,
saved_at timestamp not null default now()
);
```

Building and running the application:
```cmd
go build -o app cmd/main.go
./app
```
### API example
![image](https://github.com/BalamutDiana/crud_movie_manager/blob/main/example.gif)
