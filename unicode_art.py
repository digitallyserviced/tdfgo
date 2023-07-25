```python
import tools
import helpers

CURRENT_ART = []

def generate_unicode_art():
    global CURRENT_ART
    CURRENT_ART = []
    for i in range(16):
        row = []
        for j in range(16):
            row.append(chr(9608))
        CURRENT_ART.append(row)

def print_art():
    for row in CURRENT_ART:
        print(''.join(row))

def edit_art(x, y, char):
    global CURRENT_ART
    if helpers.is_valid_unicode(char) and helpers.is_valid_position(x, y):
        CURRENT_ART[y][x] = char

def save_art(filename):
    with open(filename, 'w', encoding='utf-8') as f:
        for row in CURRENT_ART:
            f.write(''.join(row) + '\n')

def load_art(filename):
    global CURRENT_ART
    with open(filename, 'r', encoding='utf-8') as f:
        CURRENT_ART = [list(line.rstrip('\n')) for line in f]
```