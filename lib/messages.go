package lib

import (
	"fmt"
	"strings"
)

func ErrMsgNotFound(data_name string) string {
	return fmt.Sprintf("Data %s tidak ditemukan", data_name)
}

func MsgDeleted(data_name string) string {
	return fmt.Sprintf("Data %s berhasil dihapus", data_name)
}

func MsgUpdated(data_name string) string {
	return fmt.Sprintf("Data %s berhasil diperbarui", data_name)
}
func MsgAdded(data_name string) string {
	return fmt.Sprintf("Data %s berhasil ditambahkan", data_name)
}
func MsgValidUrl(data_name string) string {
	return fmt.Sprintf("%s harus merupakan url yg valid", data_name)
}
func MsgRequired(data_name ...string) string {
	return fmt.Sprintf("%s harus diisi", strings.Join(data_name, ", "))
}

func MsgAlreadyExist(data_name string) string {
	return fmt.Sprintf("%s sudah tersedia, silahkan masukkan %s lain", data_name, data_name)
}

func MsgFailDelete(data_name string) string {
	return fmt.Sprintf("%s gagal dihapus", data_name)
}
