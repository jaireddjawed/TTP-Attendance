from google.oauth2.credentials import Credentials
from google.auth.transport.requests import Request
from google_auth_oauthlib.flow import InstalledAppFlow
from os import path
from json import load

def generateToken():
  sheet = open(path.join(path.dirname(__file__), 'sheetInfo.json'), 'r')
  sheet = load(sheet)

  creds = None

  if path.exists('token.json'):
    creds = Credentials.from_authorized_user_file(
      path.join(path.dirname(__file__), 'token.json'),
      sheet['scopes']
    )

  if creds and creds.expired and creds.refresh_token:
    creds.refresh(Request())

  elif not creds or not creds.valid:
    flow = InstalledAppFlow.from_client_secrets_file(
      path.join(path.dirname(__file__), 'TTP_CLIENT.json'),
      sheet['scopes']
    )
    creds = flow.run_local_server(port=0)

  with open(path.join(path.dirname(__file__), 'token.json'), 'w') as token:
    token.write(creds.to_json())

  return creds
