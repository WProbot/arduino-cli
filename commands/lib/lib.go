/*
 * This file is part of arduino-cli.
 *
 * arduino-cli is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 ARDUINO AG (http://www.arduino.cc/)
 */

package lib

import (
	"os"
	"strings"

	"github.com/bcmi-labs/arduino-cli/commands"
	"github.com/bcmi-labs/arduino-cli/commands/version"
	"github.com/bcmi-labs/arduino-cli/common"
	"github.com/bcmi-labs/arduino-cli/common/formatter/output"
	"github.com/bcmi-labs/arduino-cli/common/formatter/pretty_print"
	"github.com/bcmi-labs/arduino-cli/libraries"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	commandVersion string = "0.1.0-alpha.preview"
)

// Init prepares the command.
func Init(rootCommand *cobra.Command) {
	rootCommand.AddCommand(command)
	command.Flags().BoolVar(&flags.updateIndex, "update-index", false, "Updates the libraries index.")
	version.AddVersion(command.Name(), commandVersion)
}

var flags struct {
	updateIndex bool // if true, updates libraries index.
}

var command = &cobra.Command{
	Use:   "lib",
	Short: "Arduino commands about libraries.",
	Long: "Arduino commands about libraries.\n" +
		"Can be used with --update-index flag to update the libraries index too.",
	Run: run,
	Example: "" +
		"arduino lib install YoutubeApi\n" +
		"arduino lib --update-index",
}

func run(cmd *cobra.Command, args []string) {
	if flags.updateIndex {
		logrus.Info("Updating index")
		common.ExecUpdateIndex(prettyPrints.DownloadLibFileIndex())
	} else {
		logrus.Warn("Bad call, executing help command")
		cmd.Help()
		os.Exit(commands.ErrBadCall)
	}
	logrus.Info("Done")
}

func resultFromFileName(file os.FileInfo, libs *output.LibProcessResults) {
	fileName := file.Name()
	//replacing underscore in foldernames with spaces.
	fileName = strings.Replace(fileName, "_", " ", -1)
	fileName = strings.Replace(fileName, "-", " v. ", -1)
	//I use folder name
	logrus.WithField("Name", fileName).Warn("Using filename to get result")
	libs.Libraries = append(libs.Libraries, output.ProcessResult{
		ItemName: fileName,
		Status:   "",
		Error:    "Unknown Version",
	})
}

func getLibStatusContext() (*libraries.StatusContext, error) {
	var index libraries.Index
	err := libraries.LoadIndex(&index)
	if err != nil {
		logrus.WithError(err).Warn("Error during index load, pretty printing error message and trying to recover")
		status, err := prettyPrints.CorruptedLibIndexFix(index)
		if err != nil {
			logrus.WithError(err).Error("Did not recover, returning error")
			return nil, err
		}
		logrus.Warn("Recovered and status context created")
		return &status, nil
	}

	logrus.Info("Creating status context")
	status := index.CreateStatusContext()
	return &status, nil
}