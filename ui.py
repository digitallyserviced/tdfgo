```python
import curses
from settings import SETTINGS
from constants import CONSTANTS

UI_STATE = {
    "current_menu": "main",
    "current_option": 0
}

def display_ui():
    stdscr = curses.initscr()
    curses.noecho()
    curses.cbreak()
    stdscr.keypad(True)

    while True:
        stdscr.clear()
        if UI_STATE["current_menu"] == "main":
            display_main_menu(stdscr)
        elif UI_STATE["current_menu"] == "edit":
            display_edit_menu(stdscr)
        elif UI_STATE["current_menu"] == "ascii":
            display_ascii_menu(stdscr)
        elif UI_STATE["current_menu"] == "unicode":
            display_unicode_menu(stdscr)
        stdscr.refresh()

def display_main_menu(stdscr):
    stdscr.addstr(0, 0, "Main Menu")
    stdscr.addstr(1, 0, "1. Edit Art")
    stdscr.addstr(2, 0, "2. Generate ASCII Art")
    stdscr.addstr(3, 0, "3. Generate Unicode Art")
    stdscr.addstr(4, 0, "4. Exit")

def display_edit_menu(stdscr):
    stdscr.addstr(0, 0, "Edit Menu")
    stdscr.addstr(1, 0, "1. Select Tool")
    stdscr.addstr(2, 0, "2. Apply Tool")
    stdscr.addstr(3, 0, "3. Save Art")
    stdscr.addstr(4, 0, "4. Back to Main Menu")

def display_ascii_menu(stdscr):
    stdscr.addstr(0, 0, "ASCII Art Menu")
    stdscr.addstr(1, 0, "1. Generate Art")
    stdscr.addstr(2, 0, "2. Save Art")
    stdscr.addstr(3, 0, "3. Back to Main Menu")

def display_unicode_menu(stdscr):
    stdscr.addstr(0, 0, "Unicode Art Menu")
    stdscr.addstr(1, 0, "1. Generate Art")
    stdscr.addstr(2, 0, "2. Save Art")
    stdscr.addstr(3, 0, "3. Back to Main Menu")
```