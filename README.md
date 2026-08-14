### Reflections

The Vim Engine alone is ~400,000 lines of C code. I initially was surprised by this, and I still am. Having created my own editor, I no longer take the small pleasures of Vim for granted, like elegant terminal resizing as well as proper text wrapping and formatting.


### Shortcomings

Gotex has a problem with sentences that wrap around as opposed to if the user manually hits enter when the text goes too rightward. I believe this is because scanner.Scan() breaks text into segments based on newline breaks (\n).


### FIX currRow calculation