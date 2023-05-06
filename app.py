from flask import Flask, render_template, request
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
    try:
        signInInfo = request.get_json()
        print(signInInfo)

        return 'success', 201
    except:
        return 'Internal Server Error', 500

if __name__ == '__main__':
    webbrowser.open(f'http://127.0.0.1:{PORT}/')
    app.run(port=PORT)
