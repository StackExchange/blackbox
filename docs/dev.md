Developer Info
==============

Code submissions are gladly welcomed! The code is fairly easy to read.

Get the code:

```
git clone git@github.com:StackExchange/blackbox.git
```

Test your changes:

```
go test ./...
```

This runs through a number of system tests. It creates a repo,
encrypts files, decrypts files, and so on. You can run these tests to
verify that the changes you made didn't break anything. You can also
use these tests to verify that the system works with a new operating
system.

Please submit tests with code changes:

The best way to change BlackBox is via Test Driven Development. First
add a test to `tools/confidence.sh`. This test should fail, and
demonstrate the need for the change you are about to make. Then fix
the bug or add the feature you want. When you are done, `make
confidence` should pass all tests. The PR you submit should include
your code as well as the new test. This way the confidence tests
accumulate as the system grows as we know future changes don't break
old features.

Note: More info about compatibility are on the [Compatibility Page](compatibility.md)

