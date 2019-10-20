# cbgen

CBgen is yet Golang code generator, because there aren't enough of those already.

It can generate the following things: github.com/ecnepsnai/ds implementations, directory utilities, enums,
gob registrations, a state store, a statistics tracker, a basic store, and version tracking.

# Usage

```
Usage: cbgen -n <package name> [-v <package version]
-n --name     Package name
-v --version  Package version. Including will generate a version go file
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

- `func DataStoreSetup()`: Method to set up all data stores
- `func DataStoreTeardown()`: Method to tear down all data stores
- `type <name_lowercase>StoreObject struct{ table *ds.Table }`: Type associated with the store
- `var <name_titlecase>Store: <name_lowercase>StoreObject`: Global variable associated with the store, set up when `DataStoreSetup()` is called

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
|`values`|Array|The possible values for this enum|
|`values.key`|String|The key for this value|
|`values.description`|String|The description for this value|
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

### Output

- `const ( <enum name><value key> = <value> ...)`: A const of all possible enum values for this enum
- `var All<enum name> = []<enum type>`: An array of all enum values
- `var <enum name>Map = map[<enum type>]<enum type>`: A map of enum keys to enum values
- `var <enum name>NameMap = map[string]<enum type>`: A map of enum names to enum values
- `func Is<enum name>(q <enum type>) bool`: A method to validate that the given value is a valid enum value
- `var <enum name>Schema = []map[string]interface{}`: A definition of all enums, basically just the contents of the JSON config file
- `var AllEnums = map[string]interface{}`: All enums mapping to their name map

## Gobs

Gobs provide an easy way to register many objects with gob in one place, and to prevent crashing if they're already registered.

### Config

Options go into the `gob.json` file.

|Property|Type|Description|
|-|-|-|
|`type`|String|The golang type to register, including the `{}`|
|`include`|String|The include path for this type. Includes are deduplicated|

**Example:**

```json
[
    {
        "type": "time.Time{}",
        "include": "time"
    }
]
```

### Output

- `func GobSetup()`: Method to register all of your gob types

## State Store

State Store is a simple store interface for persisting configuration or other properties.

### Config

Options go into the `state.json` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of this state property|
|`type`|String|The golang type for this property|
|`default`|String|The default value of this property as it would be seen in Golang code. I.E. if it's a string include quotation marks|

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

- `var State *stateObject`: The global state object
- `func StateSetup()`: Method to set up the state store
- `func (s *stateObject) Close()`: Method to close the state store
- `func (s *stateObject) Get<property name>() <property type>`: Method to get the property value, or defualt
- `func (s *stateObject) Set<property name>(value <property type>)`: Method to set the property value

## Stats

Stats provides a statistics tracking interface for your application. It is broken up into two primary uses: counters and timers. Counters provide a simeple incrementer/decrementer to track a value. They always start at 0. Timers provide a way to track the time of an operation, and to get the average time of recent operations.

### Config

Options go into the `stats.json` file.

|Property|Type|Description|
|-|-|-|
|`counters`|Array|The counters|
|`timers`|Array|The timers|
|`counters.name`|String|The name of this counter|
|`counters.description`|String|The description of this counter|
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
    "timers": [
        {
            "name": "BuildDuration",
            "description": "Duration of builds"
        }
    ]
}
```

### Output

- `type Counters struct { <counter name> uint64 ... }`: All of your counters and their values
- `type Timers struct { <timer name> *ring.Ring ... }`: All of your timer samples
- `var Stats *statsObject`: The global stats object
- `func StatsSetup()`: Method to set up the stats interface
- `func (s *statsObject) Reset()`: Reset all counters to their default values
- `func (s *statsObject) Increment<counter name>()`: Increment the counter by 1
- `func (s *statsObject) Decrement<counter name>()`: Decrement the counter by 1
- `func (s *statsObject) Set<counter name>(newVal uint64)`: Set the value of the counter
- `func (s *statsObject) Add<timer name>(value float32)`: Add a sample to the timer
- `func (s *statsObject) GetTimerAverages() map[string]float32`: Get the averages for all timers

## Store

Store is like Data Store, but less complex and provides less features.

### Config

Options go into the `store.json` file.

|Property|Type|Description|
|-|-|-|
|`name`|String|The name of this store|
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
- `func StoreSetup()`: Method to set up your stores
- `func StoreTeardown()`: Method to tear down your stores

## Version

If you call cbgen with the `-v <version>` argument, it will a version file with a global variable to the version.