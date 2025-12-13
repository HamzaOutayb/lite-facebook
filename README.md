# Social Network Project

A **Facebook-like social network** built with a modern JavaScript frontend and a Go backend. This project demonstrates full-stack web development concepts including authentication, real-time communication, database migrations, Docker containerization, and scalable architecture.

---

## 🚀 Features

### 👤 Authentication

* User registration & login
* Session-based authentication using **cookies**
* Secure password hashing with **bcrypt**
* Persistent login until explicit logout

**Registration fields:**

* Email (required)
* Password (required)
* First Name (required)
* Last Name (required)
* Date of Birth (required)
* Avatar/Image (optional)
* Nickname (optional)
* About Me (optional)

---

### 👥 Followers

* Send follow requests
* Accept or decline follow requests
* Automatic follow for **public profiles**
* Unfollow users

---

### 🧑 Profile

* Public & private profiles
* Toggle profile visibility (public / private)
* Display:

  * User information (except password)
  * User activity
  * User posts
  * Followers & following lists

---

### 📝 Posts & Comments

* Create posts and comments
* Attach images or GIFs (JPEG, PNG, GIF)
* Post privacy levels:

  * **Public** – visible to all users
  * **Followers only** – visible to followers
  * **Custom** – visible to selected followers only

---

### 👨‍👩‍👧 Groups

* Create groups with title & description
* Invite users to join groups
* Accept or decline group invitations
* Request to join groups (approval by creator)
* Browse all groups
* Group-only posts and comments

#### 📅 Group Events

* Create events inside groups
* Event fields:

  * Title
  * Description
  * Date & time
* Event responses:

  * Going
  * Not Going

---

### 💬 Chat (WebSockets)

* Private 1-to-1 messaging
* Messaging allowed if at least one user follows the other
* Real-time delivery via **WebSockets**
* Emoji support

#### 🗨 Group Chat

* Each group has a shared chat room
* Accessible only to group members

---

### 🔔 Notifications

Notifications are visible on **every page** and are **separate from chat messages**.

Users are notified when:

* They receive a follow request (private profile)
* They receive a group invitation
* A user requests to join their group (creator only)
* A new event is created in a group they belong to
* Additional custom notifications are supported

---

## 🛠 Tech Stack

### Frontend

* JavaScript Framework (one of):

  * Next.js
* HTML, CSS
* Responsive & performance-focused UI
* Communicates with backend via HTTP & WebSockets

---

### Backend

* **Go (Golang)**
* RESTful API
* Session & cookie-based authentication
* WebSocket support for real-time chat
* Image upload & storage

---

### Database

* **SQLite**
* Structured with an Entity Relationship Diagram (ERD)
* Optimized queries and indexing

---

### 🔄 Migrations

* Managed using **golang-migrate** (or equivalent)
* Automatically applied on application startup

**Example structure:**

```
backend/
 └── pkg/
     └── db/
         ├── sqlite.go
         └── migrations/
             └── sqlite/
```

---

### 🐳 Docker

The project is fully containerized with **two Docker images**:

#### Backend Container

* Runs the Go server
* Exposes API & WebSocket ports
* Connects to SQLite database

#### Frontend Container

* Serves frontend assets
* Communicates with backend over HTTP

---

## 📦 Allowed Packages

* Go standard library
* gorilla/websocket
* golang-migrate
* sqlite3
* bcrypt
* UUID
* sql-migration / migration

---

## 📚 What You Will Learn

* Session & cookie-based authentication
* Docker & containerization
* Backend–Frontend communication
* WebSockets for real-time apps
* SQL & database design
* Database migrations
* Secure password handling
* Image uploads & storage
* Building scalable social platforms

---

## 🧪 Getting Started

### Prerequisites

* Docker & Docker Compose
* Go (for local development)
* Node.js (for frontend development)

### Run with Docker

```bash
docker-compose up --build
```

---

## 📌 Notes

* This project focuses on **learning and architecture**, not UI perfection
* Security best practices are applied where possible
* Additional features & notifications are encouraged

---

## 🧠 Something Wrong?

If something doesn’t work:

* Check Docker logs
* Verify database migrations
* Ensure WebSocket connections are open
* Confirm cookies are enabled in the browser

---

## 👨‍💻 Author

Built as a full-stack learning project to explore modern web technologies with Go and JavaScript.

---

Happy coding 🚀
