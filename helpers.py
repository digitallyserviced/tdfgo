```python
# Importing required modules
import os
import json

def use_helper(helper_name, *args, **kwargs):
    """
    Function to use a helper function based on the helper_name provided.
    """
    helpers = {
        'load_file': load_file,
        'save_file': save_file,
        'clear_screen': clear_screen,
    }

    return helpers[helper_name](*args, **kwargs)

def load_file(file_path):
    """
    Helper function to load a file.
    """
    with open(file_path, 'r') as file:
        return file.read()

def save_file(file_path, content):
    """
    Helper function to save a file.
    """
    with open(file_path, 'w') as file:
        file.write(content)

def clear_screen():
    """
    Helper function to clear the console screen.
    """
    os.system('cls' if os.name == 'nt' else 'clear')

def load_json(file_path):
    """
    Helper function to load a JSON file.
    """
    with open(file_path, 'r') as file:
        return json.load(file)

def save_json(file_path, content):
    """
    Helper function to save a JSON file.
    """
    with open(file_path, 'w') as file:
        json.dump(content, file, indent=4)
```