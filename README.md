User Authentication System (Golang + Docker + PostgreSQL)

A simple user authentication system featuring signup and login APIs, built with Golang, PostgreSQL, and containerized using Docker.
Ideal for home assignments or as a starter template.

🛠 Technical Overview

Language & Framework: Golang using net/http from the standard library.

Database: PostgreSQL, storing user data securely with SQL.

Authentication: Passwords hashed via bcrypt; JWT tokens used for secure session management.

Containerization: Docker ensures consistent development environments and easy deployment.

Configuration: Managed through environment variables.

⚙️ Setup & Run Instructions
1️⃣ Clone the Repo
git clone https://github.com/hari130303/sign-in-sign-up.git
cd sign-in-sign-up

2️⃣ Database Setup (Local PostgreSQL)

Before starting the application, create the database and table in PostgreSQL:

CREATE DATABASE task_db;

\c task_db;

CREATE TABLE public.user_master (
    user_id serial NOT NULL,
    user_name character varying(30) NOT NULL,
    password text NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    mail_id character varying(255) NOT NULL,
    PRIMARY KEY (user_id)
);

3️⃣ Environment Variables

Create a .env file in the project root:

DBUSER=postgres
DBPASS=12345
DBNAME=task_db
DBHOST=host.docker.internal
-- DBHOST=localhost   # Use this if running locally without Docker
DBPORT=5432
JWTSECRETKEY=secretkey1234567

4️⃣ Run Locally (Optional)
go mod tidy
go run main.go


Note: If you want to run locally, uncomment the following code block in main.go to load environment variables from .env:

err = godotenv.Load()
if err != nil {
    log.Fatal("Error loading .env file")
}


The server will be running at:
👉 http://localhost:8088

5️⃣ Run with Docker
Build the Docker Image
docker build -t go-auth-app .

Run with Docker Compose
docker compose up -d


In docker-compose.yml, the app is mapped to port 8088 on the host, which forwards to 8088 inside the container.

🔑 API Endpoints
1️⃣ Signup

POST http://localhost:8088/user/register

Request
{
  "username": "testuser1",
  "password": "12345",
  "mail_id": "hari@gmail.com"
}

Response
{
  "id": 3,
  "name": "testuser1",
  "email": "hari@gmail.com",
  "timestamp": "2025-09-02T01:40:52.772510518Z"
}

2️⃣ Login

POST http://localhost:8088/login

Request
{
  "password": "12345",
  "mail_id": "hari@gmail.com"
}

Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}


📜 License

This project was developed for educational purposes as part of a home assignment.