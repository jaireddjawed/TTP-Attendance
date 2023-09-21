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

'''
def checkIfNewMonthNeeded(creds):
    sheetInfo = json.loads(open('sheetInfo.json', 'r').read())
    service = build('sheets', 'v4', credentials=creds)

    # get name of the current sheet
    currentSheetName = service.spreadsheets().get(spreadsheetId=sheetInfo['spreadsheetId']).execute()['properties']['title']
    print(currentSheetName)


'''
    Adds student's sign in information to the google sheet.
'''
def addStudentSignInToGoogleSheet(creds, signInInfo):
    sheetInfo = json.loads(open('sheetInfo.json', 'r').read())
    service = build('sheets', 'v4', credentials=creds)

    service.spreadsheets().values().append(
        spreadsheetId=sheetInfo['spreadsheetId'],
        range='Sheet1!A1:E100',
        valueInputOption='USER_ENTERED',
        body={
        'values': [[
            signInInfo['first-name'],
            signInInfo['last-name'],
            signInInfo['student-id'],
            date.today().strftime("%m/%d/%Y"),
            datetime.today().strftime("%I:%M %p")
        ]]
        }
    ).execute()
