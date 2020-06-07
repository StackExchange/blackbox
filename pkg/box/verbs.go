package box

// This file implements the business logic related to a black box.

import "fmt"

func (bx *Box) AdminAdd([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminAdd")
}

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

func (bx *Box) AdminRemove([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminRemove")
}

func (bx *Box) Cat([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Cat")
}

func (bx *Box) Decrypt(names []string, overwrite bool, bulk bool, setgroup string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Decrypt")
}

func (bx *Box) Diff([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Diff")
}

func (bx *Box) Edit([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Edit")
}

func (bx *Box) Encrypt(names []string, bulk bool, setgroup string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Encrypt")
}

func (bx *Box) FileAdd(names []string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileAdd")
}

func (bx *Box) FileList() error {
	return fmt.Errorf("NOT IMPLEMENTED: FileList")
}

func (bx *Box) FileRemove(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileRemove")
}

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

	return nil
}

func (bx *Box) Init() error {
	return fmt.Errorf("NOT IMPLEMENTED: Init")
}

func (bx *Box) Reencrypt(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Reencrypt")
}

func (bx *Box) Shred(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Shred")
}

func (bx *Box) Status(names []string, mode StatusMode, nameOnly bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Status")
}
