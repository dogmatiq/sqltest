// +build cgo

package sqltest

func (sqlite3Driver) IsAvailable() bool {
	return true
}
