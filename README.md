# Unique context-based identifier

UCID is a structure for creating ID strings, which contains a context hint (like "user" or "session"),
has a timestamp to minimize conflicts and has a good readability.

## Structure

- context: string, lower case, only letters [a-z], min 2, max 8 characters
- timestamp: Unix timestamp in milliseconds, 11 hex digits, 44 bits
- random part: Random generated number, 5 hex digits, 20 bits

## Decisions

The question was, whould the hex part of the ID split, and when yes, how often.

    user_0123456789abcdef
    or
    user_01234567_89abcdef
    or
    user_0123_4567_89ab_cdef

To make it easier for the developer to see the individual parts,
4 character splits were choosen.

## Should the hex part be split?

### "_" or ":" as separator

The question was, what separator would be better to use.

    user_0123_4567_89ab_cdef
    or
    user:0123:4567:89ab:cdef

It has been decided to use "_" as separator.
The reason for this were:

- Systems often use combinations of letters, numbers and underscores as possible ids
    - We don't want to have struggles with collision, event though the risk of colliding with a system is low
- The individual parts are more easily readable

