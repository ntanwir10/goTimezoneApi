# goTimezoneApi

This implementation:
Creates an HTTP endpoint at /time
Returns the current Toronto time in JSON format
Stores each request time in a MySQL database
Properly handles the Toronto timezone
Returns both the current time and timezone information

The API will return JSON in this format:
`
{
    "current_time": "2024-03-14T14:30:00-04:00",
    "timezone": "America/Toronto"
}
`

Then we can test it by making a GET request to `http://localhost:8080/time`
Some additional considerations:
Error handling is implemented but could be more robust for production use
We might want to add environment variables for database credentials
The database connection could be moved to a separate package for better organization
We might want to add request logging middleware

For production, we have also implemented CORS handling

Explanation:
CORS Middleware: The corsMiddleware function sets the necessary CORS headers. It allows all origins (*), but you can specify a particular domain if needed for security reasons.

Preflight Requests: The middleware checks if the request method is OPTIONS, which is used for preflight requests in CORS. If so, it returns immediately after setting the headers.

Integration: The middleware is applied to the /time endpoint by wrapping the getTorontoTime handler.

This setup will allow your API to handle CORS requests properly. Adjust the Access-Control-Allow-Origin header to restrict access to specific domains if needed for your production environment.
