# IP Proximity Network Visualizer

A 3D visualization tool that demonstrates IP proximity through an interactive network graph. Built with React, Golang, GraphQL, and MariaDB.

## Features

- Create and name virtual nodes based on your public IP
- Interactive 3D network visualization
- Drag to rotate the network view
- Zoom in/out functionality
- Real-time node updates

## Tech Stack

- Frontend: React + Three.js
- Backend: Golang + GraphQL
- Database: MariaDB

## Prerequisites

- Node.js 18+
- Go 1.21+
- Docker
- MariaDB

## Setup

1. Clone the repository
2. Start the database:
   ```bash
   docker-compose up -d
   ```

3. Set up the backend:
   ```bash
   cd backend
   go mod download
   go run main.go
   ```

4. Set up the frontend:
   ```bash
   cd frontend
   npm install
   npm start
   ```

5. Open http://localhost:3000 in your browser

## Project Structure

```
.
├── frontend/           # React application
├── backend/           # Golang server
├── docker-compose.yml # Database configuration
└── README.md
``` 