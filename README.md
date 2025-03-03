# Tango

This website is a Japanese-English dictionary, it provides a simple way to look up Japanese words and their English translations.
While it's still in development the goal is to create a resource similar to Jisho.org. 

This project use the content provided by [JMDict Simplified](https://github.com/scriptin/jmdict-simplified). Without this valuable resource, building this dictionary would have been significantly more time-consuming.

## Env variables

- **TANGO_VERSION**: This is the version number of JMDict and must follow the semantic versioning rules.

- **TANGO_REBUILD**: When set to true will clear any existing data matching TANGO_VERSION and import a fresh dictionary from scratch. When enabled, the importer looks for the JSON source file `${TANGO_JMDICT_SOURCE_PATH}/jmdict-eng-{$TANGO_VERSION}.json`.

- **TANGO_JMDICT_SOURCE_PATH**: Specifies the directory where Bleve Indexes and JMDICT JSONs are located.

- **TANGO_LOCAL**: A flag to set the MongoDB connection to a local instance.

### Bash
```bash
export TANGO_VERSION="3.6.1" && export TANGO_REBUILD="true" && export TANGO_JMDICT_SOURCE_PATH="./jmdict_source" && docker compose up
```

### PowerShell
```powershell
$env:TANGO_VERSION="3.6.1"; $env:TANGO_REBUILD="true"; $env:TANGO_JMDICT_SOURCE_PATH="./jmdict_source"; docker-compose up
```
