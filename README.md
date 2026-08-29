### Reflections

The Vim Engine alone is ~400,000 lines of C code. I initially was surprised by this, and I still am. Having created my own editor, I no longer take the small pleasures of Vim for granted, like elegant terminal resizing as well as proper text wrapping and formatting.


### Shortcomings

strconv could have likely replaced my long list of supported constants.

Gotex has a problem with sentences that wrap around as opposed to if the user manually hits enter when the text goes too rightward.

The size of a tab (\t) is assumed to be 8 pixels. This can and will create problems if you terminal font size is different.

The editor reads in bytes, no runes. Therefore, text in the editor is limited to ASCII.

### Fix/Next Steps: 
- add to func to calculate E.currentWithinRowIndex every time. fix backspacing
- Saving Changes