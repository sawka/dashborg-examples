import dashborg
import asyncio
import sqlite3

class SimpleQueryModel:
    def __init__(self):
        self.dbcon = sqlite3.connect("data/yc-companies.db", isolation_level=None)
        pass
    
    async def get_rows(self, req):
        cursor = self.dbcon.cursor()
        cursor.execute("SELECT * FROM companies")
        colnames = [v[0] for v in cursor.description]
        rtn = []
        for row in cursor:
            company = dict(zip(colnames, row))
            rtn.append(company)
        return rtn

async def main():
    config = dashborg.Config(proc_name="simplequery", anon_acc=True, auto_keygen=True)
    client = await dashborg.connect_client(config)
    m = SimpleQueryModel()
    app = client.app_client().new_app("simplequery")
    app.set_app_title("Simple Query")
    app.set_html(file_name="./panels/simplequery.html", watch=True)
    app.runtime.handler("get-rows", m.get_rows, pure_handler=True)
    await client.app_client().write_app(app, connect=True)
    await client.wait_for_shutdown()
    while True:
        await asyncio.sleep(1)

asyncio.run(main())
