# Scrumer Backend

An opensource AI-powered project management assistant for software teams. This repository contains the backend services for Scrumer.

## About The Project

Scrumer Backend is the robust, scalable backend service for the Scrumer PM application. It is built with Go, leveraging PostgreSQL as its primary data store, GORM for object-relational mapping, and Gin for handling HTTP requests. The backend exposes a powerful GraphQL API, which serves as the data layer for the frontend applications.

Key responsibilities include:
- Providing a GraphQL API for data interaction
- Managing database operations and schema migrations
- Handling user authentication and authorization (future)
- Integrating with external services (future)

This project is built with Go, GORM, Gin, and GraphQL-Go, ensuring high performance and maintainability.

## Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/) (usually comes with Docker Desktop)
- [Go](https://golang.org/doc/install) (version 1.25 or higher, for development outside Docker if needed)

### Environment Variables

Create a `.env` file in the root of the project by copying `.env.example`. This file will store your database credentials.

```
cp .env.example .env
```

Edit the `.env` file with your desired values:

```
DB_NAME=scrumer
DB_USER=user
DB_PASSWORD=password
```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/BugMentor/scrumer-backend.git
   cd scrumer-backend
   ```
2. Build and run the Docker containers:
   ```sh
   docker-compose build
   ```

### Running the Application

First, ensure your database is up and configured:

```sh
docker-compose up -d db
```

#### Purging Data

To completely clear all data from the database:
```sh
docker-compose run --rm purge
```

#### Seeding Data

To populate the database with synthetic data for development/testing:
```sh
docker-compose run --rm seed
```

#### Starting the Backend

To start the GraphQL API server:
```sh
docker-compose up -d backend
```
The backend will be accessible at `http://localhost:8080/graphql`. You can use GraphiQL at `http://localhost:8080/graphql` for exploring the API.

To stop all services:
```sh
docker-compose down
```

## Features

- **GraphQL API:** A flexible and efficient API for data querying and manipulation.
- **PostgreSQL Database:** Robust and reliable data storage.
- **GORM ORM:** Simplifies database interactions in Go.
- **Gin Web Framework:** Fast and lightweight HTTP router.
- **Data Seeding:** Scripts to populate the database with synthetic data for development.
- **Data Purging:** Scripts to completely clear the database for clean testing environments.

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Donations

If you find Scrumer Backend useful and would like to support its continued development, please consider donating. Your contributions help us maintain and improve the project.

We recommend supporting us via **GitHub Sponsors**, which allows you to directly contribute to the project's maintainers and development efforts.

[Become a Sponsor on GitHub](https://github.com/sponsors/BugMentor)

Other platforms you might consider:
- [Open Collective](https://opencollective.com/bugmentor-arg/projects/scrumer)
- [Patreon](https://www.patreon.com/your_project)
- [Buy Me A Coffee](https://www.buymeacoffee.com/your_username)

## License

Distributed under the MIT License. See `LICENSE` for more information.

## 💰 Funding & Support

**Scrumer** is an open-source project developed by [BugMentor](https://bugmentor.com). We are dedicated to building a privacy-focused, vendor-lock-in-free alternative to Jira and Confluence.

Building and maintaining enterprise-grade agile tools takes significant resources. Your support directly funds server costs, development hours, and the maintenance of our open-source infrastructure.

### 🏆 Become a Sponsor (Open Collective)
This is the best way to support the project if you want public recognition on our README and website.

[![Open Collective](https://img.shields.io/opencollective/all/bugmentor-arg?label=Support%20Scrumer&logo=opencollective&color=blue)](https://opencollective.com/bugmentor-arg/projects/scrumer)

[**Click here to Donate via Open Collective**](https://opencollective.com/bugmentor-arg/projects/scrumer/donate)

---

### ⚡ Direct Support (Wise)
If you prefer to support the lead developer directly with lower fees (or for one-off contributions), you can scan the QR code or use the link below.

<a href="https://wise.com/pay/me/matiasm155">
  <img src="assets/img/wise-qr.jpg" width="200" alt="Scan to pay via Wise">
</a>
<br>

[**Send a Direct Contribution via Wise**](https://wise.com/pay/me/matiasm155)

---

### 💼 Commercial Support & Training
Need help implementing Scrumer in your company? BugMentor offers:
* **Enterprise Installation & Hosting Support**
* **QA & SDET Training**
* **Custom Feature Development**

[Contact us](https://bugmentor.com) for commercial inquiries.
