from fastapi import FastAPI, HTTPException
from service import Service
from repository import Repository
from model import *
import typing

app = FastAPI()

# init service and repository
svc = Service(Repository())

@app.get("/health")
async def root():
    return OkResponse

@app.post("/payment")
async def createPayment(body: Transaction):
    # validate request body
    try:
        body.validate()
    except ValueError as e:
        raise HTTPException(status_code = 400, detail=str(e))

    # create this new transaction
    try:
        svc.create(body)
    except Exception as e:
        raise HTTPException(status_code = 500, detail=str(e))

    return OkResponse
