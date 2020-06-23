How to add/remove a file into the system?
=========================================

# Adding files:

- If you need to, start the GPG Agent: `eval $(gpg-agent --daemon)`
- Add the file to the system:

```
blackbox file add path/to/file.name.key

# If you want to delete the old plaintext:
blackbox file add --shred path/to/file.name.key
```

Multiple file names can be specified on the command line:

Example 1: Register 2 files:

```
blackbox file add --shred file1.txt file2.txt
```

Example 2: Register all the files in `$DIR`:

```
find $DIR -type f -not -name '*.gpg' -print0 | xargs -0 blackbox file add
```


# Removing files

This command

```
blackbox file remove path/to/file.name.key
```

TODO(tlim): Add examples.

# List files

To see what files are currently enrolled in the system:

```
blackbox file list
```

You can also see their status:

```
blackbox status
blackbox status just_one_file.txt
blackbox status --type ENCRYPTED
```
