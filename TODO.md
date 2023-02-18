# progress
it works nicely when there isn't spamming going on.

if spamming:
when pressing g while a message is being received the calculation messes up and the first shown message is actually the 2nd or 3rd
when pressing G while a message is being received a message or two will be missing, im suspecting this happens because of the ToInsert gets written to the db but we already fetched from the db so ToInsert is empty even tho we haven't appened it yet. (same thing happens when going down with d)
-- potential fixes
 - save stuff to db message by message (this is slow as shit, but maybe goroutines could do something here ???)
 - same way of throttling the message receiving when a key is pressed such as g, G, d or etc..
 - maybe there's a way to optimize the readfromdb cmd, maybe the pointer stuff is slowing it down. idk
 - last resort is making this whole feature an option for people that know they have slow chats and want to keep the program open for days (idk)
