## Text Match.

Utility to compare 2 texts and gives a matching score based on:
 * Cosine Similarity.
 * Absolute Similarity (how many tokens matched).

Before matching, the text goees through the following pipeline of functions.

* Tokenisation (Split the text word by word - every punctuation and the whitespace is used as a delimiter.)
* Filtering Stopwords (English only) - The stopwords in the english language are removed so as to filter noise.
* The text is encoded to it's DoubleMetaphone encoding.
* The text is encoded to it's Stem using porter-stemmer method. 
* The Original text is also retained.

```
a := "Hello Kitty,The Quick Brown Frog Jumped"
b := "Kitty Hello Frogs are always known to be quick"
fmt.Printf("%+v", Similarity(a, b, true)) // last parameter 'true' if a needs to be treated as a master string or a source.

fmt.Printf("%+v", Similarity(a, b, true)) // last parameter 'false' if a does not to be treated as a master string or a source.

```
