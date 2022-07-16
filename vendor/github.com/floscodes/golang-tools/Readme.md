```
func CopyDir(src string, dest string) error
```
This function copies a directory from src(path to source folder) to a destination(dest).

```
func CheckIsSpace(s string) bool
```

This function checks if a string only consists of spaces. If so, the function returns true.


```
func ThousandsSeparate(N interface{}, lang string) (string, error)
```

This function adds thousands-separators to a number N. This function takes two arguments: the number itself and the language-code of the language you want to use. It will return a string containing the number with thousands separators in the language you set with the second argument. For now the package supports English and German. More languages will be added in the future.
