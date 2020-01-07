/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
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

import React, { useState, useRef } from 'react';
import { Link, withRouter, RouteComponentProps } from 'react-router-dom';
import classnames from 'classnames';

import { parseSearchString } from 'jslib/helpers/url';
import { getDigestPath } from 'web/libs/paths';
import { blacklist } from 'jslib/helpers/obj';
import SelectMenu from '../../Common/PageToolbar/SelectMenu';
import selectMenuStyles from '../../Common/PageToolbar/SelectMenu.scss';
import { Sort } from '../types';
import styles from './Toolbar.scss';

interface Props extends RouteComponentProps {
  digestUUID: string;
  sort: Sort;
  disabled?: boolean;
}

const SortMenu: React.FunctionComponent<Props> = ({
  digestUUID,
  sort,
  disabled,
  location
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const optRefs = [useRef(null), useRef(null)];
  const searchObj = parseSearchString(location.search);

  const options = [
    {
      name: 'newest',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID, blacklist(searchObj, ['sort']))}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[0]}
          tabIndex={-1}
        >
          Newest
        </Link>
      )
    },
    {
      name: 'oldest',
      value: (
        <Link
          role="menuitem"
          className={selectMenuStyles.link}
          to={getDigestPath(digestUUID, {
            ...searchObj,
            sort: Sort.Oldest
          })}
          onClick={() => {
            setIsOpen(false);
          }}
          ref={optRefs[1]}
          tabIndex={-1}
        >
          Oldest
        </Link>
      )
    }
  ];

  const isActive = sort === Sort.Oldest;

  let defaultCurrentOptionIdx: number;
  let sortText: string;
  if (sort === Sort.Oldest) {
    defaultCurrentOptionIdx = 1;
    sortText = 'Oldest';
  } else {
    defaultCurrentOptionIdx = 0;
    sortText = 'Newest';
  }

  return (
    <SelectMenu
      wrapperClassName={styles['menu-trigger']}
      defaultCurrentOptionIdx={defaultCurrentOptionIdx}
      options={options}
      disabled={disabled}
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      optRefs={optRefs}
      triggerId="sort-menu-trigger"
      menuId="sort-menu"
      headerText="Sort by"
      triggerClassName={classnames('button-no-padding', {
        [styles['active-menu-trigger']]: isActive
      })}
      triggerText={`Sort: ${sortText}`}
      alignment="right"
      direction="bottom"
    />
  );
};

export default withRouter(SortMenu);
