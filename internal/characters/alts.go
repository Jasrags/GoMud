package characters

import (
	"os"
	"strconv"

	"github.com/GoMudEngine/GoMud/internal/configs"
	"github.com/GoMudEngine/GoMud/internal/mudlog"
	"github.com/GoMudEngine/GoMud/internal/util"
	"gopkg.in/yaml.v2"
)

func AltsExists(userId int) bool {
	_, err := os.Stat(util.FilePath(string(configs.GetFilePathsConfig().DataFiles), `/users/`, strconv.Itoa(userId)+`.alts.yaml`))

	return !os.IsNotExist(err)
}

func LoadAlts(userId int) []Character {

	if !AltsExists(userId) {
		return nil
	}

	altsFilePath := util.FilePath(string(configs.GetFilePathsConfig().DataFiles), `/users/`, strconv.Itoa(userId)+`.alts.yaml`)

	altsFileBytes, err := os.ReadFile(altsFilePath)
	if err != nil {
		mudlog.Error("LoadAlts", "error", err.Error())
		return nil
	}

	altsRecords := []Character{}

	if err := yaml.Unmarshal(altsFileBytes, &altsRecords); err != nil {
		mudlog.Error("LoadAlts", "error", err.Error())
	}

	return altsRecords

}

func SaveAlts(userId int, alts []Character) bool {

	fileWritten := false
	tmpSaved := false
	tmpCopied := false
	completed := false

	defer func() {
		mudlog.Info("SaveAlts()", "userId", strconv.Itoa(userId), "wrote-file", fileWritten, "tmp-file", tmpSaved, "tmp-copied", tmpCopied, "completed", completed)
	}()

	data, err := yaml.Marshal(&alts)
	if err != nil {
		mudlog.Error("SaveAlts", "error", err.Error())
		return false
	}

	carefulSave := configs.GetFilePathsConfig().CarefulSaveFiles

	path := util.FilePath(string(configs.GetFilePathsConfig().DataFiles), `/users/`, strconv.Itoa(userId)+`.alts.yaml`)

	saveFilePath := path
	if carefulSave { // careful save first saves a {filename}.new file
		saveFilePath += `.new`
	}

	err = os.WriteFile(saveFilePath, data, 0777)
	if err != nil {
		mudlog.Error("SaveAlts", "error", err.Error())
		return false
	}
	fileWritten = true
	if carefulSave {
		tmpSaved = true
	}

	if carefulSave {
		//
		// Once the file is written, rename it to remove the .new suffix and overwrite the old file
		//
		if err := os.Rename(saveFilePath, path); err != nil {
			mudlog.Error("SaveAlts", "error", err.Error())
			return false
		}
		tmpCopied = true
	}

	completed = true

	return true

}
