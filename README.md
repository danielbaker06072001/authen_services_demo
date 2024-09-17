# Authentication project -  GoLang
Feature - API
    - Login
    - Register
    - ForgotPassword
    - Role Initialize

Structure:
    - This project use GIN for server hosting 
    - This project using CRUD method request 
    - This project also implement token using JWT 

# What is the structure of this project
- Proto file: Proto file that defined and generate the services (along side with the data that we wanted to return)
- pb file: Generated file from prpto file 

    - DTO + Models: structure of data / return message 
    - Infranstructure: Initialize all database here in this file 
    - Repositories: Mostly function that will interact with the database (insert, update, etc) 
    - Handler: Go file that handle main logic (Login, Register, ForgotPassword)
    - Services: 
            - Mostly services that performing the logic method (mostly comparing, checking, etc) but not interact with the databse yet


# How to understand this login project
- Step 1: Go from Handler file (Server.go) -> Go through other handler file
- Step 2: From handler, go through Services + Repositories file , also check the DTO + models file for the data structure
    - Step 3: In the services and repositories file , they will call  other method from services 
    - Step 4: All of the method will return response based on a Model structure
