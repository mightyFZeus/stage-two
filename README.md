# stage-two



# Countries API

A Go backend service for managing country data with features like:

Add/refresh countries

Fetch countries with filters and sorting

Fetch a single country by name

Delete a country by name

Generate a summary image (cache/summary.png)

Status endpoint showing total countries and last refresh timestamp


# Prerequisites

Go 1.21+ installed

MySQL 

gg library for image generation

go get github.com/fogleman/gg


# Other dependencies
go get github.com/go-chi/chi/v5

go get gorm.io/gorm

go get gorm.io/driver/postgres

go get github.com/google/uuid

go install github.com/air-verse/air@latest

Run Locally

Clone the repository

git clone <repo_url>
cd <repo_folder>


Install dependencies

go mod tidy


# Run the server
air

# Run docker
docker compose up --build

