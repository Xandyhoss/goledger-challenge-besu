
# Goledger-Challenge-Besu

This application provides an interface to interact with a smart contract running on a **Hyperledger Besu** network. Built using the **GoFiber** framework, it handles database interactions with **PostgreSQL** using **GORM**, an ORM library for GoLang. The database will be running in a **Docker** container.

## Project Structure
The project is located in the `/besu/app` folder and follows this structure:
```
/app
---/controllers
---/db
---/routes
---/services
```

**Controllers**
Handles the logic for the API routes.

**DB**
Manages the database connection and models using GORM.

**Routes**
Defines the available API endpoints.

**Services**
Contains the core business logic, including interactions with the smart contract. 

## Setup Instructions

1. **Navigate to the Project Folder:**
```bash
cd /besu
```
2. **Install dependencies**
```bash
go mod tidy
```
3. **Deploy the Smart Contract:** Run the following command and copy the **contract address** displayed at the and of the execution.
```bash
./startDev.sh
```
4. **Set Up Environment Variables:**
* Copy the contents of `.env.example` to a new `.env` file:
```bash
cp .env.example .env
```
* Update the `CONTRACT_ADDRESS` key in the `.env` file with the contract address from step 3.
* Update the `PRIVATE_KEY` key with the one of the private keys from the `/genesis/genesis.json` file.
5. **Start the PostgreSQL Container**
```bash
./startDb.sh
```
6. **Run the application:**
```bash
go run main.go
```

---

## API Endpoints

1. **GET /get**
Retrieves the current value stored in the smart contract.

**Response example:**
```json
{
	"result": 65,
	"status": "success"
}
```
--- 

1. **POST /set**
Set a new value into the smart contract

**Request body example:**
```json
{
	"value": 30,
}
```
**Response example:**
```json
{
	"result": 65,
	"status": "success"
}
```
--- 

3. **GET /check**
Verifies if the value in the PostgreSQL database matches the value in the smart contract

**Response example:**
```json
{
	"result": true // or false,
}
```
--- 

4. **PUT /sync**
Synchronizes the value in the PostgreSQL database with the value from the smart contract.

**Response example:**
```json
{
	"status": "contract synced",
}
```

## Notes
* Ensure the `.env` file is properly configured before running the project.
* A Insomnia request collection is contained inside the `/besu` folder as `Insomnia_Collection.json`