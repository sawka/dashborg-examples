import dashborg
import asyncio
from accmodel import AccModel
from accmodel import AccType
import re

EMAIL_REGEX = r"^[a-zA-Z0-9.-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

def my_serialize(obj):
    if isinstance(obj, AccType):
        rtn = dashborg.serialize(obj)
        rtn["foo"] = True
        return rtn
    return dashborg.serialize(obj)

class AppModel:
    def __init__(self, acc_model):
        self.acc_model = acc_model

    async def root_handler(self, req):
        if not req.check_auth(password="hello"):
            return
        await req.set_html_from_file("panels/accdemo.html")

    def refresh_accounts(self, req):
        req.set_data("$state.selaccid", None)
        req.invalidate_data(".*")

    def get_acc_list(self):
        return self.acc_model.all_accs()

    def get_acc_by_id(self, _id : str = None):
        return self.acc_model.acc_by_id(_id)

    def upgrade(self, req, _id : str):
        self.acc_model.upgrade(_id)
        req.invalidate_data(".*")

    def downgrade(self, req, _id : str):
        self.acc_model.downgrade(_id)
        req.invalidate_data(".*")

    def remove(self, req, _id : str):
        self.acc_model.remove_acc(_id)
        req.invalidate_data(".*")

    def regen_acc_list(self, req):
        self.acc_model.regen_accounts()
        req.invalidate_data(".*")

    def create_account(self, req, app_state):
        print(f"** create account {app_state}")
        name = app_state.get("create", {}).get("name", "").strip()
        email = app_state.get("create", {}).get("email", "").strip()
        errors = {}
        if name == "":
            errors["name"] = "Name must not be empty"
        elif len(name) > 40:
            errors["name"] = "Name can only be 40 characters"
        if email == "":
            errors["email"] = "Email must not be empty"
        elif len(email) > 40:
            errors["email"] = "Email can only be 40 characters"
        elif not re.search(EMAIL_REGEX, email):
            errors["email"] = "Email format is not correct"
        if len(errors) > 0:
            req.set_data("$state.create.errors", errors)
            return
        req.set_data("$state.create.errors", None)
        new_accid = self.acc_model.create_acc(name, email)
        req.set_data("$state.createAccountModal", False)
        req.set_data("$state.selaccid", new_accid)
        req.invalidate_data(".*")
        

async def main():
    config = dashborg.Config(proc_name="accdemo", anon_acc=True, auto_keygen=True)
    client = await dashborg.connect_client(config)

    m = AppModel(AccModel())
    app = client.app_client().new_app("accdemo")
    app.set_app_title("Account Demo")
    app.set_html(file_name="panels/accdemo.html", watch=True)
    app.runtime.handler("get-accounts-list", m.get_acc_list, pure_handler=True)
    app.runtime.handler("get-account", m.get_acc_by_id, pure_handler=True)
    app.runtime.handler("acc-upgrade", m.upgrade)
    app.runtime.handler("acc-downgrade", m.downgrade)
    app.runtime.handler("refresh-accounts", m.refresh_accounts)
    app.runtime.handler("regen-acclist", m.regen_acc_list)
    app.runtime.handler("acc-remove", m.remove)
    app.runtime.handler("create-account", m.create_account)
    await client.app_client().write_app(app, connect=True)
    await client.wait_for_shutdown()
    while True:
        await asyncio.sleep(1)

asyncio.run(main())


