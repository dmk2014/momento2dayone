[![Go Report Card](https://goreportcard.com/badge/github.com/dmk2014/momento2dayone)](https://goreportcard.com/report/github.com/dmk2014/momento2dayone)

# Momento - Day One Import Tool

## Introduction
This is a tool to import a Momento export into Day One. Please backup all Momento/Day One data before attempting to use this tool. It is not guaranteed to handle all edge cases. See the license for more information.

This tool is only compatible with systems running macOS.

A log file is output to the directory in which the tool is invoked. It will detail any parse errors. Import does not fail when encountering an error but will output details of the entry along with the error output from the Day One CLI.

## Usage
Build/install using standard Golang tooling. There are two flags available. The path is required. Moment count can be provided to verify the parse step when the number of expected Moments is known.

TODO: Table of flags.

## Install Day One Command Line Tools
The DayOne CLI tools are required for the import step. The application will not run if it cannot verify that the tools have been installed.

```
sudo /Applications/Day\ One.app/Contents/Resources/install_cli.sh
```

## General

Momento Tags, People and Places all become Tags in Day One. Geolocation information on places is not maintained during import.

## Verify Backup

Using a Momento backup file (SQLite), the data can be verified. There cannot be any dates that match the below regular expression on a single line.

```sql
SELECT ZNOTES, ZDATE
FROM ZMOMENT
WHERE ZNOTES REGEXP '[0-9]{1,2}\s[a-zA-Z]{3,9}\s[0-9]{4}'
```

There cannot be any times that match the below regular expression on a single line.

```sql
SELECT ZNOTES, ZDATE
FROM ZMOMENT
WHERE ZNOTES REGEXP '[0-9]{2}:[0-9]{2}'
```

## Parser
The parser is fast - it completes a ~90k line text file with about 6200 entries in less than 100ms. It returns a slice of Moments which satisfy the interface expected by the import package.

## Importer
The import will take some time. It must process each entry sequentially. There was an issue encountered during testing where photos were not correctly imported. This was fixed by having the import function wait for 10 seconds every 100 entries.

Import progress is displayed on the command line.