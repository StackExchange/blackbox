package box

// This file implements the business logic related to a black box.

import "fmt"

func (bx *box) AdminAdd([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminAdd")
}

func (bx *box) AdminList() error {

	admins, err := bx.getAdmins()
	if err != nil {
		return err
	}

	for _, v := range admins {
		fmt.Println(v)
	}
	return nil
}

func (bx *box) AdminRemove([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminRemove")
}

func (bx *box) Cat([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Cat")
}

func (bx *box) Decrypt(names []string, overwrite bool, bulk bool, setgroup string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Decrypt")
}

func (bx *box) Diff([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Diff")
}

func (bx *box) Edit([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Edit")
}

func (bx *box) Encrypt(names []string, bulk bool, setgroup string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Encrypt")
}

func (bx *box) FileAdd(names []string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileAdd")
}

func (bx *box) FileList() error {
	return fmt.Errorf("NOT IMPLEMENTED: FileList")
}

func (bx *box) FileRemove(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileRemove")
}

func (bx *box) Info() error {

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

	return nil
}

func (bx *box) Init() error {
	return fmt.Errorf("NOT IMPLEMENTED: Init")
}

func (bx *box) Reencrypt(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Reencrypt")
}

func (bx *box) Shred(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Shred")
}

func (bx *box) Status(names []string, mode StatusMode, nameOnly bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Status")
}
