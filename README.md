# URL Shortener Service

## Overview
This is a URL shortening service built with Go, Gin, and MongoDB. It provides a RESTful API to create, retrieve, update, and delete short URLs.

## Features
- Create short URLs
- Redirect to original URLs
- Update existing short URLs
- Delete short URLs
- Track URL access statistics

## Prerequisites
- Go 1.16+
- MongoDB 4.0+
- Docker (optional)

## Installation

### Local Setup
1. Clone the repository
```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener
```

2. Install dependencies
```bash
go mod download
```

3. Set MongoDB URI (optional)
```bash
export MONGODB_URI=mongodb://localhost:27017
```

4. Run the application
```bash
go run cmd/main.go
```

### Docker Setup
```bash
# Build the image
docker build -t url-shortener .

# Run the container
docker run -p 8080:8080 -e MONGODB_URI=mongodb://host.docker.internal:27017 url-shortener
```

## API Endpoints

### Create Short URL
`POST /api/v1/shorten`
```json
{
  "url": "https://www.example.com/some/long/url"
}
```

### Retrieve Original URL
`GET /api/v1/shorten/{shortCode}`

### Update Short URL
`PUT /api/v1/shorten/{shortCode}`
```json
{
  "url": "https://www.example.com/updated/url"
}
```

### Delete Short URL
`DELETE /api/v1/shorten/{shortCode}`

### Get URL Statistics
`GET /api/v1/shorten/{shortCode}/stats`

## Environment Variables
- `MONGODB_URI`: MongoDB connection string
- `PORT`: Server port (default: 8080)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss proposed changes.

## License
[MIT](https://choosealicense.com/licenses/mit/)