/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package login

import (
	"strconv"

	"github.com/dnote/dnote/pkg/cli/client"
	"github.com/dnote/dnote/pkg/cli/consts"
	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var example = `
  dnote login`

var usernameFlag, passwordFlag string

// NewCmd returns a new login command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to dnote server",
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&usernameFlag, "username", "u", "", "email address for authentication")
	f.StringVarP(&passwordFlag, "password", "p", "", "password for authentication")

	return cmd
}

// Do dervies credentials on the client side and requests a session token from the server
func Do(ctx context.DnoteCtx, email, password string) error {
	signinResp, err := client.Signin(ctx, email, password)
	if err != nil {
		return errors.Wrap(err, "requesting session")
	}

	db := ctx.DB
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if err := database.UpsertSystem(tx, consts.SystemSessionKey, signinResp.Key); err != nil {
		return errors.Wrap(err, "saving session key")
	}
	if err := database.UpsertSystem(tx, consts.SystemSessionKeyExpiry, strconv.FormatInt(signinResp.ExpiresAt, 10)); err != nil {
		return errors.Wrap(err, "saving session key")
	}

	tx.Commit()

	return nil
}

func getUsername() (string, error) {
	if usernameFlag != "" {
		return usernameFlag, nil
	}

	var email string
	if err := ui.PromptInput("email", &email); err != nil {
		return "", errors.Wrap(err, "getting email input")
	}
	if email == "" {
		return "", errors.New("Email is empty")
	}

	return email, nil
}

func getPassword() (string, error) {
	if passwordFlag != "" {
		return passwordFlag, nil
	}

	var password string
	if err := ui.PromptPassword("password", &password); err != nil {
		return "", errors.Wrap(err, "getting password input")
	}
	if password == "" {
		return "", errors.New("Password is empty")
	}

	return password, nil
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		log.Plain("Welcome to Dnote Pro (https://www.getdnote.com).\n")

		email, err := getUsername()
		if err != nil {
			return errors.Wrap(err, "getting email input")
		}

		password, err := getPassword()
		if err != nil {
			return errors.Wrap(err, "getting password input")
		}
		if password == "" {
			return errors.New("Password is empty")
		}

		log.Debug("Logging in with email: %s and password: (length %d)\n", email, len(password))

		err = Do(ctx, email, password)
		if errors.Cause(err) == client.ErrInvalidLogin {
			log.Error("wrong login\n")
			return nil
		} else if err != nil {
			return errors.Wrap(err, "logging in")
		}

		log.Success("logged in\n")

		return nil
	}

}
