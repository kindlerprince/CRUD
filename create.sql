-- CREATING DATABASE --
-- CREATE DATABASE login;
-- CREATING TABLE 
DROP TABLE UserDetails;
CREATE TABLE UserDetails (
    name varchar (200),
    email varchar(255) ,
    password varchar(255) NOT NULL,
    address varchar(255),
    PRIMARY KEY(email)
);
