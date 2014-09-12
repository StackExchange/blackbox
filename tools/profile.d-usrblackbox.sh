# Prepend to $PATH.

if type pathmunge > /dev/null 2>&1 ; then
    pathmunge /usr/blackbox/bin
elif ! echo $PATH | grep -Eq "(^|:)/usr/blackbox/bin($|:)" ; then
    PATH=/usr/blackbox/bin:$PATH
fi
