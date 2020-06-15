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
