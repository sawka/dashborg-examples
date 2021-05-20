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

class PanelModel:
    def __init__(self, acc_model):
        self.acc_model = acc_model

    async def root_handler(self, req):
        if not req.check_auth(password="hello"):
            return
        await req.set_html_from_file("panels/accdemo.html")

    def refresh_accounts(self, req):
        req.set_data("$state.selaccid", None)
        req.invalidate_data("/accounts/.*")

    def get_acc_list(self, req):
        return self.acc_model.all_accs()

    def get_acc_by_id(self, req):
        return self.acc_model.acc_by_id(req.data)

    def upgrade(self, req):
        self.acc_model.upgrade(req.data)
        req.invalidate_data("/accounts/.*")

    def downgrade(self, req):
        self.acc_model.downgrade(req.data)
        req.invalidate_data("/accounts/.*")

    def remove(self, req):
        self.acc_model.remove_acc(req.data)
        self.refresh_accounts(req)

    def regen_acc_list(self, req):
        self.acc_model.regen_accounts()
        self.refresh_accounts(req)

    def create_account(self, req):
        form_data = req.panel_state.get("create", {})
        name = form_data.get("name", "").strip()
        email = form_data.get("email", "").strip()
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
        req.invalidate_data("/accounts/.*")
        

async def main():
    config = dashborg.Config(proc_name="accdemo", anon_acc=True, auto_keygen=True)
    await dashborg.start_proc_client(config)

    panel = PanelModel(AccModel())
    await dashborg.register_panel_handler("accdemo", "/", panel.root_handler)
    await dashborg.register_data_handler("accdemo", "/accounts/list", panel.get_acc_list)
    await dashborg.register_data_handler("accdemo", "/accounts/get", panel.get_acc_by_id)
    await dashborg.register_panel_handler("accdemo", "/acc/upgrade", panel.upgrade)
    await dashborg.register_panel_handler("accdemo", "/acc/downgrade", panel.downgrade)
    await dashborg.register_panel_handler("accdemo", "/acc/refresh-accounts", panel.refresh_accounts)
    await dashborg.register_panel_handler("accdemo", "/acc/regen-acclist", panel.regen_acc_list)
    await dashborg.register_panel_handler("accdemo", "/acc/remove", panel.remove)
    await dashborg.register_panel_handler("accdemo", "/acc/create-account", panel.create_account)
    while True:
        await asyncio.sleep(1)

asyncio.run(main())


