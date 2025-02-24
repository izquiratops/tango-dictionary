# Tango

To run the Docker container with the Bleve index folder as a runtime argument, you can use:
- The -v option to mount the folder where Bleve indexes are located into the the volume
- The -e option to set the environment variable that selects the Database version

```bash
docker run -d -p 8080:8080 -e DB_VERSION=3.6.1 -v /custom/path/to/bleve:/root/database tango
```
