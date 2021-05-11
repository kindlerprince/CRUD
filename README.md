# Simple Implementation of REST API with CRUD operation  


* Execute `postgress-install.sh` to install the postgres db and client  
* Connect to the postgres 
`psql -h localhost -U postgres`  
* Create the table using the `create.sql` file  
* Create the DB  

```
CREATE DATABASE login
\c login
\i create.sql
```

```
go build
./crud
```

## Test using the given curl Commands


* Registration  
`curl -X POST -d '{"name":"prince", "password":"kumar", "address":"chhota telpa", "email":"prince.kumar@gmail.com"}' http://localhost:8200/registration`  
* Login  
`curl -X POST -d '{ "password":"kumar","email":"prince.kumar@email.com"}' http://localhost:8200/login`  
* Forgot  
`curl -X PUT -d '{ "password":"kumar1","email":"prince.kumar@email.com"}' http://localhost:8200/forgot`  
* Delete  
`curl -X DELETE -d '{ "email":"prince.kumar@email.com"}' http://localhost:8200/user`  

