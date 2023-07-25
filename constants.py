# constants.py

```python
# Constants for the Pixel art ASCII/unicode editor

# ASCII Art Constants
ASCII_CHARS = "!@#$%^&*()_+"

# Unicode Art Constants
UNICODE_CHARS = "▁▂▃▄▅▆▇█"

# UI Constants
UI_WIDTH = 80
UI_HEIGHT = 24

# Tool Constants
TOOLS = ['Pencil', 'Eraser', 'Fill', 'Rectangle', 'Circle', 'Line']

# Settings Constants
DEFAULT_SETTINGS = {
    'width': 80,
    'height': 24,
    'ascii_mode': True,
    'unicode_mode': False,
    'current_tool': 'Pencil'
}

# Current Art Constants
CURRENT_ART = [[' ' for _ in range(DEFAULT_SETTINGS['width'])] for _ in range(DEFAULT_SETTINGS['height'])]

# UI State Constants
UI_STATE = {
    'menu': True,
    'editor': False,
    'preview': False
}
```