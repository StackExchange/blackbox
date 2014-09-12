x=/usr/blackbox/bin
if type pathmunge >/dev/null 2>&1;then
pathmunge $x
elif ! grep -sqF :$x:<<<":$PATH:";then
PATH="$x:$PATH"
fi