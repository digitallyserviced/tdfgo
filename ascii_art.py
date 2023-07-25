```python
import tools
import helpers
from constants import CONSTANTS

CURRENT_ART = []

def generate_ascii_art(pixel_data):
    global CURRENT_ART
    ascii_art = ""
    for row in pixel_data:
        for pixel in row:
            ascii_art += helpers.pixel_to_ascii(pixel)
        ascii_art += "\n"
    CURRENT_ART = ascii_art
    return ascii_art

def save_art(filename):
    with open(filename, 'w') as f:
        f.write(CURRENT_ART)

def load_art(filename):
    global CURRENT_ART
    with open(filename, 'r') as f:
        CURRENT_ART = f.read()
    return CURRENT_ART
```