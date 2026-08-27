### Reflections

The Vim Engine alone is ~400,000 lines of C code. I initially was surprised by this, and I still am. Having created my own editor, I no longer take the small pleasures of Vim for granted, like elegant terminal resizing as well as proper text wrapping and formatting.


### Shortcomings

strconv could have likely replaced my long list of supported constants.

Gotex has a problem with sentences that wrap around as opposed to if the user manually hits enter when the text goes too rightward.

The size of a tab (\t) is assumed to be 8 pixels. This can and will create problems if you terminal font size is different.

The editor reads in bytes, no runes. Therefore, text in the editor is limited to ASCII.

### Fix: 
Tab character messing up checkBounds. checkBounds is also inefficiently called.

Doing the following is problematic, since it doesn't account for the length of the \t. A function could fix this.
E.cursorX = len(E.Lines[E.currentRow-1])

The Fix: Instead of treating the lines as strings, treat them as slices of runes, and if the given coordinate is the rune '/t', act accordingly.