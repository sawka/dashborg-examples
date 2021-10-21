import dashborg
import asyncio

class TodoModel:
    def __init__(self):
        self.todo_list = []
        self.next_id = 1

    async def get_todo_list(self):
        return self.todo_list

    async def add_todo(self, req, appstate):
        newtodo = req.app_state.get("newtodo")
        todo_type = req.app_state.get("todotype")
        if newtodo is None or newtodo == "":
            req.set_data("$.errors", "Please enter a Todo Item")
            return
        if todo_type is None or todo_type == "":
            req.set_data("$.errors", "Please select a Todo Type")
            return
        todo = {"Id": self.next_id, "Item": newtodo, "TodoType": todo_type, "Done": False}
        self.next_id += 1
        self.todo_list.append(todo)
        req.invalidate_data("get-todo-list")
        return

    async def mark_todo_done(self, req, todo_id : int):
        for t in self.todo_list:
            if t["Id"] == todo_id:
                t["Done"] = True
        req.invalidate_data("get-todo-list")
        return

    async def remove_todo(self, req, todo_id : int):
        self.todo_list = [t for t in self.todo_list if t["Id"] != todo_id]
        req.invalidate_data("get-todo-list")
        return

async def main():
    config = dashborg.Config(proc_name="todo", anon_acc=True, auto_keygen=True)
    client = await dashborg.connect_client(config)
    m = TodoModel()
    app = client.app_client().new_app("todo")
    app.set_app_title("Todo App")
    app.set_html(file_name="panels/todo.html", watch=True)
    app.runtime.handler("get-todo-list", m.get_todo_list, pure_handler=True)
    app.runtime.handler("add-todo", m.add_todo)
    app.runtime.handler("mark-todo-done", m.mark_todo_done)
    app.runtime.handler("remove-todo", m.remove_todo)
    await client.app_client().write_app(app, connect=True)
    await client.wait_for_shutdown()

try:
    asyncio.run(main())
except KeyboardInterrupt:
    pass


