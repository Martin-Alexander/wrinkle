DROP DATABASE IF EXISTS e2e_test_db;
DROP USER IF EXISTS e2e_test_user;

CREATE DATABASE e2e_test_db;
CREATE USER e2e_test_user WITH ENCRYPTED PASSWORD 'e2e_test_password';
GRANT ALL PRIVILEGES ON DATABASE e2e_test_db TO e2e_test_user;
