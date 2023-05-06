from flask import Flask, render_template, request

app = Flask(__name__)

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
   app.run()
