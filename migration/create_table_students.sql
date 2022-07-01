CREATE TABLE students(
    id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(40) NOT NULL,
    last_name VARCHAR(40) NOT NULL,
    age INT(3) NOT NULL,
    email VARCHAR(50) 
);

INSERT INTO students (first_name, last_name, age, email) VALUES("Dario", "Lopez", 26 , "daropl12@gmail.com" );