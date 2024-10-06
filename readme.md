# E-commerce Application

## Overview

This project is a robust and scalable e-commerce application built using Go and Docker. It provides essential features such as user authentication, product management, order processing, and integration with a MySQL database.

## Table of Contents

- [E-commerce Application](#e-commerce-application)
  - [Overview](#overview)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Technologies Used](#technologies-used)
  - [Installation](#installation)

## Features

- User authentication with JWT.
- Product management including create, read, update, and delete (CRUD) operations.
- Shopping cart functionality.
- Order management system.
- Integration with MySQL for data persistence.

## Technologies Used

- **Go**: 1.22.3
- **Docker**: Latest
- **Mongodb**: Latest
- **MySQL**: Latest
- **Alpine Linux**: For lightweight containers
- **Postman**: For API testing

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/nurcahyaari/edot.git
   cd ecommerce-app
2. Run the application:
   ```bash
   docker-compose up --build
   ```
3. Import postman collection
   Postman Collection
    To test the API, you can use the Postman collection included in the project. Import the collection file (postman_collection.json) into Postman to start exploring the API endpoints.
    - Open Postman.
    - Click on "Import".
    - Select the `ecommerce.postman_collection.json` file from the project directory.
    - Use the collection to test the API endpoints.