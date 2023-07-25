```python
import constants

TOOL_SELECTION = None

def use_tool(tool_name, art):
    global TOOL_SELECTION
    TOOL_SELECTION = tool_name
    if tool_name == constants.TOOL_PENCIL:
        return pencil_tool(art)
    elif tool_name == constants.TOOL_ERASER:
        return eraser_tool(art)
    elif tool_name == constants.TOOL_FILL:
        return fill_tool(art)
    else:
        print("Invalid tool selected.")
        return art

def pencil_tool(art):
    print("Pencil tool is selected.")
    # Add code here to implement pencil tool on the art
    return art

def eraser_tool(art):
    print("Eraser tool is selected.")
    # Add code here to implement eraser tool on the art
    return art

def fill_tool(art):
    print("Fill tool is selected.")
    # Add code here to implement fill tool on the art
    return art
```