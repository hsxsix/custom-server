### quick start test postgres server

docker run
```
docker run -d --name postgres-server -p5432:5432 -e POSTGRES_PASSWORD=postgres postgres:13.4
```

connect
```
docker run -it --rm --network host postgres:13.4 psql -h 127.0.0.1 -U postgres
```

new user and database
```
CREATE USER xxxuser WITH PASSWORD '';
CREATE DATABASE xxxdatabase OWNER xxx;
GRANT ALL PRIVILEGES ON DATABASE xxxdatabase TO xxxuser;
```

### quick start test mysql(mariadb) server

docker run
```
docker run -d --name mariadb-server -e MYSQL_ROOT_PASSWORD=mariadb -p 3306:3306 mariadb:10.5 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

connect
```
docker run -it --rm --network host mariadb:10.5 mysql -h 127.0.0.1 -U root
```