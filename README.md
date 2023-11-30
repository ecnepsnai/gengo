# GenGo

GenGo is yet another Golang code generator, _because there aren't enough of those already._

GenGo supports multiple different modules:

- Implementations of:
    - github.com/ecnepsnai/ds
    - github.com/ecnepsnai/stats
    - github.com/ecnepsnai/store
- Directory & filesystem utilities
- Enums, including TypeScript definitions
- Gob registration
- State storage system
- Statistics system

## Usage

```
Usage: gengo [Options]
-n --name <name>           Package name, defaults to 'main'
-c --config-dir <dir>      Config dir, defaults to current dir
-g --go-output-dir <dir>   Output dir for go files, defaults to current dir
-t --ts-output-dir <dir>   Output dir for ts files, defaults to current dir
-q --quiet                 Don't print out names of generated files
```

Ensure that your $GOBIN is in your $PATH and run `gengo` in the directory where your configuration JSON files are.

## Configuration

Aside from the command-line parameters, a configuration file `gengo.json` or `gengo.yaml` may be provided in the configuration directory to control the generation and output process.

|Property|Type|Description|
|-|-|-|
|`minimum_version`|String|The minimum supported version of gengo for this project in the format of `v1.2.3`. If omitted no version check is performed.|
|`file_prefix`|String|A prefix to attach to all generated files. Defaults to `gengo_`. Set to an empty string to remove the prefix.|

## Modules

### Data Store

Data Stores provide an implementation of the github.com/ecnepsnai/ds package.

#### Config

Options go into the `data_store.json` or `data_store.yaml` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of the data store iteself, as seen within code|
|`object`|String|The golang object to register to the store|
|`unordered`|Bool|Optional, if the table should not be ordered. Defaults to false.|

##### Example

```json
[
    {
        "name": "User",
        "object": "User",
        "unordered": true
    }
]
```

#### Output

- `func dataStoreSetup(storageDir string)`: Method to set up all data stores.
- `func dataStoreTeardown()`: Method to tear down all data stores
- `type <name_lowercase>StoreObject struct{ table *ds.Table }`: Type associated with the store
- `var <name_titlecase>Store: <name_lowercase>StoreObject`: Global variable associated with the store, set up when `dataStoreSetup()` is called

### Directories

Directories provides utilities for data directory management.

#### Config

Options go into the `directory.json` or `directory.yaml` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of the directory as seen within code|
|`is_data`|Bool|If this data is a data directory|
|`dir_name`|String|The name of the directory itself as seen on the filesystem|
|`required`|Bool|Optional, if this directory must already exist. If false, the directory will be created for you|
|`subdirs`|Array|Optional subdirectories under this parent|

##### Example

```json
[
    {
        "name": "Data",
        "is_data": true,
        "dir_name": "data",
        "subdirs": [
            {
                "name": "Cache",
                "dir_name": "cache"
            }
        ]
    },
]
```

#### Output

- `var operatingDirectory: String`: Global variable of the operating directory of the application
- `var dataDirectory: String`: Global variable of the data directory of the application
- `var Directories: apiDirectories`: Global variable of absolute paths to your directories
- `func fsSetup()`: Method to set up directories
- `func DirectoryExists(directoryPath string) bool`: Method to check if a directory exists
- `func MakeDirectoryIfNotExist(directoryPath string)`: Method to create a directory if it does not already exist
- `func FileExists(filePath string) bool`: Method to check if a file exists

### Enums

Enums are enums, the one thing that go really needs.

#### Config

Options go into the `enum.json` or `enum.yaml` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of the enum. This is prefixed to all associated variables and methods for this enum|
|`type`|String|The golang type for this enum|
|`description`|String|Optionally provide a description for the enum overall|
|`include_typescript`|Bool|Optionally export a typescript definition of this enum. Only supported on go primitive types|
|`values`|Array|The possible values for this enum|
|`values.key`|String|The key for this value|
|`values.description`|String|Optionally provide a description for this value|
|`values.value`|String|The value of this enum as it would be seen in Golang code. I.E. if it's a string include quotation marks|

##### Example

```json
[
    {
        "name": "Fruit",
        "type": "string",
        "values": [
            {
                "key": "Apple",
                "description": "a delicious granny smith",
                "value": "\"Apple\""
            },
            {
                "key": "Orange",
                "description": "a tasty mandarin",
                "value": "\"Orange\""
            },
            {
                "key": "Banana",
                "description": "just for scale",
                "value": "\"Banana\""
            }
        ]
    }
]
```

_Note how the values include quotes._

#### Output

##### Go

- `const ( <enum name><value key> = <value> ...)`: A const of all possible enum values for this enum.
- `var All<enum name> = []<enum type>`: An array of all enum values.
- `var <enum name>Map = map[<enum type>]<enum type>`: A map of enum keys to enum values.
- `func Is<enum name>(q <enum type>) bool`: A method to validate that the given value is a valid enum value.
- `func ForEach<enum name>(m func(value <enum type>))`: A convience method to iterate over each enum value

##### TypeScript

- `export enum <enum name>`: An exported enum definition for this enum.
- `export function <enum name>All()`: An exported function that returns an array of all enum values.
- `export function <enum name>Config()`: An exported function that returns an array of objects that describe the enum.

### Gobs

Gobs provide an easy way to register many objects with gob in one place, and to prevent a panic if they're already registered.

#### Config

Options go into the `gob.json` or `gob.yaml` file.

|Property|Type|Description|
|-|-|-|
|`type`|String|The golang type to register, including the `{}`|
|`import`|String|The import path for this type. Imports are deduplicated|

##### Example

```json
[
    {
        "type": "time.Time{}",
        "import": "time"
    }
]
```

#### Output

- `func gobSetup()`: Method to register all of your gob types

### State Store

State Store is a simple store interface for persisting configuration or other properties.

#### Config

Options go into the `state.json` or `state.yaml` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of this state property|
|`type`|String|The golang type for this property|
|`default`|String|The default value of this property as it would be seen in Golang code. I.E. if it's a string include quotation marks|
|`import`|String|The import path for this type. Imports are deduplicated|

##### Example

```json
[
    {
        "name": "TableVersion",
        "type": "int",
        "default": "0"
    }
]
```

#### Output

- `var State *gengoStateObject`: The global state object
- `func stateSetup(storageDir string)`: Method to set up the state store
- `func (s *gengoStateObject) Close()`: Method to close the state store
- `func (s *gengoStateObject) Get<property name>() <property type>`: Method to get the property value, or the default value. Threadsafe.
- `func (s *gengoStateObject) Set<property name>(value <property type>)`: Method to set the property value. Threadsafe.

### Stats

Stats provides an interface to the github.com/ecnepsnai/stats package.

#### Config

Options go into the `stats.json` or `stats.yaml` file.

|Property|Type|Description|
|-|-|-|
|`counters`|Array|The counters|
|`timed_counters`|Array|The timed counters|
|`timers`|Array|The timers|
|`counters.name`|String|The name of this counter|
|`counters.description`|String|The description of this counter|
|`timed_counters.name`|String|The name of this counter|
|`timed_counters.description`|String|The description of this counter|
|`timed_counters.max_minutes`|String|The maximum number of minutes for samples to be retained|
|`timers.name`|String|The name of this timer|
|`timers.description`|String|The description of this timer|

##### Example

```json
{
    "counters": [
        {
            "name": "NumberUsers",
            "description": "The number of users"
        }
    ],
    "timed_counters": [
        {
            "name": "FailedRequests",
            "description": "The number of failed requests",
            "max_minutes": 1440
        }
    ],
    "timers": [
        {
            "name": "BuildDuration",
            "description": "Duration of builds"
        }
    ]
}
```

#### Output

- `type gengoStatsCounters struct { <counter name> *stats.Counter ... }`: All of your counter
- `type gengoStatsTimedCounters struct { <counter name> *stats.TimedCounter ... }`: All of your timed counters
- `type gengoStatsTimers struct { <timer name> *stats.Timer ... }`: All of your timers
- `var Stats *gengoStatsObject`: The global stats object
- `func statsSetup()`: Method to set up the stats interface
- `func (s *gengoStatsObject) Reset()`: Reset all counters to their default values
- `func (s *gengoStatsObject) GetCounterValues()`: Get a map of all counter values
- `func (s *gengoStatsObject) GetTimedCounterValues()`: Get a map of all timed counter values
- `func (s *gengoStatsObject) GetTimedCounterValuesFrom(d time.Duration)`: Get a map of all timed counter values from d
- `func (s *gengoStatsObject) GetTimerAverages()`: Get a map of all timer averages
- `func (s *gengoStatsObject) GetTimerValues()`: Get a map of all timer values

### Store

Store is like Data Store, but less complex and provides less features.

#### Config

Options go into the `store.json` or `store.yaml` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of this store|
|`bucket_name`|String|Optional, specify the name of the bucket. Defaults to the name of the store.|
|`gobs`|Array|Optional, any extra types to register with gob|
|`extra_imports`|Array|Optional, any extra imports include for those types|

##### Example

```json
[
    {
        "name": "cache"
    }
]
```

#### Output

- `var <store name>Store = <store name>StoreObject{ Lock: &sync.Mutex{} }`: Global reference to your store
- `func storeSetup(storageDir string)`: Method to set up your stores
- `func storeTeardown()`: Method to tear down your stores
