GIT tips
========


# Configure git to show diffs in encrypted files

It's possible to tell Git to decrypt versions of the file before running them through `git diff` or `git log`. To achieve this do:

- Add the following to `.gitattributes` at the top of the git repository:

```
*.gpg diff=blackbox
```

- Add the following to `.git/config`:

```
[diff "blackbox"]
    textconv = gpg --use-agent -q --batch --decrypt
````

Commands like `git log -p file.gpg` and `git diff master --` will display as expected.
