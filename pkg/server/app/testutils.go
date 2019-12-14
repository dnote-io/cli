/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package app

import (
	"os"

	"github.com/dnote/dnote/pkg/clock"
	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/dnote/dnote/pkg/server/testutils"
)

// NewTest returns an app for a testing environment
func NewTest(appParams *App) App {
	emailTmplDir := os.Getenv("DNOTE_TEST_EMAIL_TEMPLATE_DIR")

	a := App{
		DB:               testutils.DB,
		WebURL:           os.Getenv("WebURL"),
		Clock:            clock.NewMock(),
		EmailTemplates:   mailer.NewTemplates(&emailTmplDir),
		EmailBackend:     &testutils.MockEmailbackendImplementation{},
		StripeAPIBackend: nil,
		OnPremise:       false,
	}

	// Allow to override with appParams
	if appParams != nil && appParams.EmailBackend != nil {
		a.EmailBackend = appParams.EmailBackend
	}
	if appParams != nil && appParams.Clock != nil {
		a.Clock = appParams.Clock
	}
	if appParams != nil && appParams.EmailTemplates != nil {
		a.EmailTemplates = appParams.EmailTemplates
	}
	if appParams != nil && appParams.StripeAPIBackend != nil {
		a.StripeAPIBackend = appParams.StripeAPIBackend
	}
	if appParams != nil && appParams.OnPremise {
		a.OnPremise = appParams.OnPremise
	}
	if appParams != nil && appParams.WebURL != "" {
		a.WebURL = appParams.WebURL
	}

	return a
}
