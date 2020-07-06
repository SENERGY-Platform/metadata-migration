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