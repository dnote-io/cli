package core

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
)

// ReduceAll reduces all actions
func ReduceAll(ctx infra.DnoteCtx, tx *sql.Tx, actionSlice []actions.Action) error {
	for _, action := range actionSlice {
		if err := Reduce(ctx, tx, action); err != nil {
			return errors.Wrap(err, "reducing action")
		}
	}

	return nil
}

// Reduce transitions the local dnote state by consuming the action returned
// from the server
func Reduce(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var err error

	switch action.Type {
	case actions.ActionAddNote:
		err = handleAddNote(ctx, tx, action)
	case actions.ActionRemoveNote:
		err = handleRemoveNote(ctx, tx, action)
	case actions.ActionEditNote:
		err = handleEditNote(ctx, tx, action)
	case actions.ActionAddBook:
		err = handleAddBook(ctx, tx, action)
	case actions.ActionRemoveBook:
		err = handleRemoveBook(ctx, tx, action)
	default:
		return errors.Errorf("Unsupported action %s", action.Type)
	}

	if err != nil {
		return errors.Wrapf(err, "reducing %s", action.Type)
	}

	return nil
}

func getBookUUIDWithTx(tx *sql.Tx, bookLabel string) (string, error) {
	var ret string
	err := tx.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&ret)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("book '%s' not found", bookLabel)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the book")
	}

	return ret, nil
}

func handleAddNote(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var data actions.AddNoteDataV1
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("reducing add_note. action: %+v. data: %+v\n", action, data)

	bookUUID, err := getBookUUIDWithTx(tx, data.BookName)
	if err != nil {
		return errors.Wrap(err, "getting book uuid")
	}

	var noteCount int
	err = tx.QueryRow("SELECT count(uuid) FROM notes WHERE uuid = ? AND book_uuid = ?", data.NoteUUID, bookUUID).Scan(&noteCount)
	if err != nil {
		return errors.Wrap(err, "counting note")
	}

	if noteCount > 0 {
		// if a duplicate exists, it is because the same action has been previously synced to the server
		// but the client did not bring the bookmark up-to-date at the time because it had error reducing
		// the returned actions.
		// noop so that the client can update bookmark
		return nil
	}

	_, err = tx.Exec(`INSERT INTO notes
	(uuid, book_uuid, content, added_on, public)
	VALUES (?, ?, ?, ?, ?)`, data.NoteUUID, bookUUID, data.Content, action.Timestamp, false)
	if err != nil {
		return errors.Wrap(err, "inserting a note")
	}

	return nil
}

func handleRemoveNote(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var data actions.RemoveNoteDataV1
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("reducing remove_note. action: %+v. data: %+v\n", action, data)

	_, err := tx.Exec("DELETE FROM notes WHERE uuid = ?", data.NoteUUID)
	if err != nil {
		return errors.Wrap(err, "removing a note")
	}

	return nil
}

func buildEditNoteQuery(ctx infra.DnoteCtx, tx *sql.Tx, noteUUID, bookUUID string, ts int64, data actions.EditNoteDataV2) (string, []interface{}, error) {
	setTmpl := "edited_on = ?"
	queryArgs := []interface{}{ts}

	if data.Content != nil {
		setTmpl = fmt.Sprintf("%s, content = ?", setTmpl)
		queryArgs = append(queryArgs, *data.Content)
	}
	if data.Public != nil {
		setTmpl = fmt.Sprintf("%s, public = ?", setTmpl)
		queryArgs = append(queryArgs, *data.Public)
	}
	if data.ToBook != nil {
		bookUUID, err := getBookUUIDWithTx(tx, *data.ToBook)
		if err != nil {
			return "", []interface{}{}, errors.Wrap(err, "getting destination book uuid")
		}

		setTmpl = fmt.Sprintf("%s, book_uuid = ?", setTmpl)
		queryArgs = append(queryArgs, bookUUID)
	}

	queryTmpl := fmt.Sprintf("UPDATE notes SET %s WHERE uuid = ? AND book_uuid = ?", setTmpl)
	queryArgs = append(queryArgs, noteUUID, bookUUID)

	return queryTmpl, queryArgs, nil
}

func handleEditNote(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var data actions.EditNoteDataV2
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("reducing edit_note v2. action: %+v. data: %+v\n", action, data)

	bookUUID, err := getBookUUIDWithTx(tx, data.FromBook)
	if err != nil {
		return errors.Wrap(err, "getting book uuid")
	}

	queryTmpl, queryArgs, err := buildEditNoteQuery(ctx, tx, data.NoteUUID, bookUUID, action.Timestamp, data)
	if err != nil {
		return errors.Wrap(err, "building edit note query")
	}
	_, err = tx.Exec(queryTmpl, queryArgs...)
	if err != nil {
		return errors.Wrap(err, "updating a note")
	}

	return nil
}

func handleAddBook(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var data actions.AddBookDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("reducing add_book. action: %+v. data: %+v\n", action, data)

	var bookCount int
	err = tx.QueryRow("SELECT count(uuid) FROM books WHERE label = ?", data.BookName).Scan(&bookCount)
	if err != nil {
		return errors.Wrap(err, "counting books")
	}

	if bookCount > 0 {
		// If book already exists, another machine added a book with the same name.
		// noop
		return nil
	}

	_, err = tx.Exec("INSERT INTO books (uuid, label) VALUES (?, ?)", utils.GenerateUUID(), data.BookName)
	if err != nil {
		return errors.Wrap(err, "inserting a book")
	}

	return nil
}

func handleRemoveBook(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var data actions.RemoveBookDataV1
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("reducing remove_book. action: %+v. data: %+v\n", action, data)

	var bookUUID string
	err := tx.QueryRow("SELECT uuid FROM books WHERE label = ?", data.BookName).Scan(&bookUUID)
	if err == sql.ErrNoRows {
		// If book does not exist, another client added and removed the book, making the add_book action
		// obsolete. noop.
		return nil
	} else if err != nil {
		return errors.Wrap(err, "querying the book")
	}

	_, err = tx.Exec("DELETE FROM notes WHERE book_uuid = ?", bookUUID)
	if err != nil {
		return errors.Wrap(err, "removing notes")
	}

	_, err = tx.Exec("DELETE FROM books WHERE uuid = ?", bookUUID)
	if err != nil {
		return errors.Wrap(err, "removing a book")
	}

	return nil
}
