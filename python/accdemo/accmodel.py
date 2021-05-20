import uuid
import random

FOOD_STRS = ["Apple", "Banana", "Guava", "Orange", "Blackberry", "Mango", "Kiwi", "Raspberry", "Pineapple", "Avacado", "Onion", "Lettuce", "Cheese", "Almond", "Cake", "Walnut"]
SUFFIX_STRS = ["LLC", "Inc", "Corp", "Corp", "Ltd", "", "", "", "", "", "", "", ""]
MODWORDS_STRS = ["Star", "Lightning", "Flash", "Media", "Data", "Micro", "Net", "Echo", "World", "Red", "Blue", "Green", "Yellow", "Purple", "Tele", "Cloud", "Insta", "Face", "Super"]
EMAIL_STRS = ["mike", "matt", "michelle", "pat", "jim", "marc", "andrew", "alan", "henry", "jenny"]

class AccModel:
    def __init__(self):
        self.regen_accounts()

    def acc_by_id(self, accid):
        for acc in self.accs:
            if acc.accid == accid:
                return acc
        return None

    def all_accs(self):
        return self.accs

    def create_acc(self, name, email):
        acc = AccType(name=name, email=email)
        self.accs.append(acc)
        return acc.accid

    def remove_acc(self, accid):
        self.accs[:] = [acc for acc in self.accs if acc.accid != accid]

    def upgrade(self, accid):
        acc = self.acc_by_id(accid)
        if acc is not None:
            acc.is_paid = True

    def downgrade(self, accid):
        acc = self.acc_by_id(accid)
        if acc is not None:
            acc.is_paid = False

    def regen_accounts(self):
        self.accs = []
        for i in range(5):
            self.accs.append(make_random_acc())
    

class AccType:
    def __init__(self, name, email, is_paid=False):
        self.accid = str(uuid.uuid4())
        if name is None:
            raise ValueError("name cannot be None")
        if email is None:
            raise ValueError("email cannot be None")
        self.name = name
        self.email = email
        self.is_paid = is_paid

    def to_dict(self):
        return {"AccId": self.accid, "AccName": self.name, "Email": self.email, "IsPaid": self.is_paid}
    

def random_word(word_list):
    return word_list[random.randrange(len(word_list))]

def make_random_acc():
    name = (random_word(MODWORDS_STRS) + " " + random_word(FOOD_STRS) + " " + random_word(SUFFIX_STRS)).strip()
    email = random_word(EMAIL_STRS) + str(random.randrange(70) + 10) + "@nomail.com"
    return AccType(name=name, email=email)
