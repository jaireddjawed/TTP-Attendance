import os.path
import json
from datetime import date, datetime

from googleapiclient.discovery import build
from google.auth.transport.requests import Request
from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow

'''
    Checks if the authentication token file exists and is valid and not expired.
'''
def checkIfTokenIsValid():
    permissions = ['https://www.googleapis.com/auth/spreadsheets']

    creds = None

    if os.path.exists('token.json'):
        creds = Credentials.from_authorized_user_file('token.json', permissions)

    if not creds or not creds.valid:
      if creds and creds.expired and creds.refresh_token:
        creds.refresh(Request())
      else:
        flow = InstalledAppFlow.from_client_secrets_file('credentials.json', permissions)
        creds = flow.run_local_server(port=0)

      # Save the credentials for the next run
      with open('token.json', 'w') as token:
          token.write(creds.to_json())

    return creds

'''
    Creates a new sheet for the current month if it does not exist
'''
def createNewSheetForMonthIfNeeded(creds):
    with open('sheetInfo.json') as sheetInfoFile:
        sheetInfo = json.loads(sheetInfoFile.read())
        service = build('sheets', 'v4', credentials=creds)

        currentDate = datetime.now()
        currentMonth, currentYear = currentDate.strftime("%b"), currentDate.strftime("%Y")
        monthSheetName = f'{currentMonth} {currentYear}'

        # get all signin sheets for each month
        sheets = service.spreadsheets().get(spreadsheetId=sheetInfo['spreadsheetId']).execute()['sheets']

        # check if the current month sheet exists, if not create it
        currentMonthSheetExists = False
        print(sheets)
        for sheet in sheets:
            if sheet['properties']['title'] == monthSheetName:
                currentMonthSheetExists = True
                return

        if not currentMonthSheetExists:
            # create the spreadsheet for the new month
            service.spreadsheets().batchUpdate(
                spreadsheetId=sheetInfo['spreadsheetId'],
                body={
                    "requests": [
                        {
                            "addSheet": {
                                "properties": {
                                    "title": monthSheetName,
                                    "gridProperties": {
                                        "rowCount": 1000,
                                        "columnCount": 6
                                    }
                                }
                            }
                        }
                    ]
                }
            ).execute()

            # add the headers to the new sheet
            service.spreadsheets().values().update(
                spreadsheetId=sheetInfo['spreadsheetId'],
                range=f'{monthSheetName}!A1:E1',
                valueInputOption='USER_ENTERED',
                body={
                    'values': [[
                        'First Name',
                        'Last Name',
                        'Student ID',
                        'Sign In Date',
                        'Sign In Time'
                    ]]
                }
            ).execute()

'''
    Adds student's sign in information to the google sheet.
'''
def addStudentSignInToGoogleSheet(creds, signInInfo):
    try:
        createNewSheetForMonthIfNeeded(creds)

        currentDate = datetime.now()
        currentMonth, currentYear = currentDate.strftime("%b"), currentDate.strftime("%Y")
        sheetName = f'{currentMonth} {currentYear}'

        with open('sheetInfo.json', 'r') as sheetInfoFile:
            sheetInfo = json.loads(sheetInfoFile.read())
            service = build('sheets', 'v4', credentials=creds)

            signInDate = date.today().strftime("%m/%d/%Y")
            signInTime = datetime.today().strftime("%I:%M %p")

            service.spreadsheets().values().append(
                spreadsheetId=sheetInfo['spreadsheetId'],
                range=f'{sheetName}!A1:E1000',
                valueInputOption='USER_ENTERED',
                body={
                'values': [[
                    signInInfo['first-name'],
                    signInInfo['last-name'],
                    signInInfo['student-id'],
                    signInDate,
                    signInTime,
                    ''
                ]]
                }
            ).execute()
    except IOError:
       print("Error: Could not retrieve attendance sheet information. Please make sure that the sheetInfo.json file exists in the main directory.")
    except Exception as e:
        print("Error:", e)
