import dashborg
import asyncio
import random

NUM_DATA = 30
COLORS = ["red", "green", "blue", "purple"]

def regen_data(req):
    rtn = []
    for i in range(NUM_DATA):
        point = {
            "X": random.randrange(50),
            "Y": random.randrange(50),
            "Val": random.randrange(50) + 1,
            "Color": COLORS[random.randrange(len(COLORS))],
        }
        rtn.append(point)
    req.set_data("$.data", rtn)

async def main():
    config = dashborg.Config(proc_name="d3", anon_acc=True, auto_keygen=True)
    client = await dashborg.connect_client(config)
    app = client.app_client().new_app("d3-test")
    app.set_app_title("D3 Demo")
    app.set_html(file_name="./panels/d3-test.html", watch=True)
    app.set_init_required(True)
    app.runtime.handler("regen-data", regen_data)
    app.runtime.init_handler(regen_data)
    await client.app_client().write_app(app, connect=True)
    await client.wait_for_shutdown()
    while True:
        await asyncio.sleep(1)

asyncio.run(main())

