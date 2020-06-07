package box

// This file implements the business logic related to a black box.

import "fmt"

// AdminAdd adds admins.
func (bx *Box) AdminAdd([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminAdd")
}

// AdminList lists the admin id's.
func (bx *Box) AdminList() error {

	admins, err := bx.getAdmins()
	if err != nil {
		return err
	}

	for _, v := range admins {
		fmt.Println(v)
	}
	return nil
}

// AdminRemove removes an id from the admin list.
func (bx *Box) AdminRemove([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminRemove")
}

// Cat outputs a file, unencrypting if needed.
func (bx *Box) Cat([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Cat")
}

// Decrypt decrypts a file.
func (bx *Box) Decrypt(names []string, overwrite bool, bulk bool, setgroup string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Decrypt")
}

// Diff ...
func (bx *Box) Diff([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Diff")
}

// Edit unencrypts, calls editor, calls encrypt.
func (bx *Box) Edit([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Edit")
}

// Encrypt encrypts a file.
func (bx *Box) Encrypt(names []string, bulk bool, setgroup string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Encrypt")
}

// FileAdd enrolls files.
func (bx *Box) FileAdd(names []string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileAdd")
}

// FileList lists the files.
func (bx *Box) FileList() error {
	return fmt.Errorf("NOT IMPLEMENTED: FileList")
}

// FileRemove de-enrolls files.
func (bx *Box) FileRemove(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileRemove")
}

// Info prints debugging info.
func (bx *Box) Info() error {

	_, err := bx.getAdmins()
	if err != nil {
		logErr.Printf("getAdmins error: %v", err)
	}

	_, err = bx.getFiles()
	if err != nil {
		logErr.Printf("getFiles error: %v", err)
	}

	fmt.Println("BLACKBOX:")
	fmt.Printf("bx.ConfigDir=%q\n", bx.ConfigDir)
	//fmt.Printf("bx.Admins=%q\n", bx.Admins)
	fmt.Printf("len(bx.Admins)=%v\n", len(bx.Admins))
	//fmt.Printf("bx.Files=%q\n", bx.Files)
	fmt.Printf("len(bx.Files)=%v\n", len(bx.Files))
	fmt.Printf("bx.Vcs=%v\n", bx.Vcs)
	fmt.Printf("bx.VcsName=%q\n", bx.VcsName)

	return nil
}

// Init initializes a repo.
func (bx *Box) Init() error {
	return fmt.Errorf("NOT IMPLEMENTED: Init")
}

// Reencrypt decrypts and reencrypts files.
func (bx *Box) Reencrypt(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Reencrypt")
}

// Shred shreds files.
func (bx *Box) Shred(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Shred")
}

// Status prints the status of files.
func (bx *Box) Status(names []string, mode StatusMode, nameOnly bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Status")
}
