from collections import namedtuple

from flask import Flask, render_template, redirect, url_for, request

from src.close_animal import get_close_animals

app = Flask(__name__)

Message = namedtuple('Message', 'dog lev_dog cat lev_cat')
messages = [Message('', 0, '', 0)]
person = '[your word]'


@app.route('/', methods=['GET'])
def hello_world():
    return render_template('index.html')


@app.route('/main', methods=['GET'])
def main():
    return render_template('main.html', messages=messages, name=person)


@app.route('/add_message', methods=['POST'])
def add_message():
    global person
    text = request.form['text']
    person = text[::]
    messages.pop()
    close_animals = get_close_animals(text)
    print(close_animals)
    messages.append(Message(close_animals[0], close_animals[1],
                            close_animals[2], close_animals[3]))
    return redirect(url_for('main'))


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=9000)
