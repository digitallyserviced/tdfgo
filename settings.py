```python
# settings.py

# Default settings for the ASCII/Unicode art editor
SETTINGS = {
    "canvas_size": (80, 24),  # Default canvas size (width, height)
    "default_tool": "pencil",  # Default tool to use
    "ascii_palette": " .:-=+*#%@",  # ASCII characters to use for art
    "unicode_palette": " ░▒▓█",  # Unicode characters to use for art
    "color_mode": "grayscale",  # Default color mode
}

def get_setting(setting_name):
    """
    Function to get a setting value.
    """
    return SETTINGS.get(setting_name)
```