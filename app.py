from flask import Flask, render_template, request
import webbrowser
import os
import gsheethelpers

app = Flask(__name__)
PORT = 5000

@app.route('/', methods=['GET'])
def renderSignInPage():
    return render_template('signin.html')

@app.route('/major', methods=['GET'])
def renderMajorPage():
    return render_template('major.html')

@app.route('/success', methods=['GET'])
def renderSuccessPage():
    return render_template('success.html')

@app.route('/submit-signin', methods=['POST'])
def submitSignIn():
    try:
        # check if google sheets authenication token is valid and not expired
        creds = gsheethelpers.checkIfTokenIsValid()

        # add the student sign in to the google sheet
        signInInfo = request.get_json()
        gsheethelpers.addStudentSignInToGoogleSheet(creds, signInInfo)

        # send a success status code and also indicate that data was created (201)
        return 'success', 201
    except Exception as e:
        # print an exception and send the error status code
        print(f'An error occurred: {e}.')
        return 'Internal Server Error', 500

if __name__ == '__main__':
    # credentials.json file is required for google sheets to be able to authenticate
    if not os.path.exists('credentials.json'):
        print('An error occrred: credentials.json file does not exist. Please download the credentials.json file from the Google API Console and place in this directory.')
        exit(0)

    # verify if token.json is not expired before running the application
    gsheethelpers.checkIfTokenIsValid()

    # open the browser to the sign in page
    webbrowser.open(f'http://127.0.0.1:{PORT}/')
    app.run(port=PORT)
