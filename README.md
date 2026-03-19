# 🐊 Gator

Gator is a CLI-based RSS feed aggregator built in Go. It allows users to register, follow feeds, and continuously fetch updates from RSS sources.

---

## 🚀 Requirements

Before running Gator, make sure you have:

- Go (>= 1.25)
- PostgreSQL

### Install Go
https://go.dev/doc/install

### Install PostgreSQL
https://www.postgresql.org/download/

---

## 📦 Installation

Install the CLI using:

```bash
go install github.com/jahidul39306/gator@latest
```

Make sure your `$GOPATH/bin` is in your `PATH`, then run:

```bash
gator
```

---

## ⚙️ Configuration

Create a config file:

```bash
~/.gatorconfig.json
```

Example:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

---

## 🗄️ Database Setup

Create database:

```sql
CREATE DATABASE gator;
```

Run schema files from:

```
sql/schema/
```

---

## 🧪 Running the Program

### Development

```bash
go run .
```

### Production

```bash
go build -o gator
./gator
```

Or if installed:

```bash
gator
```

> Go programs are statically compiled binaries. After building, you don’t need Go installed to run them.

---

## 📜 Commands

### 👤 User Commands

Register:
```bash
gator register <username>
```

Login:
```bash
gator login <username>
```

List users:
```bash
gator users
```

Reset users:
```bash
gator reset
```

---

### 📰 Feed Commands

Add feed:
```bash
gator addfeed <name> <url>
```

List feeds:
```bash
gator feeds
```

Follow feed:
```bash
gator follow <url>
```

Unfollow feed:
```bash
gator unfollow <url>
```

Show following:
```bash
gator following
```

---

### 🔄 Aggregation

Run feed scraper:

```bash
gator agg <duration>
```

Example:

```bash
gator agg 10s
```

---

## 🧠 How It Works

- Uses PostgreSQL to store users and feeds
- Fetches RSS feeds over HTTP
- Parses XML responses
- Periodically polls feeds

---

## 📁 Project Structure

```
internal/
  config/
  database/
sql/
  queries/
  schema/
main.go
```

---

## 🌐 Repository

https://github.com/jahidul39306/gator