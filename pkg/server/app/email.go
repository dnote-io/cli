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
	"fmt"
	"net/url"
	"strings"

	"github.com/dnote/dnote/pkg/server/mailer"
	"github.com/pkg/errors"
)

var defaultSender = "sung@getdnote.com"

// getSenderEmail returns the sender email
func (a *App) getSenderEmail(want string) (string, error) {
	if !a.OnPremise {
		return want, nil
	}

	addr, err := a.getNoreplySender()
	if err != nil {
		return "", errors.Wrap(err, "getting sender email address")
	}

	return addr, nil
}

func getDomainFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing url")
	}

	host := u.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return host, nil
	}
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

	return domain, nil
}

func (a *App) getNoreplySender() (string, error) {
	domain, err := getDomainFromURL(a.WebURL)
	if err != nil {
		return "", errors.Wrap(err, "parsing web url")
	}

	addr := fmt.Sprintf("noreply@%s", domain)
	return addr, nil
}

// SendVerificationEmail sends verification email
func (a *App) SendVerificationEmail(email, tokenValue string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeEmailVerification, mailer.EmailKindText, mailer.EmailVerificationTmplData{
		Token:  tokenValue,
		WebURL: a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := a.getSenderEmail(defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Verify your Dnote email address", from, []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendWelcomeEmail sends welcome email
func (a *App) SendWelcomeEmail(email string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeWelcome, mailer.EmailKindText, mailer.WelcomeTmplData{
		AccountEmail: email,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := a.getSenderEmail(defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Welcome to Dnote!", from, []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}

// SendPasswordResetEmail sends verification email
func (a *App) SendPasswordResetEmail(email, tokenValue string) error {
	body, err := a.EmailTemplates.Execute(mailer.EmailTypeResetPassword, mailer.EmailKindText, mailer.EmailResetPasswordTmplData{
		AccountEmail: email,
		Token:        tokenValue,
		WebURL:       a.WebURL,
	})
	if err != nil {
		return errors.Wrapf(err, "executing reset verification template for %s", email)
	}

	from, err := a.getSenderEmail(defaultSender)
	if err != nil {
		return errors.Wrap(err, "getting the sender email")
	}

	if err := a.EmailBackend.Queue("Reset your password", from, []string{email}, "text/plain", body); err != nil {
		return errors.Wrapf(err, "queueing email for %s", email)
	}

	return nil
}
