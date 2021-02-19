'use strict';

exports.__esModule = true;

var _extends = Object.assign || function (target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i]; for (var key in source) { if (Object.prototype.hasOwnProperty.call(source, key)) { target[key] = source[key]; } } } return target; };

var _createClass = (function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ('value' in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; })();

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError('Cannot call a class as a function'); } }

function _inherits(subClass, superClass) { if (typeof superClass !== 'function' && superClass !== null) { throw new TypeError('Super expression must either be null or a function, not ' + typeof superClass); } subClass.prototype = Object.create(superClass && superClass.prototype, { constructor: { value: subClass, enumerable: false, writable: true, configurable: true } }); if (superClass) Object.setPrototypeOf ? Object.setPrototypeOf(subClass, superClass) : subClass.__proto__ = superClass; }

var _NestedListItem = require('./NestedListItem');

var _NestedListItem2 = _interopRequireDefault(_NestedListItem);

var _UsersList = require('./UsersList');

var _UsersList2 = _interopRequireDefault(_UsersList);

var _RightPanelCard = require('./RightPanelCard');

var _RightPanelCard2 = _interopRequireDefault(_RightPanelCard);

var _SearchPanel = require('./SearchPanel');

var _SearchPanel2 = _interopRequireDefault(_SearchPanel);

var _Loaders = require('./Loaders');

var _Loaders2 = _interopRequireDefault(_Loaders);

var _TeamCreationForm = require('../TeamCreationForm');

var _TeamCreationForm2 = _interopRequireDefault(_TeamCreationForm);

var _react = require('react');

var _react2 = _interopRequireDefault(_react);

/*
 * Copyright 2007-2017 Charles du Jeu - Abstrium SAS <team (at) pyd.io>
 * This file is part of Pydio.
 *
 * Pydio is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

var _propTypes = require('prop-types');

var _propTypes2 = _interopRequireDefault(_propTypes);

var _pydio = require('pydio');

var _pydio2 = _interopRequireDefault(_pydio);

var _materialUi = require('material-ui');

var _materialUiStyles = require('material-ui/styles');

var _avatarActionsPanel = require('../avatar/ActionsPanel');

var _avatarActionsPanel2 = _interopRequireDefault(_avatarActionsPanel);

var _UserCreationForm = require('../UserCreationForm');

var _UserCreationForm2 = _interopRequireDefault(_UserCreationForm);

var _pydioHttpApi = require('pydio/http/api');

var _pydioHttpApi2 = _interopRequireDefault(_pydioHttpApi);

/**
 * High level component to browse users, groups and teams, either in a large format (mode='book') or a more compact
 * format (mode='selector'|'popover').
 * Address book allows to create external users, teams, and also to browse trusted server directories if Federated Sharing
 * is active.
 */

var _Pydio$requireLib = _pydio2['default'].requireLib('boot');

var PydioContextConsumer = _Pydio$requireLib.PydioContextConsumer;
var PydioContextProvider = _Pydio$requireLib.PydioContextProvider;

var AddressBook = (function (_React$Component) {
    _inherits(AddressBook, _React$Component);

    _createClass(AddressBook, null, [{
        key: 'propTypes',
        value: {
            /**
             * Main instance of pydio
             */
            pydio: _propTypes2['default'].instanceOf(_pydio2['default']),
            /**
             * Display mode, either large (book) or small picker ('selector', 'popover').
             */
            mode: _propTypes2['default'].oneOf(['book', 'selector', 'popover']).isRequired,
            /**
             * Use book mode but display as column
             */
            bookColumn: _propTypes2['default'].bool,
            /**
             * Callback triggered in 'selector' mode whenever an item is clicked.
             */
            onItemSelected: _propTypes2['default'].func,
            /**
             * Display users only, no teams or groups
             */
            usersOnly: _propTypes2['default'].bool,
            /**
             * Choose various user sources, either the local directory or remote ( = trusted ) servers.
             */
            usersFrom: _propTypes2['default'].oneOf(['local', 'remote', 'any']),
            /**
             * Disable the search engine
             */
            disableSearch: _propTypes2['default'].bool,
            /**
             * Theme object passed by muiThemeable() wrapper
             */
            muiTheme: _propTypes2['default'].object,
            /**
             * Will be passed to the Popover object
             */
            popoverStyle: _propTypes2['default'].object,
            /**
             * Used as a button to open the selector in a popover
             */
            popoverButton: _propTypes2['default'].object,
            /**
             * Will be passed to the Popover container object
             */
            popoverContainerStyle: _propTypes2['default'].object,
            /**
             * Will be passed to the Popover Icon Button.
             */
            popoverIconButtonStyle: _propTypes2['default'].object
        },
        enumerable: true
    }, {
        key: 'defaultProps',
        value: {
            mode: 'book',
            usersOnly: false,
            usersFrom: 'any',
            teamsOnly: false,
            disableSearch: false
        },
        enumerable: true
    }]);

    function AddressBook(props) {
        var _this = this;

        _classCallCheck(this, AddressBook);

        _React$Component.call(this, props);

        this.onFolderClicked = function (item) {
            var callback = arguments.length <= 1 || arguments[1] === undefined ? undefined : arguments[1];

            // Special case for teams
            if (_this.props.mode === 'selector' && item.IdmRole && item.IdmRole.IsTeam) {
                _this.onUserListItemClicked(item);
                return;
            }
            _this.setState({ loading: true });

            _Loaders2['default'].childrenAsPromise(item, false).then(function (children) {
                _Loaders2['default'].childrenAsPromise(item, true).then(function (children) {
                    _this.setState({ selectedItem: item, loading: false }, callback);
                });
            });
        };

        this.onUserListItemClicked = function (item) {
            if (_this.props.onItemSelected) {
                var uObject = new PydioUsers.User(item.id, item.label, item.type, item.group, item.avatar, item.temporary, item.external);
                if (item.trusted_server_id) {
                    uObject.trustedServerId = item.trusted_server_id;
                    uObject.trustedServerLabel = item.trusted_server_label;
                }
                uObject._uuid = item.uuid;
                if (item.IdmUser) uObject.IdmUser = item.IdmUser;
                if (item.IdmRole) uObject.IdmRole = item.IdmRole;
                _this.props.onItemSelected(uObject);
            } else {
                _this.setState({ rightPaneItem: item });
            }
        };

        this.onCreateAction = function (item) {
            _this.setState({ createDialogItem: item });
        };

        this.closeCreateDialogAndReload = function () {
            _this.setState({ createDialogItem: null });
            _this.reloadCurrentNode();
        };

        this.onCardUpdateAction = function (item) {
            if (item._parent && item._parent === _this.state.selectedItem) {
                _this.reloadCurrentNode();
            }
        };

        this.onDeleteAction = function (parentItem, selection) {
            var skipConfirm = arguments.length <= 2 || arguments[2] === undefined ? false : arguments[2];

            if (!skipConfirm && !confirm(_this.props.getMessage(278))) {
                return;
            }
            switch (parentItem.actions.type) {
                case 'users':
                    Promise.all(selection.map(function (user) {
                        if (_this.state.rightPaneItem === user) {
                            _this.setState({ rightPaneItem: null });
                        }
                        return _pydioHttpApi2['default'].getRestClient().getIdmApi().deleteIdmUser(user.IdmUser);
                    })).then(function () {
                        _this.reloadCurrentNode();
                    });
                    break;
                case 'teams':
                    Promise.all(selection.map(function (team) {
                        if (_this.state.rightPaneItem === team) {
                            _this.setState({ rightPaneItem: null });
                        }
                        return _pydioHttpApi2['default'].getRestClient().getIdmApi().deleteRole(team.IdmRole.Uuid);
                    })).then(function () {
                        _this.reloadCurrentNode();
                    });
                    break;
                case 'team':
                    Promise.all(selection.map(function (user) {
                        return _pydioHttpApi2['default'].getRestClient().getIdmApi().removeUserFromTeam(parentItem.IdmRole.Uuid, user.IdmUser.Login);
                    })).then(function () {
                        _this.reloadCurrentNode();
                    });
                    break;
                default:
                    break;
            }
        };

        this.openPopover = function (event) {
            _this.setState({
                popoverOpen: true,
                popoverAnchor: event.currentTarget
            });
        };

        this.closePopover = function () {
            _this.setState({ popoverOpen: false });
        };

        this.reloadCurrentNode = function () {
            _this.state.selectedItem.leafLoaded = false;
            _this.state.selectedItem.collectionsLoaded = false;
            _this.onFolderClicked(_this.state.selectedItem, function () {
                if (_this.state.rightPaneItem) {
                    (function () {
                        var rPaneId = _this.state.rightPaneItem.id;
                        var foundItem = null;
                        var leafs = _this.state.selectedItem.leafs || [];
                        var collections = _this.state.selectedItem.collections || [];
                        [].concat(leafs, collections).forEach(function (leaf) {
                            if (leaf.id === rPaneId) foundItem = leaf;
                        });
                        _this.setState({ rightPaneItem: foundItem });
                    })();
                }
            });
        };

        this.reloadCurrentAtPage = function (letterOrRange) {
            _this.state.selectedItem.leafLoaded = false;
            _this.state.selectedItem.collectionsLoaded = false;
            if (letterOrRange === -1) {
                _this.state.selectedItem.currentParams = null;
            } else if (letterOrRange.indexOf('-') !== -1) {
                _this.state.selectedItem.range = letterOrRange;
            } else {
                _this.state.selectedItem.range = null;
                _this.state.selectedItem.currentParams = { alpha_pages: 'true', value: letterOrRange };
            }
            _this.onFolderClicked(_this.state.selectedItem);
        };

        this.reloadCurrentWithSearch = function (value) {
            if (!value) {
                _this.reloadCurrentAtPage(-1);
                return;
            }
            _this.state.selectedItem.leafLoaded = false;
            _this.state.selectedItem.collectionsLoaded = false;
            _this.state.selectedItem.currentParams = { has_search: true, value: value, existing_only: true };
            _this.onFolderClicked(_this.state.selectedItem);
        };

        var pydio = props.pydio;
        var mode = props.mode;
        var usersOnly = props.usersOnly;
        var usersFrom = props.usersFrom;
        var teamsOnly = props.teamsOnly;
        var disableSearch = props.disableSearch;

        var getMessage = function getMessage(id) {
            return props.getMessage(id, '');
        };
        var authConfigs = pydio.getPluginConfigs('core.auth');
        var teamActions = {};
        // Check that user_team_create action is not disabled
        var teamsEditable = pydio.getController().actions.has("user_team_create");
        if (teamsEditable) {
            teamActions = {
                type: 'teams',
                create: '+ ' + getMessage(569),
                remove: getMessage(570),
                multiple: true
            };
        }

        var root = undefined;
        if (teamsOnly) {
            root = {
                id: 'teams',
                label: getMessage(568),
                childrenLoader: _Loaders2['default'].loadTeams,
                _parent: null,
                _notSelectable: true,
                actions: teamActions
            };

            this.state = {
                root: root,
                selectedItem: root,
                loading: false,
                rightPaneItem: null
            };

            return;
        }

        root = {
            id: 'root',
            label: getMessage(592),
            type: 'root',
            collections: []
        };
        if (usersFrom !== 'remote') {
            if (authConfigs.get('USER_CREATE_USERS')) {
                root.collections.push({
                    id: 'ext',
                    label: getMessage(593),
                    icon: 'mdi mdi-account-network',
                    itemsLoader: _Loaders2['default'].loadExternalUsers,
                    _parent: root,
                    _notSelectable: true,
                    actions: {
                        type: 'users',
                        create: '+ ' + getMessage(484),
                        remove: getMessage(582),
                        multiple: true
                    }
                });
            }
            if (!usersOnly) {
                root.collections.push({
                    id: 'teams',
                    label: getMessage(568),
                    icon: 'mdi mdi-account-multiple',
                    childrenLoader: _Loaders2['default'].loadTeams,
                    _parent: root,
                    _notSelectable: true,
                    actions: teamActions
                });
            }
            root.collections.push({
                id: 'PYDIO_GRP_/',
                label: getMessage(584),
                icon: 'mdi mdi-account-box',
                childrenLoader: _Loaders2['default'].loadGroups,
                itemsLoader: _Loaders2['default'].loadGroupUsers,
                _parent: root,
                _notSelectable: true
            });
        }

        var ocsRemotes = pydio.getPluginConfigs('core.ocs').get('TRUSTED_SERVERS');
        if (ocsRemotes && !usersOnly && usersFrom !== 'local') {
            var remotes = JSON.parse(ocsRemotes);
            var remotesNodes = {
                id: 'remotes',
                label: getMessage(594),
                //icon:'mdi mdi-server',
                collections: [],
                _parent: root,
                _notSelectable: true
            };
            for (var k in remotes) {
                if (!remotes.hasOwnProperty(k)) continue;
                remotesNodes.collections.push({
                    id: k,
                    label: remotes[k],
                    icon: 'mdi mdi-server-network',
                    type: 'remote',
                    _parent: remotesNodes,
                    _notSelectable: true
                });
            }
            if (remotesNodes.collections.length) {
                root.collections.push(remotesNodes);
            }
        }

        this.state = {
            root: root,
            selectedItem: mode === 'selector' ? root : root.collections[0],
            loading: false,
            rightPaneItem: null,
            teamsEditable: teamsEditable
        };
    }

    AddressBook.prototype.componentDidMount = function componentDidMount() {
        this.state.selectedItem && this.onFolderClicked(this.state.selectedItem);
    };

    AddressBook.prototype.render = function render() {
        var _this2 = this;

        var _props = this.props;
        var mode = _props.mode;
        var getMessage = _props.getMessage;
        var bookColumn = _props.bookColumn;

        if (mode === 'popover') {

            var popoverStyle = this.props.popoverStyle || {};
            var popoverContainerStyle = this.props.popoverContainerStyle || {};
            var iconButtonStyle = this.props.popoverIconButtonStyle || {};
            var iconButton = _react2['default'].createElement(_materialUi.IconButton, {
                style: _extends({ position: 'absolute', padding: 15, zIndex: 100, right: 0, top: 25, display: this.state.loading ? 'none' : 'initial' }, iconButtonStyle),
                iconStyle: { fontSize: 19, color: 'rgba(0,0,0,0.6)' },
                iconClassName: 'mdi mdi-book-open-variant',
                onClick: this.openPopover
            });
            if (this.props.popoverButton) {
                iconButton = _react2['default'].createElement(this.props.popoverButton.type, _extends({}, this.props.popoverButton.props, { onClick: this.openPopover }));
            }
            var WrappedAddressBook = PydioContextProvider(AddressBook, this.props.pydio);
            return _react2['default'].createElement(
                'span',
                null,
                iconButton,
                _react2['default'].createElement(
                    _materialUi.Popover,
                    {
                        open: this.state.popoverOpen,
                        anchorEl: this.state.popoverAnchor,
                        anchorOrigin: { horizontal: 'right', vertical: 'top' },
                        targetOrigin: { horizontal: 'left', vertical: 'top' },
                        onRequestClose: this.closePopover,
                        style: _extends({ marginLeft: 20 }, popoverStyle),
                        zDepth: 2
                    },
                    _react2['default'].createElement(
                        'div',
                        { style: _extends({ width: 320, height: 420 }, popoverContainerStyle) },
                        _react2['default'].createElement(WrappedAddressBook, _extends({}, this.props, { mode: 'selector', style: { height: 420 } }))
                    )
                )
            );
        }

        var _state = this.state;
        var selectedItem = _state.selectedItem;
        var root = _state.root;
        var rightPaneItem = _state.rightPaneItem;
        var createDialogItem = _state.createDialogItem;
        var teamsEditable = _state.teamsEditable;

        var leftColumnStyle = {
            backgroundColor: _materialUiStyles.colors.grey100,
            width: 256,
            overflowY: 'auto',
            overflowX: 'hidden'
        };
        var centerComponent = undefined,
            rightPanel = undefined,
            leftPanel = undefined,
            topActionsPanel = undefined,
            onEditLabel = undefined;

        if (selectedItem.id === 'search') {

            centerComponent = _react2['default'].createElement(_SearchPanel2['default'], {
                item: selectedItem,
                title: getMessage(583, ''),
                searchLabel: getMessage(595, ''),
                onItemClicked: this.onUserListItemClicked,
                onFolderClicked: this.onFolderClicked,
                mode: mode
            });
        } else if (selectedItem.type === 'remote') {

            centerComponent = _react2['default'].createElement(_SearchPanel2['default'], {
                item: selectedItem,
                params: { trusted_server_id: selectedItem.id },
                searchLabel: getMessage(595, ''),
                title: getMessage(596, '').replace('%s', selectedItem.label),
                onItemClicked: this.onUserListItemClicked,
                onFolderClicked: this.onFolderClicked,
                mode: mode
            });
        } else {

            var emptyStatePrimary = undefined;
            var emptyStateSecondary = undefined;
            var otherProps = {};
            if (selectedItem.id === 'teams') {
                if (teamsEditable) {
                    emptyStatePrimary = getMessage(571, '');
                    emptyStateSecondary = getMessage(572, '');
                } else {
                    emptyStatePrimary = getMessage('571.readonly', '');
                    emptyStateSecondary = getMessage('572.readonly', '');
                }
            } else if (selectedItem.id === 'ext') {
                emptyStatePrimary = getMessage(585, '');
                emptyStateSecondary = getMessage(586, '');
            } else if (selectedItem.IdmUser && selectedItem.IdmUser.IsGroup || selectedItem.id === 'PYDIO_GRP_/' || selectedItem.IdmRole && selectedItem.IdmRole.IsTeam) {
                otherProps = {
                    showSubheaders: true,
                    paginatorType: !(selectedItem.currentParams && selectedItem.currentParams.has_search) && 'alpha',
                    paginatorCallback: this.reloadCurrentAtPage.bind(this),
                    enableSearch: !this.props.disableSearch && !(selectedItem.IdmRole && selectedItem.IdmRole.IsTeam), // do not enable inside teams
                    searchLabel: getMessage(595, ''),
                    onSearch: this.reloadCurrentWithSearch.bind(this)
                };
            }

            if ((mode === 'book' || bookColumn) && selectedItem.IdmRole && selectedItem.IdmRole.IsTeam && teamsEditable) {
                topActionsPanel = _react2['default'].createElement(_avatarActionsPanel2['default'], _extends({}, this.props, {
                    team: selectedItem,
                    userEditable: true,
                    reloadAction: function () {
                        _this2.reloadCurrentNode();
                    },
                    onDeleteAction: function () {
                        if (confirm(_this2.props.getMessage(278))) {
                            (function () {
                                var parent = selectedItem._parent;
                                _this2.setState({ selectedItem: parent }, function () {
                                    _this2.onDeleteAction(parent, [selectedItem], true);
                                });
                            })();
                        }
                    },
                    style: { backgroundColor: 'transparent', borderTop: 0, borderBottom: 0 }
                }));
                onEditLabel = function (item, newLabel) {
                    _pydioHttpApi2['default'].getRestClient().getIdmApi().updateTeamLabel(item.IdmRole.Uuid, newLabel, function () {
                        var parent = selectedItem._parent;
                        _this2.setState({ selectedItem: parent }, function () {
                            _this2.reloadCurrentNode();
                        });
                    });
                };
            }

            centerComponent = _react2['default'].createElement(_UsersList2['default'], _extends({
                item: selectedItem,
                onItemClicked: this.onUserListItemClicked,
                onFolderClicked: this.onFolderClicked,
                onCreateAction: this.onCreateAction,
                onDeleteAction: this.onDeleteAction,
                reloadAction: this.reloadCurrentNode.bind(this),
                onEditLabel: onEditLabel,
                loading: this.state.loading,
                mode: mode,
                bookColumn: bookColumn,
                emptyStatePrimaryText: emptyStatePrimary,
                emptyStateSecondaryText: emptyStateSecondary,
                onClick: this.state.rightPaneItem ? function () {
                    _this2.setState({ rightPaneItem: null });
                } : null,
                actionsPanel: topActionsPanel,
                actionsForCell: this.props.actionsForCell,
                usersOnly: this.props.usersOnly
            }, otherProps));
        }
        var rightPanelStyle = _extends({}, leftColumnStyle, {
            position: 'absolute',
            transformOrigin: 'right',
            backgroundColor: 'white',
            right: 8,
            bottom: 8,
            top: 120,
            zIndex: 2
        });
        if (!rightPaneItem) {
            rightPanelStyle = _extends({}, rightPanelStyle, {
                //transform: 'translateX(256px)',
                transform: 'scale(0)'
            });
        }
        //                width: 0
        rightPanel = _react2['default'].createElement(_RightPanelCard2['default'], {
            pydio: this.props.pydio,
            onRequestClose: function () {
                _this2.setState({ rightPaneItem: null });
            },
            style: rightPanelStyle,
            onCreateAction: this.onCreateAction,
            onDeleteAction: this.onDeleteAction,
            onUpdateAction: this.onCardUpdateAction,
            item: rightPaneItem });
        if (mode === 'book') {
            (function () {
                var nestedRoots = [];
                root.collections.map((function (e) {
                    nestedRoots.push(_react2['default'].createElement(_NestedListItem2['default'], {
                        key: e.id,
                        selected: selectedItem.id,
                        nestedLevel: 0,
                        entry: e,
                        onClick: this.onFolderClicked
                    }));
                    nestedRoots.push(_react2['default'].createElement(_materialUi.Divider, { key: e.id + '-divider' }));
                }).bind(_this2));
                nestedRoots.pop();
                leftPanel = _react2['default'].createElement(
                    _materialUi.Paper,
                    { zDepth: 1, style: _extends({}, leftColumnStyle, { zIndex: 2 }) },
                    _react2['default'].createElement(
                        _materialUi.List,
                        null,
                        nestedRoots
                    )
                );
            })();
        }

        var dialogTitle = undefined,
            dialogContent = undefined,
            dialogBodyStyle = undefined;
        if (createDialogItem) {
            if (createDialogItem.actions.type === 'users') {
                dialogBodyStyle = { display: 'flex', flexDirection: 'column', overflow: 'hidden' };
                dialogTitle = getMessage(484, '');
                dialogContent = _react2['default'].createElement(_UserCreationForm2['default'], {
                    zDepth: 0,
                    style: { display: 'flex', flexDirection: 'column', flex: 1 },
                    newUserName: "",
                    onUserCreated: this.closeCreateDialogAndReload.bind(this),
                    onCancel: function () {
                        _this2.setState({ createDialogItem: null });
                    },
                    pydio: this.props.pydio
                });
            } else if (createDialogItem.actions.type === 'teams') {
                dialogTitle = getMessage(569, '');
                dialogContent = _react2['default'].createElement(_TeamCreationForm2['default'], {
                    onTeamCreated: this.closeCreateDialogAndReload,
                    onCancel: function () {
                        _this2.setState({ createDialogItem: null });
                    }
                });
            } else if (createDialogItem.actions.type === 'team') {
                var selectUser = function selectUser(item) {
                    _pydioHttpApi2['default'].getRestClient().getIdmApi().addUserToTeam(createDialogItem.IdmRole.Uuid, item.IdmUser.Login).then(function () {
                        _this2.reloadCurrentNode();
                    });
                };
                dialogTitle = null;
                dialogContent = _react2['default'].createElement(AddressBook, {
                    pydio: this.props.pydio,
                    mode: 'selector',
                    usersOnly: true,
                    disableSearch: true,
                    onItemSelected: selectUser
                });
            }
        }

        var style = this.props.style || {};
        return _react2['default'].createElement(
            'div',
            { style: _extends({ display: 'flex', height: mode === 'selector' ? 320 : 450 }, style) },
            leftPanel,
            centerComponent,
            rightPanel,
            _react2['default'].createElement(
                _materialUi.Dialog,
                {
                    contentStyle: { width: 380, minWidth: 380, maxWidth: 380, padding: 0 },
                    bodyStyle: _extends({ padding: 0 }, dialogBodyStyle),
                    title: _react2['default'].createElement(
                        'div',
                        { style: { padding: 20 } },
                        dialogTitle
                    ),
                    actions: null,
                    modal: false,
                    open: !!createDialogItem,
                    onRequestClose: function () {
                        _this2.setState({ createDialogItem: null });
                    }
                },
                dialogContent
            )
        );
    };

    return AddressBook;
})(_react2['default'].Component);

exports['default'] = AddressBook = PydioContextConsumer(AddressBook);
exports['default'] = AddressBook = _materialUiStyles.muiThemeable()(AddressBook);
exports['default'] = AddressBook;
module.exports = exports['default'];
