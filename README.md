# Stock Trading App

## Overview

This project is a stock trading application that allows users to trade stocks, track their portfolio, and manage their accounts. The application is built using modern technologies and follows a microservices architecture.

## Architecture

The application consists of the following components:

1. **Mobile App (React Native)**
   - User interface for trading, portfolio tracking, and account management
   - Communicates with the backend API

2. **Backend API (Golang)**
   - Handles business logic, authentication, and API requests
   - Connects to the database and IBKR API

3. **Database (PostgreSQL)**
   - Stores user accounts, portfolios, transactions, and market data

4. **IBKR API Connector**
   - Handles integration with the Interactive Brokers API
   - Manages market data streaming and order execution

## Getting Started

### Prerequisites

- Node.js (for React Native)
- Go (for backend API)
- PostgreSQL (for database)
- Interactive Brokers account (for market data and trading)

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/thetom42/trading.git
   cd trading
   ```

2. Install dependencies for the mobile app:
   ```sh
   cd mobile
   npm install
   ```

3. Install dependencies for the backend API:
   ```sh
   cd ../backend
   go mod tidy
   ```

4. Set up the database:
   ```sh
   cd ../db
   psql -f setup.sql
   ```

5. Configure the IBKR API connector:
   - Follow the instructions in the `ibkr` directory to set up the connector.

### Running the Application

1. Start the backend API:
   ```sh
   cd backend
   go run main.go
   ```

2. Start the mobile app:
   ```sh
   cd mobile
   npm start
   ```

3. Access the application:
   - Open the mobile app on your device or emulator.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.