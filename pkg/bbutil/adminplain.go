// package bbutil
//
// import (
// 	"io/ioutil"
// 	"log"
// 	"path/filepath"
// 	"sort"
// 	"strings"
//
// )
//
// func plainAdminsFile(dir string) string {
// 	return filepath.Join(dir, "blackbox-admins.txt")
// }
//
// // Administrators returns the list of administrators.
// func plainListAdmins(dir string) ([]Administrator, error) {
// 	adminFilename := plainAdminsFile(dir)
// 	d, err := ioutil.ReadFile(adminFilename)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "Could not read the list of administrators")
// 	}
//
// 	// remove a trailing \n.
// 	s := strings.TrimSuffix(string(d), "\n") // remove a single newline.
// 	names := strings.Split(s, "\n")
// 	if !sort.StringsAreSorted(names) {
// 		log.Fatalf("Admin list is corrupted. It is not sorted; %q", adminFilename)
// 	}
// 	r := make([]Administrator, len(names))
// 	for i, name := range names {
// 		r[i].Name = name
// 	}
//
// 	return r, nil
// }
//
// // plainWriteAdmins rewrites the admins file.
// func plainWriteAdmins(dir string, admins []Administrator) error {
// 	return errors.New("UNIMPLEMENTED")
// }
//
// // plainAddAdmins adds one administrator by email address.
// func plainAddAdmin(dir string, admin Administrator) error {
// 	admins, err := plainListAdmins(dir)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Add it to the list, sort it into position.
// 	admins = append(admins, admin)
// 	sort.Slice(admins, func(i, j int) bool {
// 		return admins[i].Name < admins[j].Name
// 	})
//
// 	return plainWriteAdmins(dir, admins)
// }
//
// func plainRemoveAdmin(dir string, admin Administrator) error {
// 	return errors.New("UNIMPLEMENTED")
// }
