to execute this project you need the files 'source.json' and 'target.json' in the same folder as the executable.
the location of these files can be changed by using '-source' and '-target'.

## Examples

### Help
```
./metadata-migration help
```
**Output:**
```
ask Ingo Rößner
some commands may accept parameters such as a list of ids
acceptable commands:
     all
     characteristics
     concepts
     device-types
     help
     protocols
```

### Migrate All
```
./metadata-migration all
```

### Migrate all Device-Types
```
./metadata-migration device-types
```

### Migrate some Device-Types
```
./metadata-migration device-types id1 id2 id3
```

### Export all
only source.json is needed. target.json wil be ignored.
```
./metadata-migration -export=testexport.json all
```

### Import

**simple import call:**
```
./metadata-migration import import_file.json
```

**export may be used:**
```
./metadata-migration -export=export_file.json import import_file.json
```

**import file may be online location:**
```
./metadata-migration import https://raw.githubusercontent.com/SENERGY-Platform/metadata-export/master/export.json
```


**only protocols:**

_the parameters after the import file location may be used to select resources from the file. each resource which path contains any of the given additional parameters will be send._
```
./metadata-migration import import_file.json protocols
```


**only 2 protocols by id:**
_the parameters after the import file location may be used to select resources from the file. each resource which path contains any of the given additional parameters will be send._

```
./metadata-migration import import_file.json urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e urn:infai:ses:protocol:c9a06d44-0cd0-465b-b0d9-560d604057a2
```

## Transformer
use transformations before writing to target
```
./metadata-migration -transformer=addserviceinteraction -export=testexport.json all
```

## Config Fields

| Field               | Description                  |
|---------------------|------------------------------|
| senergy_user        | user name from keycloak      |
| password            | password from keycloak       |
| auth_client         | keycloak auth client         |
| auth_client_secret  | keycloak autch client secret |
| auth_url            | url to keycloak              |
| device_manager_url  | url to device-manager        |
| source_list_url     | url to permission-search     |
| source_semantic_url | url to semantic-repository   |

**JSON:**
```
{
  "senergy_user":"user",
  "password":"password",
  "auth_client":"keycloak auth client",
  "auth_client_secret":"keycloak autch client secret",
  "auth_url":"url to keycloak",
  "device_manager_url":"url to device-manager",
  "source_list_url":"url to permission-search",
  "source_semantic_url":"url to semantic-repository"
}
```