# Movie list management app
  ### Tools:
  - go 1.19
  - postgres

 ### How to use this on Windows
 Run containers:

```cmd
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
docker-compose up --build crud_movie_manager
```

#### Json example for POST query:
```json
{
    "title": "Witcher",
    "release": "2020",
    "streamingService": "Netflix"
}
```

### API example
![image](https://github.com/BalamutDiana/crud_movie_manager/blob/main/example.gif)
