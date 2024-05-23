import unittest

from model import *

class TesTransaction(unittest.TestCase):
    def test_validation(self):
        @dataclass
        class TestCase:
            name: str
            card: str
            # date should be in format "2023-12-23"
            expiry_date: str
            # cvv should be three digit numbers
            cvv: str
            # amount should look like "12.89"
            amount: str
            currency: str
            merchant_id: str
            # expected
            hasException: bool
            expectedException: str = ''

        testCases = [
            TestCase(name='test 1', card="5555555555554443", expiry_date='2023-12-13', cvv='123', amount='12.89', currency='CNY', merchant_id='mr_123', hasException=False),
            TestCase(name='test 2', card="5555555555554443", expiry_date='2023-12-13', cvv='12', amount='12.89', currency='CNY', merchant_id='mr_123', hasException=True, expectedException='cvv should be three digits'),
            TestCase(name='test 3', card="5555555555554443", expiry_date='2023-12-13', cvv='ji8', amount='12.89', currency='CNY', merchant_id='mr_123', hasException=True, expectedException='cvv should be three digits'),
            TestCase(name='test 4', card="5555555555554443", expiry_date='2023/12/13', cvv='123', amount='12.89', currency='CNY', merchant_id='mr_123', hasException=True, expectedException='expiry date format shoud be YYYY-MM-DD'),
            TestCase(name='test 5', card="5555555555554443", expiry_date='2023-12-13', cvv='123', amount='er', currency='CNY', merchant_id='mr_123', hasException=True, expectedException='amount should be in format 12.89'),
            TestCase(name='test 6', card="5555555555554443", expiry_date='2023-12-13', cvv='123', amount='18.899', currency='CNY', merchant_id='mr_123', hasException=True, expectedException='amount decimal should be two')
        ]

        for ca in testCases:
            tran = Transaction(card=ca.card, expiry_date=ca.expiry_date, cvv=ca.cvv, amount=ca.amount, currency=ca.currency, merchant_id=ca.merchant_id)
            self.assertEqual(tran.card, ca.card, ca.name)
            if ca.hasException:
                with self.assertRaises(ValueError) as e:
                    tran.validate()
                self.assertEqual(str(e.exception), ca.expectedException, "failed test '{}' expected '{}', actual '{}'".format(
                    ca.name, ca.expectedException, str(e.exception)
                ),)

if __name__ == '__main__':
    unittest.main()
