# Toronto Time API

A Go-based API that provides the current time in Toronto and logs requests to MySQL.

**This implementation**:
Creates an HTTP endpoint at /time
Returns both the current Toronto time and timezone information. It returns the current time in JSON format, in both ISO and user-friendly human-readable formats.
Stores each request time in a MySQL database
Properly handles the Toronto timezone

**Some additional considerations**:
Error handling is implemented.
Added environment variables for database credentials.
The database connection could be moved to a separate package for better organization
Added request logging middleware.
For production, we have also implemented CORS handling.

**Explanation**:

**CORS Middleware**: The corsMiddleware function sets the necessary CORS headers. It allows all origins (*), but you can specify a particular domain if needed for security reasons.

**Preflight Requests**: The middleware checks if the request method is OPTIONS, which is used for preflight requests in CORS. If so, it returns immediately after setting the headers.

**Integration**: The middleware is applied to the /time endpoint by wrapping the getTorontoTime handler.

This setup will allow your API to handle CORS requests properly. Adjust the Access-Control-Allow-Origin header to restrict access to specific domains if needed for your production environment.

## Pre-requisites

- Go and git are installed.
- Web browser or API testing tool (like Postman or cURL)
- Docker and docker-compose are installed and running
- Install MySQL workbench or MySQL-Shell.
- MySQL instance is running and you've connected to your local mysql database.
  - Ensure that your connection string is same as used in this code.
  - You may modify the connection string in the `main.go` file to match your DB's username and password.

## Features

- Returns current Toronto time in multiple formats (ISO and human-readable)
- Stores each request in MySQL database
- CORS enabled
- Dockerized application with Docker Compose
- Error handling and logging


## Project Structure

    bash
    goTimezoneApi/
    ├── Dockerfile
    ├── docker-compose.yml
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── README.md
    ├── schema.sql
    └── test.html

## Installation & Running

1. Clone the repository:
    `git clone https://github.com/ntanwir10/goTimezoneApi`
    `cd goTimezoneApi`

2. Start the application using Docker Compose:
    `docker-compose up --build`

3. Wait for the following messages:
   - "Server starting on :8080..."
   - MySQL initialization complete

## Development

### Environment Variables

The application uses these environment variables (set in docker-compose.yml):

- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`

### Making Changes

1. Stop the running containers
2. Make your changes
3. Rebuild and start:
    `docker-compose up --build`

## Testing the API

### Method 1: Using Web Browser

1. Open `test.html` in your web browser
2. You should see the current Toronto time displayed

### Method 2: Using cURL

    `curl http://localhost:8080/time`

### Method 3: Using Browser Console

1. Open your browser's developer tools (F12)
2. In the console, paste:
  `using JavaScript`

   `fetch('<http://localhost:8080/time>')`
    `.then(response => response.json())`
    `.then(data => console.log(data))`
    `.catch(error => console.error('Error:', error));`

### Method 4: Using Postman

1. Create a new GET request
2. Enter URL: `http://localhost:8080/time`
3. Send the request
   ![time api in Postman](image-1.png)

## Expected Response

The API returns JSON in this format:

    `    {
            "current_time": "2024-03-14T14:30:00-04:00",
            "readable_time": "2:30 PM",
            "readable_date": "Thursday, March 14",
            "timezone": "America/Toronto"
        }
    `

## Verifying Database Logs

1. Connect to MySQL container:
`docker exec -it goTimezoneApi-db-1 mysql -uroot -padmin`

2. Query the logs:
`USE timedb;
SELECT FROM time_logs ORDER BY id DESC LIMIT 5;`

## Stopping the Application

`docker-compose down`

To remove all data and start fresh: `docker-compose down -v`

## Troubleshooting

### Port Conflicts

If you see port conflict errors:

1. Ensure local MySQL is not running
2. Or modify ports in docker-compose.yml

### Database Connection Issues

1. Wait a few seconds after starting containers
2. Check if MySQL container is running: `docker ps`

### CORS Issues

If testing from a web page and getting CORS errors:

1. Ensure you're using the correct URL
2. Check browser console for specific error messages
