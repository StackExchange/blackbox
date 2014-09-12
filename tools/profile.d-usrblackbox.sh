x=/usr/blackbox/bin
if type pathmunge>/dev/null 2>&1;then
pathmunge $x
else case ":$PATH:" in *:$x:*);;*)PATH="$x:$PATH";;esac;fi