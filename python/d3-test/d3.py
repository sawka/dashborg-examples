import dashborg
import asyncio
import random

NUM_DATA = 30
COLORS = ["red", "green", "blue", "purple"]

async def root_handler(req):
    await req.set_html_from_file("panels/d3-test.html")
    regen_data(req)

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
    await dashborg.start_proc_client(config)
    await dashborg.register_panel_handler("d3-test", "/", root_handler)
    await dashborg.register_panel_handler("d3-test", "/regen-data", regen_data)
    while True:
        await asyncio.sleep(1)

asyncio.run(main())

