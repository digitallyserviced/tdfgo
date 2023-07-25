```python
import ascii_art
import unicode_art
import tools
import constants
import settings
import helpers

CURRENT_ART = None
TOOL_SELECTION = None

def edit_art(art_type, art_data):
    global CURRENT_ART
    global TOOL_SELECTION

    if art_type == 'ascii':
        CURRENT_ART = ascii_art.generate_ascii_art(art_data)
    elif art_type == 'unicode':
        CURRENT_ART = unicode_art.generate_unicode_art(art_data)
    else:
        raise ValueError("Invalid art type. Please choose 'ascii' or 'unicode'.")

def select_tool(tool_name):
    global TOOL_SELECTION
    TOOL_SELECTION = tool_name

def apply_tool(x, y):
    global CURRENT_ART
    global TOOL_SELECTION

    if TOOL_SELECTION is None:
        raise ValueError("No tool selected. Please select a tool first.")

    tool_function = tools.use_tool(TOOL_SELECTION)
    CURRENT_ART = tool_function(CURRENT_ART, x, y)

def save_art(filename):
    with open(filename, 'w') as f:
        f.write(CURRENT_ART)

def load_art(filename):
    global CURRENT_ART
    with open(filename, 'r') as f:
        CURRENT_ART = f.read()
```