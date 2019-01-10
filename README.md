# Bayar ![](https://img.shields.io/badge/language-golang-lightblue.svg)
Bayar is a personal API service to log expenses into a Google Sheets spreadsheet. It is written with the goal of simplifying the process of logging expenses, and also as an opportunity to learn [Go](https://golang.org).

It is used in conjunction with the [Apple Shortcuts](https://support.apple.com/en-my/guide/shortcuts/welcome/ios) app to create a convenient interface for logging expenses into a custom Google Sheets spreadsheet.

## Usage

**Step 1:** Create a Bayar application directory somewhere in your system
```bash
$ mkdir /path/to/dir
```

**Step 2:** Create a Bayar configuration file.  See [Configuration](#Configuration).
```bash
$ vim bayar.json
```

**Step 3:** Set the `BAYAR_CONFIG` environment variable to point to the configuration file.
```bash
$ export BAYAR_CONFIG=/path/to/config
```

**Step 4:**  Download a [Google API OAuth 2.0](https://developers.google.com/api-client-library/python/auth/installed-app#creatingcred) configuration file into your Bayar application directory

**Step 5:** Compile and run the Bayar binary file. See [Building](#Building).
```bash
$ ./bayar
```

**Step 6:** Visit `http://localhost:8888/startAuthorization` to authorize the Bayar application with Google

**Step 7:** You're done! You can start making API calls at http://localhost:8888/newExpense


### Configuration

Bayar refers to a JSON configuration file specified by the `BAYAR_CONFIG` environment variable.
```js
{
  "ApplicationDirectory": "", // (Required) Bayar will store application data here
  "SpreadsheetID": "", // (Required) The Google Sheets spreadsheet ID
  "SheetName": "", // (Required) The name of the sheet to log expenses
  "HostAddress": "localhost", // (Optional) Defaults to localhost
  "PortNumber": 8888, // (Optional) Defaults to 8888
  "GoogleConfigurationFileName": "client_secret.json", // (Optional) The name of the Google API service file, defaults to client_secret.json
}
```

## Building

Run any one of the scripts in the `scripts/` directory from the project root directory. The binary files will be built in the `dist\` directory.

```cmd
> scripts\build-windows_amd64.bat
```

## API Reference

#### POST `/newExpense`
Log a new expense into the spreadsheet.  
* Request Body: `JSON`
```js
{
  "label": "McDonalds", // optional
  "category": "Food",
  "cost": 9.95
}
```