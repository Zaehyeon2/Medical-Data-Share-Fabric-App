# -*- coding: utf-8 -*- 

from flask import Flask, jsonify, request
from flask_cors import CORS
import uuid
from pprint import pprint
import random
import datetime
import json
import requests
import ast
import os
import ipfshttpclient

# from pick import pick
# initialize Steem class

# configuration
DEBUG = True

# instantiate the app
app = Flask(__name__)
app.config.from_object(__name__)

# enable CORS
CORS(app)

def decrypt_file(key, in_file, out_filename, chunksize=24 * 1024):
    origsize = struct.unpack('<Q', in_file[:struct.calcsize('Q')])[0]
    iv = in_file[struct.calcsize('Q'):struct.calcsize('Q')+16]
    in_file = in_file[struct.calcsize('Q')+16:]
    decryptor = AES.new(key, AES.MODE_CBC, iv)
    with open(out_filename, 'wb') as outfile:
        while True:
            chunk = in_file[:chunksize]
            if len(chunk) == 0:
                break
            outfile.write(decryptor.decrypt(chunk))
            in_file = in_file[chunksize:]
        outfile.truncate(origsize)

def encrypt_file(key, in_file, chunksize=65536):
    iv = b'initialvector123'
    encryptor = AES.new(key, AES.MODE_CBC, iv)
    filesize = len(in_file)
    ret = (struct.pack('<Q', filesize))
    ret += iv
    while True:
        chunk = in_file[:chunksize]
        if len(chunk) == 0:
            break
        elif len(chunk) % 16 != 0:
            chunk += b' ' * (16 - len(chunk) % 16)
        ret += encryptor.encrypt(chunk)
        in_file = in_file[chunksize:]
    return ret

print("server On")

# Add response head 'Access-Control'
@app.after_request
def after_request(response):
	response.headers.add('Access-Control-Allow-Origin', '*')
	response.headers.add('Access-Control-Allow-Headers', 'Content-Type,Authorization')
	response.headers.add('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE,OPTIONS')
	return response

@app.route('/', methods=['GET', 'POST'])
def default():
	post_data = request.get_json()
	print(post_data)
	client = ipfshttpclient.connect()
	file = {'file': client.cat(post_data['hash'])}
	client.close()

	res = requests.post(post_data['dest_ip']+"/getData", files=file, data=post_data)
	print(res.text)
	return res.text

@app.route("/getData", methods=['GET', 'POST'])
def getData():
	f = request
	f.files['file'].save(f.form.to_dict()['hash'])
	return "Complete"

if __name__ == '__main__':
	app.run(host='0.0.0.0', port=int("9002"), debug=True)
