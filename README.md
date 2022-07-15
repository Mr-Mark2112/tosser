 Tosser
## What is tosser?
A programm for synchronizing files between two directories (source and destination).

---
 Help
    ```
    tosser help
    ``` 
---

## Installation

* Via a GO install
  ```shell
  go get -u github.com/Mr-Mark2112/tosser
  ```
---

## Building From Source

 In order to build tosser from source you must:

 1. Clone the repo
 2. Build and run the executable

      ```shell
      make build && ./execs/tosser
      ```
---
## How to use it?
* To initialize the default config, you need to run the program once, then you can specify the necessary directories and other parameters;
* By default, a config for Tosser is generated in /etc/tosser/config.yaml; also you can create a config by yourself, name it "config.yaml", configure the necessarily parameters (shown below) or you can use the config.yaml file from this repo. You can specify the path to a config.yaml file by flag `-c`.
* To launch TOsser use command 
```
tosser run 
```
* logs are writting in 
```
/etc/tosser/logs
```

## Configuration
src_dir: /tmp
dst_dir: /tmp_test_tosser
max_copy_threads: 4
rescaninterval: 3600 #in seconds
loglevel: INFO

| Parameter               | required               | Comment                                     |
|-------------------------|------------------------|---------------------------------------------|
| src_dir                 | `+`                    | A source directory to be synced             |
| dst_dir                 | `+`                    | A destination directory for syncing         |
| max_copy_threads        | `+`                    | The more threads the faster will be syncing |
| rescaninterval          | `+`                    | Interval in seconds when syncing strarts    |
| loglevel                | `+`                    | Level of logging (DEBUG INFO ERROR)         |


