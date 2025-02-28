# Chatly

A real-time messaging platform designed for seamless and instant communication.

## 🚀 Quick Start

Ensure you have either [Docker](https://www.docker.com/get-started) or the following installed and configured:

- [Go](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/)
- [PostgreSQL](https://www.postgresql.org/download/)

### Clone the project

```bash
git clone https://github.com/danilovict2/chatly.git
cd chatly
```

### Set environment variables

```bash
cp .env.example .env
```

### Create a Pusher Account

1. Sign up for a [Pusher](https://pusher.com/) account.
2. Create an application on the Pusher dashboard and copy your App Keys.
3. Open the `.env` file and paste the keys into the corresponding fields.

### Run with Docker

```bash
docker compose up --build
```

### Run locally

1. Create the `real-time-chat` database in PostgreSQL.
2. Open the `.env` file located in the project root directory. Comment out the current `DATABASE_DSN` line and uncomment the alternative `DATABASE_DSN` line.
3. Install Frontend Dependencies

```bash
npm install
```

4. Build assets

```bash
make css
```

5. Run the Application

```bash
make run
```

### Access the Platform

Open your web browser and navigate to [http://127.0.0.1:8000](http://127.0.0.1:8000) to access the platform.


## ✨ Features

- 💬 Real-Time Chatting.
- 🔐 User authentication and authorization.
- 🛠️ Customizable User Profiles.
- ⚙️ Easy setup with Docker and local development options.
- 📡 Pusher API Integration.
- 🎨 Styling with Tailwind.

## 🤝 Contributing

### Build assets

```bash
npx tailwindcss -i ./views/css/app.css -o ./public/index.css
```

### Build the project

```bash
templ generate
go build -o bin/app main.go
```

### Run the project

```bash
./bin/app
```

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
