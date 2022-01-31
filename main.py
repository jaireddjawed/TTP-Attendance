from googleapiclient.discovery import build
from generateToken import generateToken
from os import path, system
from json import load
from datetime import date, datetime
from time import sleep

def clearTerminal():
  system('cls' if path.exists('cls') else 'clear')

def addStudentManually(spreadsheetId, service):
  firstName = input('Please enter your first name: ')
  lastName = input('Please enter your last name: ')
  studentId = input('Please enter your student ID: ')

  service.spreadsheets().values().append(
    spreadsheetId=spreadsheetId,
    range='Sheet1!A1:E100',
    valueInputOption='USER_ENTERED',
    body={
      'values': [[
        firstName,
        lastName,
        studentId,
        date.today().strftime("%m/%d/%Y"),
        datetime.today().strftime("%I:%M %p")
      ]]
    }
  ).execute()

  clearTerminal()
  sleep(0.5)
  print('Successfully signed in.')
  sleep(2)
  clearTerminal()

def main():
  creds = generateToken()
  service = build('sheets', 'v4', credentials=creds)

  sheet = open(path.join(path.dirname(__file__), 'sheetInfo.json'), 'r')
  sheet = load(sheet)

  while True:
    print('Welcome to the TTP Attendance Tracker.')
    print('Please swipe your card or press "Enter" to add details manually.')
    print("By signing-in you acknowledge that you have read and agree to abide by the Engineering Transfer Center (WCH 103) Space Policies.\n")

    addStudentManually(sheet['spreadsheetId'], service)

if __name__ == '__main__':
  main()
