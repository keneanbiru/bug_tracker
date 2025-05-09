# Bug Tracker Application

A full-stack bug tracking application built with Vue.js and Go, featuring role-based access control and real-time bug management.

## Features

- 🔐 Role-based authentication (Admin, Developer, Manager)
- 🐛 Bug tracking and management
- 📊 Real-time status updates
- 🔍 Advanced filtering and search
- 📱 Responsive design
- ✅ Comprehensive test coverage

## System Workflow

### For Managers
- Create and assign bugs to developers
- Set bug priorities (High, Medium, Low)
- Track bug resolution progress
- Generate reports on bug statistics
- Monitor team performance
- Close or reopen bugs
- Add comments and attachments
- Filter bugs by status, priority, and assignee

### For Developers
- View assigned bugs
- Update bug status (New, In Progress, Resolved, Closed)
- Add technical details and solutions
- Upload screenshots or error logs
- Comment on bug progress
- Mark bugs as resolved
- Request additional information
- Filter bugs by status and priority

## Tech Stack

### Frontend
- Vue.js 3
- Pinia for state management
- Vue Router
- Jest for testing
- Tailwind CSS for styling

### Backend
- Go (Golang)
- MongoDB
- JWT Authentication
- RESTful API
- Clean Architecture

## Prerequisites

- Node.js (v16 or higher)
- Go (v1.21 or higher)
- MongoDB
- npm or yarn

## Installation

1. Clone the repository
```bash
git clone https://github.com/keneanbiru/bug-tracker.git
cd bug-tracker
```

2. Set up the backend
```bash
cd backend
go mod download
# Create .env file with your MongoDB connection string
go run main.go
```

3. Set up the frontend
```bash
cd bug-tracking
npm install
npm run dev
```

## Testing

### Frontend Tests
```bash
cd bug-tracking
npm test
```

### Backend Tests
```bash
cd backend
go test ./...
```

## Project Structure

```
bug-tracker/
├── backend/           # Go backend
│   ├── controller/    # HTTP controllers
│   ├── models/        # Data models
│   ├── repository/    # Database operations
│   └── usecase/       # Business logic
└── bug-tracking/      # Vue.js frontend
    ├── src/
    │   ├── components/
    │   ├── stores/
    │   └── views/
    └── tests/
```

## API Documentation

### Authentication Endpoints
- POST /api/auth/register - Register new user
- POST /api/auth/login - User login
- GET /api/auth/me - Get current user

### Bug Management Endpoints
- GET /api/bugs - List all bugs
- POST /api/bugs - Create new bug
- GET /api/bugs/:id - Get bug details
- PUT /api/bugs/:id - Update bug
- DELETE /api/bugs/:id - Delete bug

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Project Link: [https://github.com/keneanbiru/bug-tracker](https://github.com/keneanbiru/bug-tracker)
