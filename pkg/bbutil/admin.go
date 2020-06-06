// package bbutil
//
// // Administrator is a description of the admininstrators.
// type Administrator struct {
// 	Name string
// }
//
// // Administrators returns the administrators of this repo.
// func (bbu *RepoInfo) Administrators() ([]Administrator, error) {
// 	return plainListAdmins(bbu.BlackboxConfigDir)
// }
//
// // AddAdmin adds an administrator to this repo.
// func (bbu *RepoInfo) AddAdmin(admin Administrator) error {
// 	return plainAddAdmin(bbu.BlackboxConfigDir, admin)
// }
//
// // RemoveAdminByName removes an administrator from this repo.
// func (bbu *RepoInfo) RemoveAdminByName(admin Administrator) error {
// 	return plainRemoveAdmin(bbu.BlackboxConfigDir, admin)
// }
