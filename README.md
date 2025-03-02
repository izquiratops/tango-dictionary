# Tango

This website is a Japanese-English dictionary, it provides a simple way to look up Japanese words and their English translations.
While it's still in development the goal is to create a resource similar to Jisho.org. 

This project use the content provided by [JMDict Simplified](https://github.com/scriptin/jmdict-simplified). Without this valuable resource, building this dictionary would have been significantly more time-consuming.

To run the Docker container with the Bleve index folder as a runtime argument, you can use:
- The -v option to mount the folder where Bleve indexes are located into the the volume
- The -e option to set the environment variable that selects the Database version

```bash
docker run -d -p 8080:8080 -e DB_VERSION=3.6.1 -v /custom/path/to/bleve:/root tango -- --rebuild-database=true
```
