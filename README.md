# TTP-Attendance

## Description
Desktop application that records student attendance at the TTP Center at UC Riverside (Winston Chung Hall Room 103).

## Setup
* Install the following dependencies:
```
pip3 install --upgrade google-api-python-client google-auth-httplib2 google-auth-oauthlib
```

* Setup a <a href="https://console.cloud.google.com/">Google Cloud Platform Project</a> and enable the Google Sheets API.

* Setup a desktop application to access the Google Cloud Platform and download your ```credentials.json``` file.

* Create a Google Sheet and add it's ID to the ```sheetInfo.json``` file.

* Run the application
```
python3 main.py
```
