Enabling Blackbox on a Repo
===========================

Overview:
1. Run the initialization command
2. Add at least one admin.
3. Add files. (don't add files before the admins)

The long version:

1. If you don't have a GPG key, set it up using instructions such as:
[Set up GPG key](https://help.github.com/articles/generating-a-new-gpg-key/). \
Now you are ready to go.

1. `cd` into a Git, Mercurial, Subversion or Perforce repository and run `blackbox init`.

1. Add yourself with `blackbox admin add YOUR@EMAIL`

1. Commit the files as directed.

That's it!

At this point you should encrypt a file and make sure you can decrypt
it. This verifies that everything is working as expected.


1. Pick a file to be encrypted. Since this is a test, you might want
   to create a test file.  Call it `secret.txt` and edit the file
   so that it includes your mother's maiden name.  Just kidding!
   Store this sentence: `This is my test file.`

2. Run `blackbox file add secret.txt`

3. Decode the encrypted version: `blackbox cat secret.txt`

The "cat" subcommand only accesses the encrypted (`.gpg`) file and is
a good way to see that the file was encrypted properly.  You should
see `This is my test file.` 

4  Verify that editing the file works.

To view and/or edit a file, run `blackbox edit --shred secret.txt`

Now encrypt it and shred the original:

```
blackbox encrypt --shred secret.txt
```

Now make sure you can decrypt the new file:

```
blackbox cat secret.txt
```

You should see the changed text.

Now commit and push `secret.txt.gpg` and you are done!
