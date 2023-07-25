```python
import editor
import ui
import ascii_art
import unicode_art
import tools
import constants
import settings
import helpers

def start_program():
    UI_STATE = settings.get_setting('UI_STATE')
    while UI_STATE != 'EXIT':
        ui.display_ui()
        if UI_STATE == 'EDIT':
            CURRENT_ART = editor.edit_art()
        elif UI_STATE == 'ASCII':
            ascii_art.generate_ascii_art(CURRENT_ART)
        elif UI_STATE == 'UNICODE':
            unicode_art.generate_unicode_art(CURRENT_ART)
        elif UI_STATE == 'TOOL':
            TOOL_SELECTION = tools.use_tool()
        UI_STATE = settings.get_setting('UI_STATE')

if __name__ == "__main__":
    start_program()
```