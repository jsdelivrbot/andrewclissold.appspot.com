27 Apr 2014
Some nice code I found out in the wild:

``` javascript
if (!theForm.onsubmit || (theForm.onsubmit() == false))
```
