from flask import Flask, render_template, request
from datetime import datetime
import webbrowser

app = Flask(__name__)
PORT = 5000

@app.route('/', methods=['GET'])
def renderSignInPage():
    return render_template('signin.html')

@app.route('/success', methods=['GET'])
def renderSuccessPage():
    return render_template('success.html')

@app.route('/submit-signin', methods=['POST'])
def submitSignIn():
    # todo: submit sign in info to google sheet
    # for now, simply write the sign in info to a text file
    try:
        signInInfo = request.get_json()
        with open('signins.txt', 'a') as f:
            if signInInfo['student-id'] == '':
                signInInfo['student-id'] = 'N/A'

            f.write(f'{signInInfo["first-name"]}, {signInInfo["last-name"]}, {signInInfo["student-id"]}, {datetime.now().strftime("%m/%d/%Y %H:%M:%S %p")} \n')

        # send a success status code and also indicate that data was created (201)
        return 'success', 201
    except:
        # send a server error status code
        return 'Internal Server Error', 500

if __name__ == '__main__':
    # todo: check if google credentials exist or are expired
    webbrowser.open(f'http://127.0.0.1:{PORT}/')
    app.run(port=PORT)
