
Each test does the following:
1. Copy the files from testdata/NNNN
2. Run the command in test_NNNN.sh
3. 


TEST ENROLLMENT:

PHASE 'Alice creates a repo.  She creates secret.txt.'
PHASE 'Alice wants to be part of the secret system.'
PHASE 'She creates a GPG key...'
PHASE 'Initializes BB...'
PHASE 'and adds herself as an admin.'
PHASE 'Bob arrives.'
PHASE 'Bob creates a gpg key.'
PHASE 'Alice does the second part to enroll bob.'
PHASE 'She enrolls bob.'
PHASE 'She enrolls secrets.txt.'
PHASE 'She decrypts secrets.txt.'
PHASE 'She edits secrets.txt.'
PHASE 'Alice copies files to a non-repo directory. (NO REPO)'
PHASE 'Alice shreds these non-repo files. (NO REPO)'
PHASE 'Alice decrypts secrets.txt (NO REPO).'
PHASE 'Alice edits secrets.txt. (NO REPO EDIT)'
PHASE 'Alice decrypts secrets.txt (NO REPO EDIT).'
PHASE 'appears.'
#PHASE 'Bob makes sure he has all new keys.'

TEST INDIVIDUAL COMMANDS:

PHASE 'Bob postdeploys... default.'
PHASE 'Bob postdeploys... with a GID.'
PHASE 'Bob cleans up the secret.'
PHASE 'Bob removes Alice.'
PHASE 'Bob reencrypts files so alice can not access them.'
PHASE 'Bob decrypts secrets.txt.'
PHASE 'Bob edits secrets.txt.'
PHASE 'Bob decrypts secrets.txt VERSION 3.'
PHASE 'Bob exposes a secret in the repo.'
PHASE 'Bob corrects it by registering it.'
PHASE 'Bob enrolls my/path/to/relsecrets.txt.'
PHASE 'Bob decrypts relsecrets.txt.'
PHASE 'Bob enrolls !important!.txt'
PHASE 'Bob enrolls #andpounds.txt'
PHASE 'Bob enrolls stars*bars?.txt'
PHASE 'Bob enrolls space space.txt'
PHASE 'Bob checks out stars*bars?.txt.'
PHASE 'Bob checks out space space.txt.'
PHASE 'Bob shreds all exposed files.'
PHASE 'Bob updates all files.'
PHASE 'Bob DEregisters mistake.txt'
PHASE 'Bob enrolls multiple files: multi1.txt and multi2.txt'
PHASE 'Alice returns. She should be locked out'
PHASE 'Alice tries to decrypt secret.txt. Is blocked.'
