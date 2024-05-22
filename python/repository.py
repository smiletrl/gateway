import time
import typing
from model import *

class Repository:
    # a temporary mock db storage engine
    mockDB:  typing.Dict[int, Transaction] = {}

    # create a new transaction
    def createTransaction(self, body: Transaction) -> int:
        tran_id = self.generateTransactionID()

        body.status = TransactionStatus.PENDING
        
        self.mockDB[tran_id] = body
        return tran_id

    # update transaction status
    def updateTransactionStatus(self, tran_id: int, status: TransactionStatus):
        # assign a new status for this specific transaction
        self.mockDB[tran_id].status = status

    # generate a unique transaction id
    def generateTransactionID(self) -> int:
        transactionID = time.time_ns()
        return transactionID
