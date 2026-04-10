go get -u github.com/stretchr/testify/assert

https://github.com/stretchr/testify

buffer = buffer[parsedBytes:] — This re-slices the buffer. It moves the start pointer forward, but the underlying array stays the same size. The problem is that your buffer now has a reduced capacity. When you later check if the buffer is full and try to grow it, len(buffer) shrinks but the old memory before parsedBytes is still allocated and wasted. Over many iterations, you're slowly losing usable buffer space without actually freeing anything.

copy(buffer, buffer[parsedBytes:]) — This shifts the unparsed data to the front of the buffer, keeping the full capacity available. The buffer stays the same length, and you just adjust readToIndex to reflect how much valid data is in it.
