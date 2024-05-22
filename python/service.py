from repository import Repository
from model import *

class Service:
    repo: Repository

    def __init__(self, repo: Repository) -> None:
        self.repo = repo

    def create(self, body: Transaction):
        # create transaction firstly
        try:
            tran_id: int = self.repo.createTransaction(body)
        except Exception as e:
            raise Exception("error creating transaction", e)
    
        # if transaction should be approved
        approved: bool = self.isAcquirerApproved(body.card)

        # update transaction status
        try:
            status = TransactionStatus.APPROVED
            if not approved:
                status = TransactionStatus.DENIED
            self.repo.updateTransactionStatus(tran_id, status)
        except Exception as e:
            raise Exception("error updating transaction status")

    # should this transaction be approved
    def isAcquirerApproved(self, card: str) -> bool:
        last_digit = card[-1]
        return int(last_digit) % 2 == 0

