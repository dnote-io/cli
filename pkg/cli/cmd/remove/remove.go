/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

package remove

import (
	"fmt"
	"strconv"

	"github.com/dnote/dnote/pkg/cli/context"
	"github.com/dnote/dnote/pkg/cli/database"
	"github.com/dnote/dnote/pkg/cli/infra"
	"github.com/dnote/dnote/pkg/cli/log"
	"github.com/dnote/dnote/pkg/cli/output"
	"github.com/dnote/dnote/pkg/cli/ui"
	"github.com/dnote/dnote/pkg/cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var targetBookName string

var example = `
  * Delete a note by its id
  dnote delete 2

  * Delete a book
  dnote delete -b js`

// NewCmd returns a new remove command
func NewCmd(ctx context.DnoteCtx) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a note or a book",
		Aliases: []string{"rm", "d", "delete"},
		Example: example,
		RunE:    newRun(ctx),
	}

	f := cmd.Flags()
	f.StringVarP(&targetBookName, "book", "b", "", "The book name to delete")

	return cmd
}

func newRun(ctx context.DnoteCtx) infra.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		if targetBookName != "" {
			if err := removeBook(ctx, targetBookName); err != nil {
				return errors.Wrap(err, "removing the book")
			}

			return nil
		}

		var noteRowIDArg string
		if len(args) == 2 {
			log.Plain(log.ColorYellow.Sprintf("DEPRECATED: you no longer need to pass book name to the view command. e.g. `dnote view 123`.\n\n"))

			noteRowIDArg = args[1]
		} else if len(args) == 1 {
			noteRowIDArg = args[0]
		} else {
			return errors.New("Missing argument")
		}

		noteRowID, err := strconv.Atoi(noteRowIDArg)
		if err != nil {
			return errors.Wrap(err, "invalid rowid")
		}

		if err := removeNote(ctx, noteRowID); err != nil {
			return errors.Wrap(err, "removing the note")
		}

		return nil
	}
}

func removeNote(ctx context.DnoteCtx, noteRowID int) error {
	db := ctx.DB

	noteInfo, err := database.GetNoteInfo(db, noteRowID)
	if err != nil {
		return err
	}

	output.NoteInfo(noteInfo)

	ok, err := ui.Confirm("remove this note?", false)
	if err != nil {
		return errors.Wrap(err, "getting confirmation")
	}
	if !ok {
		log.Warnf("aborted by user\n")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if _, err = tx.Exec("UPDATE notes SET deleted = ?, dirty = ?, body = ? WHERE uuid = ?", true, true, "", noteInfo.UUID); err != nil {
		return errors.Wrap(err, "removing the note")
	}
	tx.Commit()

	log.Successf("removed from %s\n", noteInfo.BookLabel)

	return nil
}

func removeBook(ctx context.DnoteCtx, bookLabel string) error {
	db := ctx.DB

	bookUUID, err := database.GetBookUUID(db, bookLabel)
	if err != nil {
		return errors.Wrap(err, "finding book uuid")
	}

	ok, err := ui.Confirm(fmt.Sprintf("delete book '%s' and all its notes?", bookLabel), false)
	if err != nil {
		return errors.Wrap(err, "getting confirmation")
	}
	if !ok {
		log.Warnf("aborted by user\n")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "beginning a transaction")
	}

	if _, err = tx.Exec("UPDATE notes SET deleted = ?, dirty = ?, body = ? WHERE book_uuid = ?", true, true, "", bookUUID); err != nil {
		return errors.Wrap(err, "removing notes in the book")
	}

	// override the label with a random string
	uniqLabel := utils.GenerateUUID()
	if _, err = tx.Exec("UPDATE books SET deleted = ?, dirty = ?, label = ? WHERE uuid = ?", true, true, uniqLabel, bookUUID); err != nil {
		return errors.Wrap(err, "removing the book")
	}

	tx.Commit()

	log.Success("removed book\n")

	return nil
}
