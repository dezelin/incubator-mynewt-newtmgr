/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package run

import (
	"os"

	"github.com/spf13/cobra"
	"mynewt.apache.org/newt/newt/builder"
	"mynewt.apache.org/newt/newt/cli"
	"mynewt.apache.org/newt/newt/image"
	"mynewt.apache.org/newt/newt/target"
	"mynewt.apache.org/newt/util"
)

func runRunCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cli.NewtUsage(cmd, util.NewNewtError("Must specify target"))
	}

	t := target.ResolveTargetName(args[0])
	if t == nil {
		cli.NewtUsage(cmd, util.NewNewtError("Invalid target name"+args[0]))
	}

	b, err := builder.NewBuilder(t)
	if err != nil {
		cli.NewtUsage(cmd, err)
	}

	err = b.Build()
	if err != nil {
		cli.NewtUsage(cmd, err)
	}

	/*
	 * Run create-image if version number is specified. If no version number,
	 * remove .img which would'be been created. This so that download script will
	 * barf if it needs an image for this type of target, instead of downloading
	 * an older version.
	 */
	if len(args) > 1 {
		image, err := image.NewImage(b)
		if err != nil {
			cli.NewtUsage(cmd, err)
		}
		err = image.SetVersion(args[1])
		if err != nil {
			cli.NewtUsage(cmd, err)
		}
		err = image.Generate()
		if err != nil {
			cli.NewtUsage(cmd, err)
		}
		err = image.CreateManifest(t)
		if err != nil {
			cli.NewtUsage(cmd, err)
		}
	} else {
		os.Remove(b.AppImgPath())
	}
	err = b.Download()
	if err != nil {
		cli.NewtUsage(cmd, err)
	}
	err = b.Debug()
	if err != nil {
		cli.NewtUsage(cmd, err)
	}
}

func AddCommands(cmd *cobra.Command) {
	runHelpText := "Same as running\n" +
		" - build <target>\n" +
		" - create-image <target> <version>\n" +
		" - download <target>\n" +
		" - debug <target>\n\n" +
		"Note if version number is omitted, create-image step is skipped\n"
	runHelpEx := "  newt run <target-name> [<version>]\n"

	runCmd := &cobra.Command{
		Use:     "run",
		Short:   "build/create-image/download/debug <target>",
		Long:    runHelpText,
		Example: runHelpEx,
		Run:     runRunCmd,
	}
	cmd.AddCommand(runCmd)
}