from pydantic import BaseModel
# from dateutil import parse
import datetime
from dataclasses import dataclass
from typing import Optional
from enum import Enum

OkResponse = {"data": "ok"}

@dataclass
class TransactionStatus(Enum):
    PENDING = 'pending'
    APPROVED = 'approved'
    DENIED = 'denied' 

class Transaction(BaseModel):
    card: str
    # date should be in format "2023-12-23"
    expiry_date: str
    # cvv should be three digit numbers
    cvv: str
    # amount should look like "12.89"
    amount: str
    currency: str
    merchant_id: str

    id: Optional[str] = None
    status: Optional[TransactionStatus] = None

    def validate(self) -> None:
        # validate expiry_date
        try:
            datetime.date.fromisoformat(self.expiry_date)
        except ValueError:
            raise ValueError("expiry date format shoud be YYYY-MM-DD")

        # cvv might have leading 0. It should be three digit numbers
        if len(self.cvv) != 3:
            raise ValueError("cvv should be three digits")
        try:
            int(self.cvv)
        except ValueError:
            raise ValueError("cvv should be three digits")

        # amount validate
        try:
            amount = float(self.amount)
        except ValueError:
            raise ValueError("amount should be in format 12.89")

        if amount != round(amount, 2):
            raise ValueError("amount decimal should be two")
