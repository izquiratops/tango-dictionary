# Tango

This website is a Japanese-English dictionary, it provides a simple way to look up Japanese words and their English translations.
While it's still in development the goal is to create a resource similar to Jisho.org. 

This project use the content provided by [JMDict Simplified](https://github.com/scriptin/jmdict-simplified). Without this valuable resource, building this dictionary would have been significantly more time-consuming.

```bash
$env:TANGO_VERSION="3.6.1"; $env:TANGO_REBUILD="true"; $env:TANGO_BLEVE_PATH="."; docker-compose up
```
