from fastapi import FastAPI, HTTPException
from neo4j import GraphDatabase, basic_auth
from xmlrpc.client import ServerProxy
from pydantic import BaseModel, SecretStr
import hashlib
import socket
import json

app = FastAPI()

class User(BaseModel):
    username: str
    password: SecretStr

# Connect to supervisord
server = ServerProxy('http://' + config_vars["host"] + ':9001/RPC2')

config_vars = json.load(open('config.json')

# Neo4j connection
driver = GraphDatabase.driver("bolt://" + config_vars["host"] + ":7687", auth=basic_auth(config_vars["neouname"], config_vars["neopword"]))

def md5_hash(string):
    return hashlib.md5(string.encode()).hexdigest()

def get_user(username, password):
    password = md5_hash(password)
    with driver.session() as session:
        result = session.run("MATCH (a:account) WHERE u.username = $username AND u.password = $password AND u.permissions=31 RETURN u", username=username, password=password)
        return result.single() is not None

def stop_process(process_name):
    try:
        server.supervisor.stopProcess(process_name)
        return True
    except socket.error as e:
        print(f'Could not connect to supervisord: {e}')
        return False
    except Exception as e:
        print(f'Could not stop process: {e}')
        return False

def start_process(process_name):
    try:
        server.supervisor.startProcess(process_name)
        return True
    except socket.error as e:
        print(f'Could not connect to supervisord: {e}')
        return False
    except Exception as e:
        print(f'Could not start process: {e}')
        return False

def restart_process(process_name):
    if not stop_process(process_name):
        return False
    return start_process(process_name)

@app.get("/stop_nexus/")
async def stop_nexus(user: User)):
    if not get_user(user.username, user.password.get_secret_value()):
        raise HTTPException(status_code=400, detail="Invalid credentials")

    if not stop_process('nexus'):
        raise HTTPException(status_code=500, detail="Could not stop process")
    return {"message": "Process stopped successfully"}

@app.get("/start_nexus/")
async def start_nexus(user: User)):
    if not get_user(user.username, user.password.get_secret_value()):
        raise HTTPException(status_code=400, detail="Invalid credentials")

    if not start_process('nexus'):
        raise HTTPException(status_code=500, detail="Could not start process")
    return {"message": "Process started successfully"}

@app.get("/restart_nexus/")
async def restart_nexus(user: User)):
    if not get_user(username, password):
        raise HTTPException(status_code=400, detail="Invalid credentials")

    if not restart_process('nexus'):
        raise HTTPException(status_code=500, detail="Could not restart process")
    return {"message": "Process restarted successfully"}

@app.get("/clean_shutdown")
def call_go_endpoint(user: User)):
    if not get_user(username, password):
        raise HTTPException(status_code=400, detail="Invalid credentials")
    url = "http://127.0.0.1:1234/clean_shutdown"
    headers = {"Content-Type": "application/json"}
    data = {"token": config_vars["resttoken"]}
    response = requests.post(url, headers=headers, data=json.dumps(data))

    if response.status_code == 200:
        return {"message": "Valid token received"}
    else:
        return {"message": "Invalid token"}