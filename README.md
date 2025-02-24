# Tango

To run the Docker container with the Bleve index folder as a runtime argument, you can use the -v option to mount the volume and the -e option to set the environment variable:

```bash
docker run -d -p 8080:8080 -e DB_VERSION=3.6.1 -v /custom/path/to/bleve:/root/database tango
```
