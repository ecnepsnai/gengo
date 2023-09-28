# cbgen

CBgen is yet another Golang code generator, because there aren't enough of those already.

It can generate the following things: github.com/ecnepsnai/ds implementations, directory utilities, enums,
gob registrations, a state store, a statistics tracker, github.com/ecnepsnai/store implementations, and version tracking.

What's CB stand for? *I don't remember.*

# Usage

```
Usage: cbgen [Options]
-n --name <name>        Package name, defaults to 'main'
-c --config-dir <dir>   Config dir, defaults to current dir
-o --output-dir <dir>   Output dir, defaults to current dir
```

Ensure that your $GOBIN is in your $PATH and run `cbgen` in the directory where your configuration JSON files are.

## Data Store

Data Stores provide an implementation of the github.com/ecnepsnai/ds package.

### Config

Options go into the `data_store.json` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of the data store iteself, as seen within code|
|`object`|String|The golang object to register to the store|
|`unordered`|Bool|Optional, if the table should not be ordered. Defaults to false.|

**Example:**

```json
[
    {
        "name": "User",
        "object": "User",
        "unordered": true
    }
]
```

### Output

- `func dataStoreSetup()`: Method to set up all data stores
- `func dataStoreTeardown()`: Method to tear down all data stores
- `type <name_lowercase>StoreObject struct{ table *ds.Table }`: Type associated with the store
- `var <name_titlecase>Store: <name_lowercase>StoreObject`: Global variable associated with the store, set up when `dataStoreSetup()` is called

## Directories

Directories provides utilities for data directory management.

### Config

Options go into the `directory.json` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of the directory as seen within code|
|`is_data`|Bool|If this data is a data directory|
|`dir_name`|String|The name of the directory itself as seen on the filesystem|
|`required`|Bool|Optional, if this directory must already exist. If false, the directory will be created for you|
|`subdirs`|Array|Optional subdirectories under this parent|

**Example:**

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

### Output

- `var operatingDirectory: String`: Global variable of the operating directory of the application
- `var dataDirectory: String`: Global variable of the data directory of the application
- `var Directories: apiDirectories`: Global variable of absolute paths to your directories
- `func fsSetup()`: Method to set up directories
- `func DirectoryExists(directoryPath string) bool`: Method to check if a directory exists
- `func MakeDirectoryIfNotExist(directoryPath string)`: Method to create a directory if it does not already exist
- `func FileExists(filePath string) bool`: Method to check if a file exists

## Enums

Enums are enums, the one thing that go really needs.

### Config

Options go into the `enum.json` file.

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

**Example:**

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

### Output

- `const ( <enum name><value key> = <value> ...)`: A const of all possible enum values for this enum.
- `var All<enum name> = []<enum type>`: An array of all enum values.
- `var <enum name>Map = map[<enum type>]<enum type>`: A map of enum keys to enum values.
- `func Is<enum name>(q <enum type>) bool`: A method to validate that the given value is a valid enum value.
- `func ForEach<enum name>(m func(value <enum type>))`: A convience method to iterate over each enum value

## Gobs

Gobs provide an easy way to register many objects with gob in one place, and to prevent a panic if they're already registered.

### Config

Options go into the `gob.json` file.

|Property|Type|Description|
|-|-|-|
|`type`|String|The golang type to register, including the `{}`|
|`import`|String|The import path for this type. Imports are deduplicated|

**Example:**

```json
[
    {
        "type": "time.Time{}",
        "import": "time"
    }
]
```

### Output

- `func gobSetup()`: Method to register all of your gob types

## State Store

State Store is a simple store interface for persisting configuration or other properties.

### Config

Options go into the `state.json` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of this state property|
|`type`|String|The golang type for this property|
|`default`|String|The default value of this property as it would be seen in Golang code. I.E. if it's a string include quotation marks|
|`import`|String|The import path for this type. Imports are deduplicated|

**Example:**

```json
[
    {
        "name": "TableVersion",
        "type": "int",
        "default": "0"
    }
]
```

### Output

- `var State *cbgenStateObject`: The global state object
- `func stateSetup()`: Method to set up the state store
- `func (s *cbgenStateObject) Close()`: Method to close the state store
- `func (s *cbgenStateObject) Get<property name>() <property type>`: Method to get the property value, or the default value. Threadsafe.
- `func (s *cbgenStateObject) Set<property name>(value <property type>)`: Method to set the property value. Threadsafe.

## Stats

Stats provides an interface to the github.com/ecnepsnai/stats package.

### Config

Options go into the `stats.json` file.

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

**Example:**

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

### Output

- `type cbgenStatsCounters struct { <counter name> *stats.Counter ... }`: All of your counter
- `type cbgenStatsTimedCounters struct { <counter name> *stats.TimedCounter ... }`: All of your timed counters
- `type cbgenStatsTimers struct { <timer name> *stats.Timer ... }`: All of your timers
- `var Stats *cbgenStatsObject`: The global stats object
- `func statsSetup()`: Method to set up the stats interface
- `func (s *cbgenStatsObject) Reset()`: Reset all counters to their default values
- `func (s *cbgenStatsObject) GetCounterValues()`: Get a map of all counter values
- `func (s *cbgenStatsObject) GetTimedCounterValues()`: Get a map of all timed counter values
- `func (s *cbgenStatsObject) GetTimedCounterValuesFrom(d time.Duration)`: Get a map of all timed counter values from d
- `func (s *cbgenStatsObject) GetTimerAverages()`: Get a map of all timer averages
- `func (s *cbgenStatsObject) GetTimerValues()`: Get a map of all timer values

## Store

Store is like Data Store, but less complex and provides less features.

### Config

Options go into the `store.json` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of this store|
|`bucket_name`|String|Optional, specify the name of the bucket. Defaults to the name of the store.|
|`gobs`|Array|Optional, any extra types to register with gob|
|`extra_imports`|Array|Optional, any extra imports include for those types|

**Example:**

```json
[
    {
        "name": "cache"
    }
]
```

### Output

- `var <store name>Store = <store name>StoreObject{ Lock: &sync.Mutex{} }`: Global reference to your store
- `func storeSetup()`: Method to set up your stores
- `func storeTeardown()`: Method to tear down your stores

## Version

If you call cbgen with the `-v <version>` argument, it will a version file with a global variable to the version.
