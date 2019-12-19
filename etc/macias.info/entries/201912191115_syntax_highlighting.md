# Code syntax highlighting available!

After almost two years without updating this blog, I have changed the Markdown
processor from [Blackfriday](github.com/russross/blackfriday) to
[Goldmark](github.com/yuin/goldmark), which includes an extension to highlight
the code syntax according to the [Alec Thomas' Chroma engine and styles](https://github.com/alecthomas/chroma).

For example, the following preformatted Markdown entry:

```
 ```c
 #include<stdio.h>
   
 void main() {
     printf("Hello, world!");
 }
 `` `
```

Will be rendered as:

```c
#include<stdio.h>

void main() {
   printf("Hello, world!");
}
```

Also, few new changes have been added:

* Moved away from [Clean Blog](https://startbootstrap.com/themes/clean-blog/)
  template to a simpler, self-made template. You may think it's horrible, and I
  probably agree 😅.
  
* Added metadata to a header, and a page title that coincides with the entry title.