# gin-file-server
A simple file server built with Gin framework, featuring file uploads, downloads, and tracking of upload history in MySQL using GORM, as well as event notification using Kafka-go and logging with Zap.



## Development Plan

- [x]  Use the Gin module to write an HTTP server with two endpoints. Use a POST request to upload a file to the server, and a GET request to download a specific file from the server.
- [X] Use the GORM module to log the history of all file uploads to MySQL, and add an endpoint to query these records from the database.
- [X] Use the Segmentio/kafka-go module to push file upload events to Kafka.
- [X] Use the Zap module to log server running messages for easy troubleshooting and issue resolution.
- [ ] Build a simple CLI interface for user interaction, such as accepting user input and performing corresponding actions, displaying results, and other operations. Use the Go standard library or other third-party libraries to implement features such as menu selectors and input prompts.
