# Go Event Booking API

Ein modulares **Event-Buchungssystem in Go**, bestehend aus mehreren Microservices.  
Derzeit umfasst das Projekt folgende Services:

- **User-Service** – Authentifizierung und Registrierung  
- **Event-Service** – Verwaltung von Events und Teilnehmer-Buchungen  

Jeder Service läuft als eigenständiger Container im Docker-Compose-Setup und nutzt eine eigene SQLite-Datenbank (PostgreSQL-Unterstützung ist geplant).

---

## Features

- Saubere Service-Struktur in Go mit `gin-gonic`
- Gemeinsame `.env`-Konfiguration (über `.env.example`)
- Multi-Service-Setup mit **Docker Compose**
- Bereit für zukünftige **Kubernetes-Deployments**
- Optionale **Swagger-Dokumentation** (in Planung)

---

## Installation & Setup

1. **Repository klonen**
   ```
   git clone https://github.com/rearatrox/go-event-booking-api.git
   cd go-event-booking-api
   ```

2. **.env-Dateien anpassen**  
   Erstelle aus der `.env.example` eine `.env`-Datei und passe sie an:
   ```
   cp .env.example user-service/.env
   cp .env.example event-service/.env
   ```

3. **Container starten**
   ```
   docker compose up -d
   ```

4. **Services testen**   
   - Event-Service: [http://localhost:8081 (oderer anderer gewählter Port)](http://localhost:8081)
   - User-Service: [http://localhost:8082 (oderer anderer gewählter Port)](http://localhost:8082) 

---

## ⚙️ Umgebungsvariablen (`.env.example`)

| Variable | Beschreibung | Beispielwert |
|-----------|---------------|---------------|
| **API_PREFIX** | Gemeinsamer API-Präfix für alle Services | `/api/v1` |
| **JWT_SECRET** | Geheimschlüssel für JWT-Token-Signierung | `supersecret` |

### 🪵 Logger

| Variable | Beschreibung | Beispielwert |
|-----------|---------------|---------------|
| **LOG_LEVEL** | Log-Level (z. B. `debug`, `info`, `warn`, `error`) | `info` |
| **LOG_FORMAT** | Format der Logs (`text` oder `json`) | `json` |
| **LOG_OUTPUT** | Zielausgabe der Logs (`stdout`, `file`, etc.) | `stdout` |
| **REQUEST_ID_HEADER** | Header-Name für Request-IDs (Tracing) | `X-Request-Id` |

### 🧩 Services

| Variable | Beschreibung | Beispielwert |
|-----------|---------------|---------------|
| **EVENTSERVICE_PORT** | Externer Port des Event-Service | `8081` |
| **USERSERVICE_PORT** | Externer Port des User-Service | `8082` |

### 🗄️ Datenbank

| Variable | Beschreibung | Beispielwert |
|-----------|---------------|---------------|
| **DB_USERNAME** | Benutzername für PostgreSQL | `admin` |
| **DB_PASSWORD** | Passwort für PostgreSQL | `password123` |
| **DB_NAME** | Name der Datenbank | `api_db` |
| **DB_PORT** | Port der PostgreSQL-Instanz | `5432` |
| **DB_SSLMODE** | SSL-Modus der Verbindung (`disable`, `require`, etc.) | `disable` |

💡 **Hinweis:**  
Die DATABASE_URL wird automatisch mit den obigen Angaben generiert

---

## Swagger

Eine Integration mit [**swaggo/gin-swagger**](https://github.com/swaggo/gin-swagger) ist vorgesehen.  

---

## Kubernetes 

In Zukunft werden Kubernetes-Manifeste unter  
`/k8s/` bereitgestellt, um eine einfache Bereitstellung der Services auf einem Cluster zu ermöglichen.

---

## Grundlage 

Als Grundlage des Projekts diente der folgende Udemy-Kurs: [Go - The Complete Guide](https://www.udemy.com/course/go-the-complete-guide/)
