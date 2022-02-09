from googleapiclient.discovery import build
from generateToken import generateToken
from os import path, system
from json import load
from datetime import date, datetime
from time import sleep

def clearTerminal():
  system('cls' if path.exists('cls') else 'clear')

def addStudentManually(spreadsheetId, service):
  # todo: validate the user's input
  firstName = input('Please enter your first name: ')
  lastName = input('Please enter your last name: ')
  studentId = input('Please enter your SID: ')

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
  sleep(0.5)
  clearTerminal()

def parseCardSwipeInput(card):
  fullName = card.split('^')[1].split('/')

  # first and last name are in reverse order
  firstName = fullName[1]
  lastName = fullName[0]

  return firstName, lastName

def addStudentCardSwipe(spreadsheetId, service, firstName, lastName):
  service.spreadsheets().values().append(
    spreadsheetId=spreadsheetId,
    range='Sheet1!A1:E100',
    valueInputOption='USER_ENTERED',
    body={
      'values': [[
        firstName,
        lastName,
        '', # no student id is provided from the card swipe
        date.today().strftime("%m/%d/%Y"),
        datetime.today().strftime("%I:%M %p")
      ]]
    }
  ).execute()

  clearTerminal()
  sleep(0.5)
  print('Successfully signed in.')
  sleep(0.5)
  clearTerminal()

def main():
  creds = generateToken()
  service = build('sheets', 'v4', credentials=creds)

  sheet = open(path.join(path.dirname(__file__), 'sheetInfo.json'), 'r')
  sheet = load(sheet)

  while True:
    print('Welcome to the TTP Attendance Tracker.')
    print('By signing-in you acknowledge that you have read and agree to abide by the Engineering Transfer Center (WCH 103) Space Policies.\n')

    choice = input('Please swipe your card or press "m" and "Enter" to add details manually: ')

    if choice == 'm':
      addStudentManually(sheet['spreadsheetId'], service)
    elif choice == 'q':
      break
    else:
      # get the first name and last name of student from card swipe
      firstName, lastName = parseCardSwipeInput(choice)
      addStudentCardSwipe(sheet['spreadsheetId'], service, firstName, lastName)

if __name__ == '__main__':
  main()
